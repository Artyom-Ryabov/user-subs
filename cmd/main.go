package main

import (
	"log"
	"usersubs/internal"
)

// @title User subscriptions API
// @version 1

// @BasePath /

func main() {
	if err := internal.Start(); err != nil {
		log.Printf("%v\n", err)
	}
}
