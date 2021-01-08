package main

import (
	"go-home-camera/localserver"
	"go-home-camera/picamera"
)

func main() {
	// goroutine で camera撮影をする
	picam := picamera.PiCamera{
		Mode:           6,
		VerticalFlip:   true,
		HorizontalFlip: true,
		SavePath:       "./assets/image.jpg",
		Interval:       1000,
	}

	go picam.SequentialCapture()
	localserver.MyServer()
}
