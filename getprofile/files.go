package getprofile

import (
    "io/ioutil"
)

func setConfig(repoUrl string) error {
    return ioutil.WriteFile(configPath, []byte(repoUrl), 0600)
}

func getConfig() (string, error) {
    if bytes, err := ioutil.ReadFile(configPath); err != nil {
        return "", err
    } else {
        return string(bytes), nil
    }
}