package main

import (
	"errors"
	"fmt"
	"io"
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
	Target string
}

// ReadBufferSize is the size of the buffer used when reading files.
const ReadBufferSize = 512

var configHome string
var dataHome string

func init() {
	defaultConfig := fmt.Sprintf("%s/.config", getEnv("HOME", "~"))
	defaultData := fmt.Sprintf("%s/.local/share", getEnv("HOME", "~"))

	configHome = fmt.Sprintf("%s/jump", defaultConfig)
	dataHome = fmt.Sprintf("%s/jump", defaultData)

	if _, err := os.Stat(configHome); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(configHome, 0777)
		} else {
			printError(err)
			os.Exit(1)
		}
	}

	if _, err := os.Stat(dataHome); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dataHome, 0777)
		} else {
			printError(err)
			os.Exit(1)
		}
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	return val
}

type JumpMap map[string][]string

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

	// need to compile paths
	paths := make([]string, 0, 100)
	envPaths := strings.Split(os.Getenv("JUMP_PATH"), ":")
	// TODO: Add file paths
	paths = append(paths, envPaths...)

	jm := indexPaths(paths)

	if opts.Back {
		path, err := back()
		if path != "" {
			cd(path)
		}

		if err != nil {
			printError(err)
		}

		os.Exit(0)
		return
	}

	if opts.Target != "" {
		path, err := jm.Jump(opts.Target)
		if err != nil {
			printError(err)
			os.Exit(1)
			return
		}

		cd(path)
		pushHistory()
		return
	}
}

// Jump attempts to change directories given a name.
func (j JumpMap) Jump(name string) (string, error) {
	if paths, ok := j[name]; ok {
		var path string
		path = paths[0]
		if len(path) > 1 {
			// TODO: Prompt for choice
		}

		return path, nil
	}

	return "", errors.New("directory does not exist")
}

func pushHistory() error {
	file, err := os.OpenFile(fmt.Sprintf("%s/history", dataHome), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	file.WriteString(fmt.Sprintf("%s\n", wd))
	return file.Close()
}

func popHistory() (string, error) {
	f, err := os.OpenFile(fmt.Sprintf("%s/history", dataHome), os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buffer []byte
	info, err := f.Stat()
	if err != nil {
		return "", err
	}

	fsize := info.Size()
	if fsize == 0 {
		return "", errors.New("there is no history to go back to")
	}

	if fsize < ReadBufferSize {
		buffer = make([]byte, fsize)
	} else {
		buffer = make([]byte, ReadBufferSize)
	}

	count, err := f.ReadAt(buffer, fsize-int64(cap(buffer)))
	if err != nil {
		if err != io.EOF {
			return "", err
		}
	}

	idx := int64(count - 2)
	for {
		if buffer[idx] == byte('\n') || idx == 0 {
			break
		}
		idx--
	}

	if idx == 0 {
		err := f.Truncate(0)
		return string(buffer[0 : count-1]), err
	}

	err = f.Truncate(idx + 1)
	return string(buffer[idx+1 : count-1]), err
}

func back() (string, error) {
	return popHistory()
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

func indexPaths(paths []string) JumpMap {
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
