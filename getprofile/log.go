package getprofile

import (
    "fmt"
    "strings"
)

var logLevel = 2

func logFormat(args ...interface{}) string {
    s := make([]string, len(args))
    for i, arg := range args {
        s[i] = fmt.Sprintf("%v", arg)
    }
    return strings.Join(s, " ")
}

func dbg(args ...interface{}) {
    if logLevel > 1 {
        return
    }
    fmt.Println("[DBG] " + logFormat(args...))
}
func dbgf(s string, args ...interface{}) {
    if logLevel > 1 {
        return
    }
    dbg(fmt.Sprintf(s, args...))
}
func inf(args ...interface{}) {
    if logLevel > 2 {
        return
    }
    fmt.Println("[INF] " + logFormat(args...))
}

func inff(s string, args ...interface{}) {
    if logLevel > 2 {
        return
    }
    inf(fmt.Sprintf(s, args...))
}