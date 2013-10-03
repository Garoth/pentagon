package main

import (
    "flag"
    "log"
    "encoding/json"
    "net/http"

    "code.google.com/p/go.net/websocket"
    "github.com/Garoth/go-signalhandlers"
    "github.com/Garoth/pentagon-model"
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

    http.Handle(HTTP_WEBSOCKET, websocket.Handler(HandleWebSocket))

    if err := http.ListenAndServe(*ADDR, nil); err != nil {
        log.Fatalln("Can't start server:", err)
    }
}

func HandleWebSocket(ws *websocket.Conn) {
    for {
        componentInfo := &pentagonmodel.ComponentInfo{}

        var message string
        if err := websocket.Message.Receive(ws, &message); err != nil {
            log.Println("Reading Socket Error:", err)
            log.Println("Closing connection with", ws.RemoteAddr())
            break
        }

        if err := json.Unmarshal([]byte(message), &componentInfo); err != nil {
            log.Fatalln("Decoding Message Error:", err)
        }

        log.Printf("WebSocket Read: %+v", componentInfo)
    }
}