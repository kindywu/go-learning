package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "docker-init")

	//ps aux | grep docker-inity || true

	var outputBuf1 bytes.Buffer
	cmd1.Stdout = &outputBuf1
	if err := cmd1.Start(); err != nil {
		fmt.Printf("Error: the first command can't be startup %s\n", err)
		return
	}
	if err := cmd1.Wait(); err != nil {
		fmt.Printf("Error: the first command can't be done %s\n", err)
		return
	}
	// println(outputBuf1.String())

	cmd2.Stdin = &outputBuf1
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error: the second command can't be startup %s\n", err)
		return
	}
	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Error: the second command can't be done %s\n", err)
		return
	}

	println(outputBuf2.String())
}
