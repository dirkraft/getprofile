package getprofile

import (
    "os"
    "github.com/urfave/cli"
    "strings"
    "errors"
)

func RunCli() {
    if err := initialize(); err != nil {
        panic(err)
    }

    app := cli.NewApp()
    app.Name = "getprofile"
    app.Version = "0.0.1"
    app.Flags = []cli.Flag{
        cli.BoolFlag{
            Name:"verbose",
            Usage: "Show verbose logging",
        },
    }
    app.Before = func(ctx *cli.Context) error {
        if ctx.IsSet("verbose") {
            logLevel = 1
        }
        return nil
    }
    app.Commands = []cli.Command{
        {
            Name: "config",
            Usage: "Set up getprofile",
            ArgsUsage: "REPOSITORY_URL",
            Description: usageConfigDesc,
            Action: func(ctx *cli.Context) error {
                if repoUrl := strings.TrimSpace(ctx.Args().First()); repoUrl == "" {
                    return errors.New("REPOSITORY_URL is required. See --help")
                } else if err := setConfig(repoUrl); err != nil {
                    return err
                } else {
                    return Init()
                }
            },
        }, {
            Name: "track",
            Usage: "Track or untrack a file",
            Flags: []cli.Flag{
                cli.BoolFlag{
                    Name: "untrack, u",
                    Usage: "Stop tracking a tracked file. The file is not deleted from the local machine.",
                },
                // TODO recursive, r
            },
            ArgsUsage: "FILE",
            Action: func(ctx *cli.Context) error {
                if _, err := getConfig(); err != nil {
                    dbg(err)
                    return errors.New("Run 'config' first")
                } else if file := strings.TrimSpace(ctx.Args().First()); file == "" {
                    return errors.New("FILE is required. See --help")
                } else if ctx.Bool("delete") {
                    return Untrack(file)
                } else {
                    return Track(file)
                }
            },
        }, {
            Name: "sync",
            Usage: "Synchronize profile",
            Flags: []cli.Flag{
                cli.BoolFlag{
                    Name: "watch, w",
                    Usage: "Continuously watch and synchronize changes",
                },
            },
            Action: func(ctx *cli.Context) error {
                if _, err := getConfig(); err != nil {
                    dbg(err)
                    return errors.New("Run 'config' first")
                } else {
                    return Sync()
                }
            },
        },
    }
    app.Run(os.Args)
}