//go:build arm64

package remarkable

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectArm64PaperPure(t *testing.T) {
	dir := t.TempDir()
	machine := filepath.Join(dir, "machine")
	hostname := filepath.Join(dir, "hostname")
	if err := os.WriteFile(machine, []byte("reMarkable Paper Pure\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hostname, []byte("imx93-tatsu\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := detectArm64Device(machine, hostname); got != RemarkablePaperPure {
		t.Fatalf("detectArm64Device() = %s", got)
	}
}

func TestDetectPaperPureInputDevices(t *testing.T) {
	root := t.TempDir()
	for event, name := range map[string]string{"event2": "Elan marker input", "event3": "Elan touch input"} {
		dir := filepath.Join(root, event, "device")
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "name"), []byte(name), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	pen, touch := detectInputDevices(root, "/dev/input/event2", "/dev/input/event3")
	if pen != "/dev/input/event2" || touch != "/dev/input/event3" {
		t.Fatalf("unexpected input devices: pen=%s touch=%s", pen, touch)
	}
}

func TestPaperPureUsesLiveBGRAFramebufferShape(t *testing.T) {
	previousModel := Model
	previousWidth := ScreenWidth
	previousHeight := ScreenHeight
	previousSize := ScreenSizeBytes
	previousConfig := Config
	applyArm64DeviceProfile(RemarkablePaperPure)
	t.Cleanup(func() {
		Model = previousModel
		ScreenWidth = previousWidth
		ScreenHeight = previousHeight
		ScreenSizeBytes = previousSize
		Config = previousConfig
	})

	if Config.Width != 1404 || Config.Height != 1872 || Config.StridePixels != 1404 {
		t.Fatalf("unexpected dimensions: %#v", Config)
	}
	if !Config.UseBGRA || Config.BytesPerPixel != 4 || Config.SizeBytes != 1404*1872*4 {
		t.Fatalf("unexpected framebuffer format: %#v", Config)
	}
}
