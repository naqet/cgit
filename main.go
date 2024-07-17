package main

import (
	"cgit/internal/cargs"
	"cgit/internal/commands"
	"fmt"
)

func main() {
    args := cargs.InitArgs()
    args.AddCommand(commands.HelpCmd{})

    cmd, err := args.Process()

    if err != nil {
        fmt.Println(err)
        return
    }

    err = cmd.Process(args)

    if err != nil {
        fmt.Println(err)
        return
    }
}
