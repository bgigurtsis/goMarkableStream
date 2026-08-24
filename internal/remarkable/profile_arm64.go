//go:build arm64

package remarkable

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	paperPureWidth  = 1404
	paperPureHeight = 1872
)

func init() {
	applyArm64DeviceProfile(detectArm64Device("/sys/devices/soc0/machine", "/etc/hostname"))
	PenInputDevice, TouchInputDevice = detectInputDevices("/sys/class/input", PenInputDevice, TouchInputDevice)
}

func detectArm64Device(machinePath, hostnamePath string) DeviceModel {
	identity := ""
	for _, path := range []string{machinePath, hostnamePath} {
		value, err := os.ReadFile(path)
		if err == nil {
			identity += " " + strings.ToLower(strings.TrimSpace(string(value)))
		}
	}
	if strings.Contains(identity, "paper pure") || strings.Contains(identity, "imx93-tatsu") || strings.Contains(identity, "rmppure") {
		return RemarkablePaperPure
	}
	return RemarkablePaperPro
}

func applyArm64DeviceProfile(model DeviceModel) {
	Model = model
	if model != RemarkablePaperPure {
		return
	}
	ScreenWidth = paperPureWidth
	ScreenHeight = paperPureHeight
	ScreenSizeBytes = paperPureWidth * paperPureHeight * BytesPerPixelBGRA
	Config = FramebufferConfig{
		Width:          paperPureWidth,
		Height:         paperPureHeight,
		StridePixels:   paperPureWidth,
		BytesPerPixel:  BytesPerPixelBGRA,
		SizeBytes:      ScreenSizeBytes,
		UseBGRA:        true,
		TextureFlipped: false,
	}
	log.Printf("Detected reMarkable Paper Pure profile (%dx%d BGRA)", paperPureWidth, paperPureHeight)
}

func detectInputDevices(root, defaultPen, defaultTouch string) (string, string) {
	pen, touch := defaultPen, defaultTouch
	events, _ := filepath.Glob(filepath.Join(root, "event*"))
	for _, event := range events {
		nameBytes, err := os.ReadFile(filepath.Join(event, "device", "name"))
		if err != nil {
			continue
		}
		name := strings.ToLower(string(nameBytes))
		device := filepath.Join("/dev/input", filepath.Base(event))
		if strings.Contains(name, "pen") || strings.Contains(name, "stylus") || strings.Contains(name, "wacom") || strings.Contains(name, "marker") {
			pen = device
		}
		if strings.Contains(name, "touch") || strings.Contains(name, "goodix") {
			touch = device
		}
	}
	return pen, touch
}
