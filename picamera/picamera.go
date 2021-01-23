package picamera

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var exposureMargin int = 2
var cameraMode int = 4
var standardBrightness int = 50

// alternative implementation of constant slice
func exmodes() []string {
	return []string{"auto", "night", "nightpreview", "backlight", "spotlight", "sports", "snow", "beach", "verylong", "fixedfps", "antishake", "fireworks"}
}

// PiCamera : Configuration for ManualPiCamera
type PiCamera struct {
	Brightness  int `json:"brightness"`
	Exposure    int `json:"exposure"`
	AWBR        int `json:"awbr"`
	AWBB        int `json:"awbb"`
	AnalogGain  int `json:"analoggain"`
	DigitalGain int `json:"digitalgain"`
	Interval    int `json:"interval"`
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
func (picam *PiCamera) Capture(savename string, vflip bool, hflip bool) {
	var cmd []string

	// Constant: for manual capture
	cmd = append(cmd,
		"-n",
		"-awb", "off",
		"-ex", "off",
	)

	// Variables
	if vflip {
		cmd = append(cmd, "-vf")
	}

	if hflip {
		cmd = append(cmd, "-hf")
	}

	if picam.Interval*1000 < exposureMargin*picam.Exposure {
		picam.Interval = exposureMargin * picam.Exposure / 1000
	}

	cmd = append(cmd,
		"-ss", strconv.Itoa(picam.Exposure),
		"-md", strconv.Itoa(cameraMode),
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

func matchSlice(key string, array []string) bool {
	for _, v := range array {
		if key == v {
			return true
		}
	}
	return false
}

// alternative implementation of constant map
func parameterMap() map[string]int {
	return map[string]int{
		"Exposure":    8,
		"AnalogGain":  11,
		"DigitalGain": 14,
		"AWBR":        18,
		"AWBB":        20,
	}
}

func extractParameter(stdout []byte) (map[string]int, error) {
	commands := strings.Split(string(stdout), "\n")
	last := strings.Join(commands[len(commands)-3:len(commands)-1], " ")
	for _, v := range []string{",", "/256"} {
		last = strings.Replace(last, v, "", -1)
	}
	last = strings.Replace(last, "=", " ", -1)
	lastArr := strings.Split(last, " ")

	if len(lastArr) != 21 {
		log.Fatal("Could not calibrate")
		return nil, errors.New("Calibrate error")
	}

	m := make(map[string]int)
	for k, v := range parameterMap() {
		i, _ := strconv.Atoi(lastArr[v])
		m[k] = i
	}

	return m, nil
}

func mapToStruct(arr map[string]int) PiCamera {
	var cam PiCamera
	tmp, _ := json.Marshal(arr)
	json.Unmarshal(tmp, &cam)
	return cam
}

func saveJSON(cam *PiCamera, savename string) error {
	b, _ := json.MarshalIndent(cam, "", "\t")
	if ioutil.WriteFile(savename, b, 0755) != nil {
		log.Fatal("Could not save setting json")
		return errors.New("Could not save setting json")
	}
	return nil
}

// Calibrate : get calibrated camera setting
func Calibrate(ex string, savename string) {
	if !matchSlice(ex, exmodes()) {
		log.Fatalf("The exposure mode %s does not exit", ex)
		return
	}

	stdout, _ := exec.Command("raspistill", "-ex", ex, "-set", "-md", strconv.Itoa(cameraMode)).CombinedOutput()
	params, err := extractParameter(stdout)

	if err != nil {
		log.Fatal("Could not calibrate")
		return
	}

	params["Brightness"] = standardBrightness
	params["Interval"] = params["Exposure"] * exposureMargin / 1000

	cam := mapToStruct(params)

	if saveJSON(&cam, savename) != nil {
		log.Fatal("Could not save json file")
	}

	return
}
