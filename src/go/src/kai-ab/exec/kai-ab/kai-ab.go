package main

import (
	"os"
	"fmt"
	"flag"
	"path/filepath"
)

var (
	Cmd    string
	CurDir string
)

func kaiab() error {
	switch Cmd {
	case "help":
		fmt.Println(STR_HELP)
		return nil
	case "init":
		return cmd_init(CurDir)
	}

	if !IsKaiabDir(CurDir) {
		return fmt.Errorf("Is not kai-ab directory, please 'kai-ab init'.")
	}
	switch Cmd {
	case "template":
		date := flag.Arg(1)
		return cmd_template(date)

	case "autofil":
		path := flag.Arg(1)
		return cmd_autofilter(path)
	case "mfil":
		//path := flag.Arg(1)
		fmt.Println("havn't function, yet")
	case "check":
		fmt.Println("havn't function, yet")
	case "calc":
		return cmd_calc()
	default:
		return fmt.Errorf("undefined operation: '%s'\nshow 'kai-ab help'", Cmd)
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
