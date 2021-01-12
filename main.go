package main

import (
	"fmt"
	"go-home-camera/localserver"
	"go-home-camera/picamera"
)

func main() {
	// 結局プリセットをほとんど活用しないのでリファクタリングすること
	var picam picamera.PiCamera
	picamera.ApplyPreset("./settings/daily.json", &picam) //Initialize

	calib := make(chan string)

	savename := "./assets/image.jpg"

	go func() {
		for {
			select {
			case exmode := <-calib:
				fmt.Printf("Start Calibration: mode = %s\n", exmode)
				picamera.Calibrate(exmode, "./settings/special.json")
			default:
				picamera.ApplyPreset("./settings/special.json", &picam)
				picam.Capture(savename, true, true)
			}
		}
	}()

	localserver.ResponsiveServer(calib)
}
