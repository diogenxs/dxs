package utils

import (
	"os/exec"
	"strconv"
)

// Notification is a notification
type Notification struct {
	Title    string
	Message  string
	Duration int
}

// Notify create a new notification
func (n *Notification) Notify() error {
	send, err := exec.LookPath("kdialog")
	if err != nil {
		return err
	}

	c := exec.Command(send, "--title", n.Title, "--passivepopup", n.Message, strconv.Itoa(n.Duration))
	err = c.Run()
	if err != nil {
		return err
	}
	return nil
}

// NewNotification creates a new notification
func NewNotification(title, message string) error {
	n := &Notification{Title: title, Message: message, Duration: 10}
	return n.Notify()
}
