package remarkable

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOrientationFromAcceleration(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		y        int
		previous int
		want     int
	}{
		{name: "portrait", x: -400, y: -9000, want: 0},
		{name: "inverted portrait", x: 500, y: 9000, want: 180},
		{name: "clockwise landscape", x: -9000, y: 300, want: 90},
		{name: "counterclockwise landscape", x: 9000, y: -300, want: 270},
		{name: "flat keeps last orientation", x: 100, y: -100, previous: 180, want: 180},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := orientationFromAcceleration(test.x, test.y, test.previous); got != test.want {
				t.Fatalf("orientationFromAcceleration() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestDetectAccelerometerRoot(t *testing.T) {
	root := t.TempDir()
	for device, name := range map[string]string{
		"iio:device0": "imx93-adc",
		"iio:device6": "lis2dw12_accel",
	} {
		path := filepath.Join(root, device)
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(path, "name"), []byte(name), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if got := detectAccelerometerRoot(root); got != filepath.Join(root, "iio:device6") {
		t.Fatalf("detectAccelerometerRoot() = %q", got)
	}
}
