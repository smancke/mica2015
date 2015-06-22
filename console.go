package main

import (
    "os"
    "os/exec"
)

func waitForAnyKey() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
}
