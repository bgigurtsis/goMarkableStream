package remarkable

import (
	"os"
	"path/filepath"
	"strings"
)

const orientationMotionThreshold = 2048

func orientationFromAcceleration(x, y, previous int) int {
	if abs(x) < orientationMotionThreshold && abs(y) < orientationMotionThreshold {
		return previous
	}
	if abs(y) >= abs(x) {
		if y < 0 {
			return 0
		}
		return 180
	}
	if x < 0 {
		return 90
	}
	return 270
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func detectAccelerometerRoot(root string) string {
	devices, _ := filepath.Glob(filepath.Join(root, "iio:device*"))
	for _, device := range devices {
		name, err := os.ReadFile(filepath.Join(device, "name"))
		if err == nil && strings.TrimSpace(string(name)) == "lis2dw12_accel" {
			return device
		}
	}
	return ""
}
