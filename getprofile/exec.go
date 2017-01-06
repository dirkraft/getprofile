package getprofile

import (
    "os/exec"
    "syscall"
)

func execWithDebug(name string, arg ...string) error {
    dbgf("Command: %v % v", name, arg)
    out, err := exec.Command(name, arg...).CombinedOutput()
    dbgf("%s", out)
    return err
}

func execExitStatus(err error, status int) bool {
    if exitErr, ok := err.(*exec.ExitError); ok {
        exit := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
        return status == exit
    } else {
        return false
    }
}