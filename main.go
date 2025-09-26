package main

import (
	"log"

	"github.com/aiven/go-api-schemas/cmd"
)

func main() {
	c := cmd.NewCmdRoot()
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
