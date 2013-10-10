package git

import (
    "encoding/json"
    "log"

    "github.com/Garoth/pentagon-model"
)

func Init() {
}

func Channels() (chan string, chan string) {
    watch, reply := make(chan string), make(chan string)

    go func() {
        for {
            message := <-watch

            command := &pentagonmodel.GitWatchMessage{}
            if err := json.Unmarshal([]byte(message), &command); err != nil {
                log.Println("Decoding Message:", message, "Error:", err)
                continue
            }

            go doWatch(command, reply)
        }
    }()

    return watch, reply
}

func doWatch(cmd *pentagonmodel.GitWatchMessage, reply chan string) {
    git, err := NewGit(cmd.URL)
    if err != nil {
        log.Println(err)
    }

    git.URL = git.URL
}
