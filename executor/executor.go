package executor

import (
	"log"
	"os"
	"os/exec"
)

func ExecuteCommand(args []string) {
	name := args[0]
	args = args[1:]
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalf("error in cmd.Run: %v", err)
		}
	}
}
