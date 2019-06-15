package main

import (
	`flag`
	`fmt`
	"log"
	`os`
)

var (
	showVer   bool
	repo      = ""
	commit    = ""
	version   = "v0.0.0"
	buildTime = ""
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
	s, err := New()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Run())
}