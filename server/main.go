package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	showVer   bool
	repo      string
	commit    string
	version   string
	buildTime string
)

func init() {
	flag.BoolVar(&showVer, "v", false, "show build version")
	flag.Parse()

	if showVer {
		fmt.Printf("repo: %s\ncommit: %s\nversion: %s\nbuildTime: %s\n", repo, commit, version, buildTime)
		os.Exit(0)
	}
}

func main() {
	robots, err := loadRobots("robots")
	if err != nil {
		log.Fatal(err)
	}

	if s, err := New(robots); err != nil {
		log.Fatal(err)
	} else {
		log.Fatal(s.Run(":9613"))
	}
}
