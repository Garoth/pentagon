package main

import (
    "flag"
    "log"
    "encoding/json"
    "net/http"

    "code.google.com/p/go.net/websocket"
    "github.com/Garoth/go-signalhandlers"
    "github.com/Garoth/pentagon-model"

    "pentagon/mail"
    "pentagon/keyvalue"
)

const (
    HTTP_WEBSOCKET = "/websocket"
)

var (
    ADDR = flag.String("port", ":9217", "listening port")
    MAIL_CHANNEL_MAIN chan string
    KV_CHANNEL_READ, KV_CHANNEL_WRITE, KV_CHANNEL_REPLY chan string
)

func main() {
    log.SetFlags(log.Ltime)
    flag.Parse()

    go signalhandlers.Interrupt()
    go signalhandlers.Quit()

    MAIL_CHANNEL_MAIN = mail.Start()
    KV_CHANNEL_READ, KV_CHANNEL_WRITE, KV_CHANNEL_REPLY = keyvalue.Start()

    http.Handle(HTTP_WEBSOCKET, websocket.Handler(HandleWebSocket))

    if err := http.ListenAndServe(*ADDR, nil); err != nil {
        log.Fatalln("Can't start server:", err)
    }
}

// TODO reply messages might get sent to wrong socket, based on whoever
// reads them first. There needs to be a more explicit way of saying
// who you're replying to
func HandleWebSocket(ws *websocket.Conn) {
    closed := make(chan bool)

    defer func() {
        log.Println("Closing connection with", ws.RemoteAddr())
        closed <- true
        ws.Close()
    }()

    go func() {
        for {
            select {
            case msg := <-KV_CHANNEL_REPLY:
                log.Println("Sending reply to client:", msg)
                err := websocket.Message.Send(ws, msg)
                if err != nil {
                    log.Println("Couldn't send reply:", err)
                    return
                }

            case <-closed:
                log.Println("Reply thread noticed conn closed")
                return
            }
        }
    }()

    for {
        h := &pentagonmodel.ClientHeader{}

        var message string
        if err := websocket.Message.Receive(ws, &message); err != nil {
            if err.Error() != "EOF" {
                log.Println("Reading Socket Error:", err)
            }
            return
        }

        if err := json.Unmarshal([]byte(message), &h); err != nil {
            log.Println("Decoding Message:", message, "Error:", err)
            continue
        }

        if h.Component == pentagonmodel.COMPONENT_EMAIL {
            if err := websocket.Message.Receive(ws, &message); err != nil {
                log.Println("Error reading mail message:", err)
                continue
            }

            if (h.Subcomponent == pentagonmodel.SUBCOMPONENT_EMAIL_MAIN) {
                MAIL_CHANNEL_MAIN <- message
            }

        } else if h.Component == pentagonmodel.COMPONENT_KV {
            if err := websocket.Message.Receive(ws, &message); err != nil {
                log.Println("Error reading mail message:", err)
                continue
            }

            if h.Subcomponent == pentagonmodel.SUBCOMPONENT_KV_READ {
                KV_CHANNEL_READ <- message
            } else if h.Subcomponent == pentagonmodel.SUBCOMPONENT_KV_WRITE {
                KV_CHANNEL_WRITE <- message
            }

        } else {
            log.Println("Invalid component message receieved")
        }
    }
}
