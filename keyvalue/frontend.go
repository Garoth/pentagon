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

func Init() {
    FILE_ACCESS_LOCK <- true
}

func Channels() (chan string, chan string, chan string) {
    read := make(chan string)
    write := make(chan string)
    reply := make(chan string)

    go func() {
        for {
            message := <-read

            cmd := &pentagonmodel.KeyValueReadMessage{}
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                log.Println("Decoding Message:", message, "Error:", err)
                continue
            }

            doRead(cmd, reply)
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

            doWrite(cmd, reply)
        }
    }()

    return read, write, reply
}

func doRead(command *pentagonmodel.KeyValueReadMessage, reply chan string) {
    <-FILE_ACCESS_LOCK

    replyMsg := &pentagonmodel.KeyValueResponse{}
    val, err := BACKEND.Read(command.Category, command.Key)
    if err != nil {
        replyMsg.Success = false
        replyMsg.Error = err.Error()
    } else {
        replyMsg.Success = true
        replyMsg.Key = command.Key
        replyMsg.Value = val
    }

    bytes, err2 := json.Marshal(replyMsg)
    if err2 != nil {
        log.Fatalln("Failed marshalling kv reply", err2)
    }

    reply <- string(bytes)

    FILE_ACCESS_LOCK <- true
}

func doWrite(command *pentagonmodel.KeyValueWriteMessage, reply chan string) {
    <-FILE_ACCESS_LOCK

    BACKEND.Write(command.Category, command.Key, command.Value)

    FILE_ACCESS_LOCK <- true
}
