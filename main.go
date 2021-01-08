package main

import (
	"fmt"
	"go-home-camera/picamera"
)

var picam picamera.PiCamera

func main() {
	// goroutine で camera撮影をする
	picamera.ApplyPreset("./settings/daily.json", &picam)
	picam.Capture()

	fmt.Printf("%+v\n", picam)
}
