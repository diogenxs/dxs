package utils

import (
	"os"
	"os/exec"
)

func WriteToClipboard(message string) {
	c1 := exec.Command("echo", message)
	c2 := exec.Command("xclip", "-i", "-selection", "clipboard")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = os.Stdout
	_ = c2.Start()
	_ = c1.Run()
	_ = c2.Wait()
}
