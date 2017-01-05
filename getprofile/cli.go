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
    app.Version = "0.0.1-dev"
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
            Usage: "Set up getprofile.",
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
            Name: "add",
            Usage: "Track a file.",
            ArgsUsage: "FILE",
            Action: func(ctx *cli.Context) error {
                if _, err := getConfig(); err != nil {
                    dbg(err)
                    return errors.New("Run 'config' first")
                } else if file := strings.TrimSpace(ctx.Args().First()); file == "" {
                    return errors.New("FILE is required. See --help")
                } else {
                    return Track(file)
                }
            },
        }, {
            Name: "rm",
            Usage: "Untrack a file. The local copy is not deleted.",
            ArgsUsage: "FILE",
            Action: func(ctx *cli.Context) error {
                if _, err := getConfig(); err != nil {
                    dbg(err)
                    return errors.New("Run 'config' first")
                } else if file := strings.TrimSpace(ctx.Args().First()); file == "" {
                    return errors.New("FILE is required. See --help")
                } else {
                    return Untrack(file)
                }
            },
        }, {
            Name: "push",
            Usage: "Push profile updates to remote.",
            Action: func(ctx *cli.Context) error {
                if _, err := getConfig(); err != nil {
                    dbg(err)
                    return errors.New("Run 'config' first")
                } else {
                    return Push()
                }
            },
        }, {
            Name: "pull",
            Usage: "Pull profile update from remote. Overwrites local files with updates. " +
                "Does nothing if there is no remote update.",
            Flags: []cli.Flag{
                cli.BoolFlag{
                    Name: "force, f",
                    Usage: "Copy from repo to local whether or not there is an update",
                },
            },
            Action: func(ctx *cli.Context) error {
                if _, err := getConfig(); err != nil {
                    dbg(err)
                    return errors.New("Run 'config' first")
                } else {
                    return Pull(ctx.IsSet("force"))
                }
            },
        },
    }
    app.Run(os.Args)
}