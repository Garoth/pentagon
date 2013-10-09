package keyvalue

import (
    "log"
    "encoding/json"

    "github.com/Garoth/pentagon-model"
)

var (
    FILE_ACCESS_LOCK = make(chan bool, 1)
    BACKEND = NewFileBackend("/Users/athorp/pentagondb")
)

// TODO need some kind of atomic test & write
// TODO need to figure out how to delete a key

func Start() (chan string, chan string) {
    read, write := make(chan string), make(chan string)

    FILE_ACCESS_LOCK <- true

    go func() {
        for {
            message := <-read

            cmd := &pentagonmodel.KeyValueReadMessage{}
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                log.Println("Decoding Message:", message, "Error:", err)
                continue
            }

            handleRead(cmd)
        }
    }()

    go func() {
        for {
            message := <-write

            cmd := &pentagonmodel.KeyValueWriteMessage{}
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                log.Println("Decoding Message:", message, "Error:", err)
                continue
            }

            handleWrite(cmd)
        }
    }()

    return read, write
}

func handleRead(command *pentagonmodel.KeyValueReadMessage) {
    <-FILE_ACCESS_LOCK

    BACKEND.Read(command.Category, command.Key)

    FILE_ACCESS_LOCK <- true
}

func handleWrite(command *pentagonmodel.KeyValueWriteMessage) {
    <-FILE_ACCESS_LOCK

    BACKEND.Write(command.Category, command.Key, command.Value)

    FILE_ACCESS_LOCK <- true
}
