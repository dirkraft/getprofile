package getprofile

import "github.com/daaku/go.copyfile"

var cp = &copyfile.Copy{}

func err2(_ interface{}, err error) error {
    return err
}