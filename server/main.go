package main

import (
	"log"
	"os"
	`os/signal`
	`syscall`
)

func main() {
	if len(os.Args) == 1 {
		serverRun()
		return
	}

	if err := NewClient().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func serverRun() {
	robots, err := loadRobots("robots")
	if err != nil {
		log.Fatal(err)
	}

	if err := savePid(); err != nil {
		log.Fatal(err)
	}

	server := NewServer()
	setupSignalHandler(server)
	server.SetupRobots(robots)
	if err := server.Run(":9613"); err != nil {
		log.Fatal(err)
	}
}

func setupSignalHandler(server *Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		for {
			sig := <-ch
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				cleanPid();
				signal.Stop(ch)
				server.Shutdown()
				log.Printf("system exit.")
				return
			case syscall.SIGUSR1:
				// hot reload
				robots, err := loadRobots("robots")
				if err != nil {
					log.Println(robots)
					return
				}
				server.SetupRobots(robots)
				log.Printf("config reload.")
			}
		}
	}()
}
