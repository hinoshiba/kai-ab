package main

import (
	"os"
	"fmt"
	"flag"
	"path/filepath"
)

const (
	VERSION_TEMPLATE string = "v0.0.1"
)

var (
	Cmd    string
	CurDir string
)

func kaiab() error {
	switch Cmd {
	case "help":
		fmt.Println("show https://github.com/hinoshiba/kai-ab")
	case "init":
		return cmd_init(CurDir)
	case "template":
		if !isKaiabDir(CurDir) {
			return fmt.Errorf("Is not kai-ab directory, please 'kai-ab init'.")
		}
	case "autofil":
		if !isKaiabDir(CurDir) {
			return fmt.Errorf("Is not kai-ab directory, please 'kai-ab init'.")
		}
	case "mfil":
		if !isKaiabDir(CurDir) {
			return fmt.Errorf("Is not kai-ab directory, please 'kai-ab init'.")
		}
	case "check":
		if !isKaiabDir(CurDir) {
			return fmt.Errorf("Is not kai-ab directory, please 'kai-ab init'.")
		}
	case "calc":
		if !isKaiabDir(CurDir) {
			return fmt.Errorf("Is not kai-ab directory, please 'kai-ab init'.")
		}
	default:
		return fmt.Errorf("undefined operation: '%s'", Cmd)
	}

	return nil
}

func die(s string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, s + "\n" , msg...)
	os.Exit(1)
}

func init() {
	flag.Parse()

	cmd := flag.Arg(0)
	if cmd == "" {
		die("empty subcommand, Usage: kai-ab <sub-command>")
	}
	Cmd = cmd

	cur_dir, err := os.Getwd()
	if err != nil {
		die("cannot get current directory: %s", err)
	}
	CurDir = filepath.Clean(cur_dir)
}

func main() {
	if err := kaiab(); err != nil {
		die("%s", err)
	}
}
