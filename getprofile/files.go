package getprofile

import (
    "os"
    "os/user"
    "path"
    "io/ioutil"
)

func getConfigPath() (string, string, error) {
    curUser, err := user.Current()
    if err != nil {
        return "", "", err
    }

    basePath := path.Join(curUser.HomeDir, ".getprofile")
    if err := os.MkdirAll(basePath, 0700); err != nil {
        return "", "", err
    }

    return basePath, path.Join(basePath, "config"), nil
}

func setConfig(repoUrl string) error {
    if _, configPath, err := getConfigPath(); err != nil {
        return err
    } else {
        return ioutil.WriteFile(configPath, []byte(repoUrl), 0600)
    }
}

func getConfig() (string, error) {
    if _, configPath, err := getConfigPath(); err != nil {
        return "", err
    } else if bytes, err := ioutil.ReadFile(configPath); err != nil {
        return "", err
    } else {
        return string(bytes), nil
    }
}