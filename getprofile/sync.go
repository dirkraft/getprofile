package getprofile

import (
    "errors"
    "os/user"
    "path/filepath"
)

type syncer interface {
    Supports(string) bool
    Init() error
    Track(absPath, relPath string) error
    Untrack(relPath string) error
    Push() error
    Pull(force bool) error
}

var syncers = []syncer{
    gitSyncerInstance,
}

func getSyncer() (syncer, error) {
    if repoUrl, err := getConfig(); err != nil {
        return nil, err
    } else {
        for _, impl := range syncers {
            if impl.Supports(repoUrl) {
                return impl, nil
            }
        }
        return nil, errors.New("No supporting syncer for " + repoUrl)
    }
}

func Init() error {
    if syncer, err := getSyncer(); err != nil {
        return err
    } else {
        return syncer.Init()
    }
}

func Track(p string) error {
    if relPath, err := relativize(p); err != nil {
        return err
    } else if syncer, err := getSyncer(); err != nil {
        return err
    } else if absPath, err := filepath.Abs(p); err != nil {
        return err
    } else {
        return syncer.Track(absPath, relPath)
    }
}

func Untrack(p string) error {
    if relPath, err := relativize(p); err != nil {
        return err
    } else if syncer, err := getSyncer(); err != nil {
        return err
    } else {
        return syncer.Untrack(relPath)
    }
}

func Push() error {
    if syncer, err := getSyncer(); err != nil {
        return err
    } else {
        return syncer.Push()
    }
}

func Pull(force bool) error {
    if syncer, err := getSyncer(); err != nil {
        return err
    } else {
        return syncer.Pull(force)
    }
}

func relativize(p string) (string, error) {
    if abs, err := filepath.Abs(p); err != nil {
        return "", err
    } else if curUser, err := user.Current(); err != nil {
        return "", err
    } else {
        return filepath.Rel(curUser.HomeDir, abs)
    }
}

