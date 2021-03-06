package getprofile

import (
    "os/user"
    "path"
    "os"
)

var homePath string
var basePath string
var configPath string
var repoPath string

func initialize() error {
    if curUser, err := user.Current(); err != nil {
        return err
    } else {
        homePath = curUser.HomeDir
        basePath = path.Join(homePath, ".getprofile")
        configPath = path.Join(basePath, "config")
        repoPath = path.Join(basePath, "repo")
    }

    dbg("homePath:%v", homePath)
    dbg("basePath:%v", basePath)
    dbg("configPath:%v", configPath)
    dbg("repoPath:%v", repoPath)

    return os.MkdirAll(basePath, 0700)
}