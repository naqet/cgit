package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/naqet/cgit/internal/repository"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("You need to specify command")
		return
	}

	switch args[0] {
	case "init":
		cmd_init(args[1:])
	}
}

func cmd_init(args []string) {
    path := "."

    if len(args) > 0 {
        path = args[0]
    }
    _, err := repository.InitRepository(path)
    
    if err != nil {
        slog.Error(err.Error())
        os.Exit(1)
    }
}
