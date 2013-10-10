package git

import (
    "net/url"
    "path"
    "log"
    "os/exec"
    "bytes"

    "github.com/Garoth/pentagon-model"
)

var (
    ACTIVE_GITS map[string] Git
    GIT = "git"
)

type Git struct {
    URL, LocalPath string
}

func NewGit(targetUrl string) (*Git, error) {
    me := &Git{}
    me.URL = targetUrl

    parsedUrl, err := url.Parse(targetUrl)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    me.LocalPath = pentagonmodel.GetConfig().Workdir + "/git/" +
        path.Base(parsedUrl.Path)

    if err = me.Clone(); err != nil {
        log.Println(err)
        return nil, err
    }

    return me, nil
}

func (me *Git) Clone() error {
    cmd := exec.Command(GIT, "clone", me.URL, me.LocalPath)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    if err != nil {
        log.Println(err)
        return err
    }

    // TODO check for error
    return nil
}

// TODO add functions to look at most recent commit
