package svc

import (
	"errors"
	"fmt"
	"os"

	"github.com/eriktate/jump"
)

type jumpSvc struct {
	jump.Jumper // TODO: Remove once this struct fully implements the interface.
	names       map[string][]string
}

// NewJumpSvc returns a properly initialized jumpSvc.
func NewJumpSvc(paths []string) jumpSvc {
	return jumpSvc{
		names: indexPaths(paths),
	}
}

// Jump attempts to change directories given a name.
func (j jumpSvc) Jump(name string) (string, error) {
	if paths, ok := j.names[name]; ok {
		var path string
		path = paths[0]
		if len(path) > 1 {
			// TODO: Prompt for choice
		}

		return path, nil
	}

	return "", errors.New("directory does not exist")
}

func indexPaths(paths []string) map[string][]string {
	result := make(map[string][]string)
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			// TODO: Log error in verbose mode.
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			// TODO: Log error in verbose mode.
			continue
		}

		files, err := file.Readdir(0)
		if err != nil {
			// TODO: Log error in verbose mode.
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				continue
			}

			name := file.Name()
			if len(name) == 0 || name[0] == '.' { // TODO: Add option to include hidden directories?
				continue
			}

			if _, ok := result[name]; ok {
				result[name] = append(result[name], fmt.Sprintf("%s/%s", path, name)) // TODO: Use something faster than Sprintf
			} else {
				result[name] = []string{fmt.Sprintf("%s/%s", path, name)}
			}
		}
	}

	return result
}
