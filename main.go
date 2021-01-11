package main

import (
	"fmt"
	"go-home-camera/localserver"
	"go-home-camera/picamera"
)

func main() {
	var picam picamera.PiCamera
	picamera.ApplyPreset("./settings/daily.json", &picam) //Initialize

	calib := make(chan string)

	savename := "./assets/image.jpg"

	go func() {
		for {
			select {
			case <-calib:
				fmt.Println("Start Calibration")
				picamera.Calibrate("night", "./settings/special.json")
			default:
				// applypresetをhandlerで書き換えるように変更する
				fmt.Printf("%+v\n", picam)
				picamera.ApplyPreset("./settings/special.json", &picam)
				picam.Capture(savename, true, true)
			}
		}
	}()

	localserver.ResponsiveServer(calib)
}
