package remarkable

import "testing"

func TestParseFramebufferSpyConfig(t *testing.T) {
	config, err := parseFramebufferSpyConfig("0xffff8c15b000,1404,1872,2,5616,0\n")
	if err != nil {
		t.Fatal(err)
	}
	if config.address != 0xffff8c15b000 || config.width != 1404 || config.height != 1872 || config.pixelType != 2 || config.bytesPerLine != 5616 || config.requiresReload {
		t.Fatalf("unexpected config: %#v", config)
	}
}

func TestParseFramebufferSpyConfigRejectsUnavailableFramebuffer(t *testing.T) {
	if _, err := parseFramebufferSpyConfig("NULL\n"); err == nil {
		t.Fatal("expected an error")
	}
}
