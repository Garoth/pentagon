package keyvalue

import (
    "os"
    "log"
    "io/ioutil"
)

type FileBackend struct {
    location string
}

func NewFileBackend(location string) *FileBackend {
    me := &FileBackend{location}

    // TODO check location is sane
    if err := os.MkdirAll(location, os.ModeDir | 0777); err != nil {
        log.Fatalln("Couldn't make db location", location, err)
    }
    log.Println("Started database")

    return me
}

// TODO need to return err
func (me *FileBackend) Write(category, key, value string) {
    if err := os.MkdirAll(me.location + "/" + category, os.ModeDir | 0777);
            err != nil {
        log.Fatalln("Couldn't make db category dir", category, err)
    }

    filePath := me.location + "/" + category + "/" + key
    file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_TRUNC, 0666)
    if os.IsNotExist(err) {
        if file, err = os.Create(filePath); err != nil {
            log.Fatalln("Couldn't create db file", filePath)
        }
    } else if err != nil {
        log.Fatalln("Unexpected error opening db file", filePath, err)
    }
    defer file.Close()

    if _, err := file.WriteString(value); err != nil {
        log.Fatalln("Failed to db write", value, err)
    }
}

// TODO need to return err
func (me *FileBackend) Read(category, key string) {
    filePath := me.location + "/" + category + "/" + key

    bytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        log.Println("Couldn't db read file", filePath, err)
    } else {
        log.Println("Read value", string(bytes), "for c", category, "k", key)
    }
}
