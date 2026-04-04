package main

import (
	"log"

	"github.com/skygenesisenterprise/aether-bank/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
