package main

import (
	"github.com/TskFok/OpenAi/bootstrap"
	"github.com/TskFok/OpenAi/cmd"
)

func main() {
	bootstrap.Init()

	cmd.Execute()
}
