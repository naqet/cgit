package main

import (
	"cgit/internal/cargs"
	"cgit/internal/commands"
	"fmt"
)

func main() {
    args := cargs.InitArgs()
    args.AddCommand(commands.HelpCmd{})
    args.AddCommand(commands.InitCmd{})
    args.AddCommand(commands.HashObjectCmd{})
    args.AddCommand(commands.CatFileCmd{})
    args.AddCommand(commands.LsTreeCmd{})
    // TODO: after implementing commits
    //args.AddCommand(commands.LogCmd{})

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
