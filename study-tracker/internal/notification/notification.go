package notification

import (
	"os/exec"
)

func Send(title, message string) {
	cmd := exec.Command("notify-send", title, message)
	cmd.Run()
}
