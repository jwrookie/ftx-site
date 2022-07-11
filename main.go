package main

import "github.com/foxdex/ftx-site/cmd"

func main() {
	if err := cmd.ApiCmd.Execute(); err != nil {
		panic(err)
	}
}
