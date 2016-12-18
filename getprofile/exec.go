package getprofile

import (
    "os/exec"
    "syscall"
)

func execExitStatus(err error, status int) bool {
    if exitErr, ok := err.(*exec.ExitError); ok {
        exit := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
        return status == exit
    } else {
        return false
    }
}