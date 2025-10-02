package notification

import (
	"os/exec"
)

func Send(title, message string) {
	cmd := exec.Command("notify-send", title, message)
	cmd.Run()
}

func PlaySound() {
	// You can change this path to your preferred sound file.
	soundFilePath := "/usr/share/sounds/freedesktop/stereo/complete.oga"
	cmd := exec.Command("paplay", soundFilePath)
	cmd.Run()
}
