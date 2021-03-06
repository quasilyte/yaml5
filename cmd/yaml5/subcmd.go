package main

import (
	"fmt"
)

type subCommand struct {
	name     string
	main     func() (int, error)
	summary  string
	examples []subCommandExample
}

type subCommandExample struct {
	comment string
	line    string
}

func findSubCommand(list []*subCommand, name string) *subCommand {
	for _, cmd := range list {
		if cmd.name == name {
			return cmd
		}
	}
	return nil
}

func printSupportedCommands(list []*subCommand) {
	fmt.Printf("Supported sub-commands:\n")
	for _, cmd := range list {
		fmt.Printf("\n\tyaml5 %s\n", cmd.name)
		fmt.Printf("\tDescription: %s.\n", cmd.summary)
		for _, ex := range cmd.examples {
			fmt.Printf("\t%s:\n", ex.comment)
			fmt.Printf("\t\t$ yaml5 %s %s\n", cmd.name, ex.line)
		}
	}
}
