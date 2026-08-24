//go:build linux && arm64

package remarkable

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const accelerometerDevicesRoot = "/sys/bus/iio/devices"

var lastOrientation atomic.Int32
var accelerometerRoot string
var accelerometerOnce sync.Once

// CurrentOrientation returns the clockwise correction needed for an upright frame.
func CurrentOrientation() int {
	previous := int(lastOrientation.Load())
	accelerometerOnce.Do(func() {
		accelerometerRoot = detectAccelerometerRoot(accelerometerDevicesRoot)
	})
	if accelerometerRoot == "" {
		return previous
	}
	x, xErr := readAcceleration(filepath.Join(accelerometerRoot, "in_accel_x_raw"))
	y, yErr := readAcceleration(filepath.Join(accelerometerRoot, "in_accel_y_raw"))
	if xErr != nil || yErr != nil {
		return previous
	}
	orientation := orientationFromAcceleration(x, y, previous)
	lastOrientation.Store(int32(orientation))
	return orientation
}

func readAcceleration(path string) (int, error) {
	value, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(value)))
}
