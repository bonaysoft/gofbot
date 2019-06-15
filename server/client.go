package main

import (
	`fmt`
	`github.com/urfave/cli`
	`log`
	`syscall`
)

var (
	repo      string
	commit    string
	version   string
	buildTime string
)

// define some flags
var flags = []cli.Flag{

}
// define some commands
var commands = []cli.Command{
	{
		Name:  "reload",
		Usage: "reload for the config",
		Action: func(c *cli.Context) {
			p, err := findProcess()
			if err != nil {
				log.Println(err)
				return
			}

			if err := p.Signal(syscall.SIGUSR1); err != nil {
				log.Println(err)
				return
			}
		},
	},
}

func NewClient() *cli.App {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("repo: %s\ncommit: %s\nversion: %s\nbuildTime: %s\n", repo, commit, version, buildTime)
	}
	app := cli.NewApp()
	app.Flags = flags
	app.Commands = commands
	return app
}
