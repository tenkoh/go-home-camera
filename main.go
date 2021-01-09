package main

import (
	"fmt"
	"go-home-camera/picamera"
)

var picam picamera.PiCamera
var savename string

func main() {
	// goroutine で camera撮影をする
	savename = "./assets/image.jpg"
	picamera.ApplyPreset("./settings/daytime_no_light.json", &picam)
	picam.Capture(savename)

	fmt.Printf("%+v\n", picam)
}
