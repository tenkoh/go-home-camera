package picamera

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
)

var exposureMargin int = 2

// PiCamera : Configuration for ManualPiCamera
type PiCamera struct {
	Mode           int  `json:"mode"`
	Brightness     int  `json:"brightness"`
	Exposure       int  `json:"exposure"`
	AWBR           int  `json:"awbr"`
	AWBB           int  `json:"awbb"`
	AnalogGain     int  `json:"analoggain"`
	DigitalGain    int  `json:"digitalgain"`
	Interval       int  `json:"interval"`
	VerticalFlip   bool `json:"verticalflip"`
	HorizontalFlip bool `json:"horizontalflip"`
}

// ApplyPreset : Overwrite PiCamera setting by preset file
func ApplyPreset(file string, picam *PiCamera) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Can not find the setting file: %s", file)
	}

	json.Unmarshal(buf, picam)
	return
}

// Capture : See https://www.raspberrypi.org/documentation/raspbian/applications/camera.md
func (picam *PiCamera) Capture(savename string) {
	var cmd []string

	// Constant: for manual capture
	cmd = append(cmd,
		"-n",
		"-awb", "off",
		"-ex", "off",
	)

	// Variables
	if picam.VerticalFlip {
		cmd = append(cmd, "-vf")
	}

	if picam.HorizontalFlip {
		cmd = append(cmd, "-hf")
	}

	if picam.Interval*1000 < exposureMargin*picam.Exposure {
		picam.Interval = exposureMargin * picam.Exposure / 1000
	}

	cmd = append(cmd,
		"-ss", strconv.Itoa(picam.Exposure),
		"-md", strconv.Itoa(picam.Mode),
		"-o", savename,
		"-br", strconv.Itoa(picam.Brightness),
		"-awbg", fmt.Sprintf("%.2f,%.2f", float32(picam.AWBR)/256.0, float32(picam.AWBB)/256.0),
		"-ag", fmt.Sprintf("%.2f", float32(picam.AnalogGain)/256.0),
		"-dg", fmt.Sprintf("%.2f", float32(picam.DigitalGain)/256.0),
		"-t", strconv.Itoa(picam.Interval),
	)

	exec.Command("raspistill", cmd...).Run()
	// fmt.Println(cmd)
}
