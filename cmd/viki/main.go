package main

import (
	"log"

	"github.com/evertras/viki/cmd/viki/cmds"
)

func main() {
	if err := cmds.RootCmd.Execute(); err != nil {
		log.Fatalln("Failed:", err)
	}
}
