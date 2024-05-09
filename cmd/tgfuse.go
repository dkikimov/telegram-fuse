package main

import (
	"log"

	"telegram-fuse/cmd/app"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		log.Fatalf("couldn't start program: %s", err)
	}
}
