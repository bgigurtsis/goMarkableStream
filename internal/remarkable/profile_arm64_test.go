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
	for event, name := range map[string]string{"event4": "Wacom Pen", "event7": "Goodix Touchscreen"} {
		dir := filepath.Join(root, event, "device")
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "name"), []byte(name), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	pen, touch := detectInputDevices(root, "/dev/input/event2", "/dev/input/event3")
	if pen != "/dev/input/event4" || touch != "/dev/input/event7" {
		t.Fatalf("unexpected input devices: pen=%s touch=%s", pen, touch)
	}
}
