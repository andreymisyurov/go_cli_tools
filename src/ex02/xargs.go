package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		os.Exit(1)
	}

	command := args[0]
	var flags []string
	if len(args) > 1 {
		flags = args[1:]
	}

	var arg_buffer []string
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeNamedPipe) != 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			arg_buffer = append(arg_buffer, scanner.Text())
		}
	}

	arg_buffer = append(flags, arg_buffer...)
	cmd := exec.Command(command, arg_buffer...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

}
