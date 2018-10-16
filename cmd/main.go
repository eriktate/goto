package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type options struct {
	Help   bool
	Back   bool
	Alias  string
	Clean  bool
	Add    bool
	Remove bool
	Path   string
}

func main() {
	log.Printf("Arg count: %d", len(os.Args))
	opts, err := parseArgs(os.Args)
	if err != nil {
		log.Print(err)
	}

	fmt.Printf("%+v", opts)
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
		}
	}

	return opts, nil
}
