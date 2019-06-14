package main

import (
	"gofrobot/server"
	"log"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Run())
}
