package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eriktate/jump/svc"
)

type options struct {
	Help   bool
	Back   bool
	Alias  string
	Clean  bool
	Add    bool
	Remove bool
	Path   string
	Target string
}

func printError(err error) {
	fmt.Printf("echo \"%s\"", err)
}

func cd(path string) {
	fmt.Printf("cd %s", path)
}

func main() {
	opts, err := parseArgs(os.Args)
	if err != nil {
		printError(err)
	}

	envPaths := strings.Split(os.Getenv("JUMP_PATH"), ":")

	j := svc.NewJumpSvc(envPaths)

	if opts.Target != "" {
		path, err := j.Jump(opts.Target)
		if err != nil {
			printError(err)
			os.Exit(1)
			return
		}

		cd(path)
	}
}

func parseArgs(args []string) (options, error) {
	var opts options
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		switch parts[0] {
		case "-b", "--back":
			opts.Back = true
		case "-h", "--help":
			opts.Help = true
		case "-l", "--alias":
			if len(parts) > 1 {
				opts.Alias = parts[1]
			} else {
				return opts, errors.New("You must provide an alias")
			}
		case "-a", "--add":
			if len(parts) > 1 && parts[1] != "" {
				opts.Alias = parts[1]
			}
			opts.Add = true
		case "--clean":
			opts.Clean = true
		case "--remove":
			if len(parts) > 1 && parts[1] != "" {
				opts.Path = parts[1]
			}
			opts.Remove = true
		default:
			if parts[0][0] == '-' {
				return opts, fmt.Errorf("unrecognized flag: %s", parts[0])
			}

			opts.Target = parts[0]
		}
	}

	return opts, nil
}
