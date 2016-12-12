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
    app.Commands = []cli.Command{
        {
            Name: "config",
            Usage: "Set up getprofile",
            ArgsUsage: "REPOSITORY_URL",
            Description: usageConfigDesc,
            Action: func(ctx *cli.Context) error {
                if repoUrl := strings.TrimSpace(ctx.Args().First()); repoUrl == "" {
                    return errors.New("REPOSITORY_URL is required")
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
                if file := strings.TrimSpace(ctx.Args().First()); file == "" {
                    return errors.New("FILE is required")
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
                return Sync()
            },
        },
    }
    app.Run(os.Args)
}