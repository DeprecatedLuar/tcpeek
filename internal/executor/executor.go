package executor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Run(command string) error {
	log.Printf("[INFO] Executing: %s", command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Dir, _ = os.UserHomeDir()

	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] Command failed: %s (%v)", command, err)
		notifyError(command, err)
		return err
	}
	return nil
}

func notifyError(command string, err error) {
	msg := fmt.Sprintf("%s\n(%v)", command, err)
	exec.Command("notify-send", "-u", "critical", "tcpeek: command failed", msg).Run()
}
