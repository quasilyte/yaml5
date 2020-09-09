package main

import (
	"fmt"
	"log"
	"os"

	"github.com/quasilyte/yaml5/cmd/internal/lintcmd"
)

var commands []*subCommand

func main() {
	log.SetFlags(0)

	commands = []*subCommand{
		{
			name:    "help",
			main:    cmdHelp,
			summary: "print linter documentation based on the subject",
		},

		{
			name:    "version",
			main:    cmdVersion,
			summary: "print yaml5 version and exit",
		},

		{
			name:    "lint",
			main:    lintcmd.Main,
			summary: "check whether a YAML complies with the YAML5 restrictions",
		},
	}

	if len(os.Args) < 2 {
		log.Printf("Usage: yaml5 <subcmd> [args...]\n\n")
		printSupportedCommands(commands)
		os.Exit(1)
	}

	subcmdName := os.Args[1]
	subcmd := findSubCommand(commands, subcmdName)
	if subcmd == nil {
		fmt.Printf("Sub-command %s doesn't exist\n\n", subcmdName)
		printSupportedCommands(commands)
		os.Exit(1)
	}

	subcmdIdx := 1 // [0] is program name
	// Erase sub-command argument (index=1) to make it invisible for
	// sub commands themselves.
	os.Args = append(os.Args[:subcmdIdx], os.Args[subcmdIdx+1:]...)

	status, err := subcmd.main()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(status)
}

// Build* are filled during the go build with -ldflags.
var (
	BuildTime    string
	BuildOSUname string
	BuildCommit  string
)

func cmdHelp() (int, error) {
	printSupportedCommands(commands)
	return 0, nil
}

func cmdVersion() (int, error) {
	semver := "v0.5.0"
	description := "(no build info available)"
	if BuildTime != "" {
		description = fmt.Sprintf("(%s %s %s)", BuildCommit, BuildTime, BuildOSUname)
	}
	fmt.Printf("yaml5 %s %s\n", semver, description)
	return 0, nil
}
