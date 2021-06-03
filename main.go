package main

import (
	"context"
	"os/exec"
	"time"
)

func main() {
	const textColor = "FFFFFFCC"
	lock := exec.Command("i3lock",
		"-n", "-B5", "-k", "--pass-media-keys", "--ignore-empty-password",
		"--indicator",
		"--time-color="+textColor, "--date-color="+textColor)

	if err := lock.Start(); err != nil {
		println("Error starting command:", err.Error())

		return
	}

	turnOn := make(chan bool, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		select {
		case <-time.After(1 * time.Minute):
			if err := exec.Command("xset", "dpms", "force", "off").Run(); err != nil {
				println("Error on dmps off:", err.Error())
			}

			return
		case <-ctx.Done():
			return
		}
	}(ctx)

	go func(turnOn chan bool) {
		lampStatus := exec.Command("yeelight-cli", "get", "power")
		status, err := lampStatus.Output()

		on := err == nil && string(status) == "on\n"
		turnOn <- on

		if on {
			if err := exec.Command("yeelight-cli", "off").Run(); err != nil {
				println("Error turning off lamp:", err.Error())
			}
		}
	}(turnOn)

	if err := lock.Wait(); err != nil {
		println("Lock command error:", err.Error())

		return
	}

	if <-turnOn {
		if err := exec.Command("yeelight-cli", "on").Run(); err != nil {
			println("Error turning on lamp:", err.Error())
		}
	}
}
