package picamera

import (
	"log"
	"os/exec"
	"strconv"
)

// PiCamera : Configuration for PiCamera
type PiCamera struct {
	Mode           int
	VerticalFlip   bool
	HorizontalFlip bool
	Interval       int
	SavePath       string
}

// Capture : Capture 1 frame image
func (picam *PiCamera) Capture() {
	var cmd []string

	if picam.VerticalFlip {
		cmd = append(cmd, "-vf")
	}

	if picam.HorizontalFlip {
		cmd = append(cmd, "-hf")
	}

	cmd = append(cmd,
		"-md", strconv.Itoa(picam.Mode),
		"-o", picam.SavePath,
		"-t", strconv.Itoa(picam.Interval),
	)

	exec.Command("raspistill", cmd...).Run()
}

// SequentialCapture : Capture frames by a constant interval
func (picam *PiCamera) SequentialCapture() {
	if picam.Interval <= 0 {
		log.Fatalln("Interval must be positive")
	}

	for {
		picam.Capture()
	}
}
