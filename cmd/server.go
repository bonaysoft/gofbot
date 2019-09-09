package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"

	"github.com/saltbo/gofbot/api"
	"github.com/saltbo/gofbot/pkg/process"
	"github.com/saltbo/gofbot/robot"
)

var (
	repo      string
	commit    string
	version   string
	buildTime string
)

var pidCtrl = process.New("gofbot.lock")

// define some flags
var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "robots",
		Value: "deployments/robots",
	},
}

// define some commands
var commands = []cli.Command{
	{
		Name:   "reload",
		Usage:  "reload for the config",
		Action: reloadAction,
	},
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("repo: %s\ncommit: %s\nversion: %s\nbuildTime: %s\n", repo, commit, version, buildTime)
	}
	app := cli.NewApp()
	app.Compiled = time.Now()
	app.Copyright = "(c) 2019 yanbo.me"
	app.Flags = flags
	app.Commands = commands
	app.Action = serverRun
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func serverRun(c *cli.Context) {
	robotsPath := c.String("robots")
	robots, err := robot.LoadAndParse(robotsPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := pidCtrl.Save(); err != nil {
		log.Fatal(err)
	}
	defer pidCtrl.Clean()

	server := api.NewServer()
	server.SetupRobots(robots)
	setupSignalHandler(server, robotsPath)

	// startup
	if err := server.Run(":9613"); err != nil {
		log.Fatal(err)
	}

	log.Println("normal exited.")
}

func reloadAction(c *cli.Context) {
	p, err := pidCtrl.Find()
	if err != nil {
		log.Println(err)
		return
	}

	if err := p.Signal(syscall.SIGUSR1); err != nil {
		log.Println(err)
		return
	}
}

func setupSignalHandler(server *api.Server, robotsPath string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		for {
			switch <-ch {
			case syscall.SIGINT, syscall.SIGTERM:
				_ = pidCtrl.Clean()
				signal.Stop(ch)
				_ = server.Shutdown()
				log.Print("system exit.")
				return
			case syscall.SIGUSR1:
				// hot reload
				robots, err := robot.LoadAndParse(robotsPath)
				if err != nil {
					log.Println(err)
					return
				}
				server.SetupRobots(robots)
				log.Printf("config reload.")
			}
		}
	}()
}
