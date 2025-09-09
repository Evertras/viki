package main

import "github.com/evertras/viki/cmd/viki/cmds"

func main() {
	if err := cmds.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
