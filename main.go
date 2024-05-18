package main

import (
	"fmt"
	"os"

	"github.com/naqet/cgit/internal/commands"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		panic("You need to specify command")
	}

	cmd, ok := commands.Commands[args[0]]

	if !ok {
		fmt.Println("Invalid command")
		return
	}

	err := cmd.Process(args[1:])

	if err != nil {
		panic(err)
	}
}
