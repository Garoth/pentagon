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
    "pentagon/git"
)

const (
    HTTP_WEBSOCKET = "/websocket"
)

var (
    ADDR = flag.String("port", ":9217", "listening port")
)

func main() {
    log.SetFlags(log.Ltime)
    flag.Parse()

    go signalhandlers.Interrupt()
    go signalhandlers.Quit()

    mail.Init()
    keyvalue.Init()
    git.Init()

    http.Handle(HTTP_WEBSOCKET, websocket.Handler(HandleWebSocket))

    if err := http.ListenAndServe(*ADDR, nil); err != nil {
        log.Fatalln("Can't start server:", err)
    }
}

func HandleWebSocket(ws *websocket.Conn) {
    closed := make(chan bool)

    defer func() {
        log.Println("Closing connection with", ws.RemoteAddr())
        closed <- true
        ws.Close()
    }()

    mail := mail.Channels()
    kvRead, kvWrite, kvReply := keyvalue.Channels()
    gitWatch, gitReply := git.Channels()

    go func() {
        for {
            select {
            case msg := <-kvReply:
                err := websocket.Message.Send(ws, msg)
                if err != nil {
                    log.Println("Couldn't send reply:", err)
                    return
                }

            case msg := <-gitReply:
                err := websocket.Message.Send(ws, msg)
                if err != nil {
                    log.Println("Couldn't send reply:", err)
                    return
                }

            case <-closed:
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

        if err := websocket.Message.Receive(ws, &message); err != nil {
            log.Println("Error reading message:", err)
            continue
        }

        if h.Component == pentagonmodel.COMPONENT_EMAIL {
            if (h.Subcomponent == pentagonmodel.SUBCOMPONENT_EMAIL_MAIN) {
                mail <- message
            }

        } else if h.Component == pentagonmodel.COMPONENT_KV {
            if h.Subcomponent == pentagonmodel.SUBCOMPONENT_KV_READ {
                kvRead <- message
            } else if h.Subcomponent == pentagonmodel.SUBCOMPONENT_KV_WRITE {
                kvWrite <- message
            }

        } else if h.Component == pentagonmodel.COMPONENT_GIT {
            if (h.Subcomponent == pentagonmodel.SUBCOMPONENT_GIT_WATCH) {
                gitWatch <- message
            }

        } else {
            log.Println("Invalid component message receieved")
        }
    }
}
