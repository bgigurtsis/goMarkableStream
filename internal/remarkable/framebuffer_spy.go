package remarkable

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const framebufferSpyCommand = "printf '%s\\n' '>eframebuffer-spy$getConfigString:' > /run/xovi-mb; cat /run/xovi-mb-out"

type framebufferSpyConfig struct {
	address        int64
	width          int
	height         int
	pixelType      int
	bytesPerLine   int
	requiresReload bool
}

func getFramebufferSpyPointer() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, "/bin/sh", "-c", framebufferSpyCommand).Output()
	if ctx.Err() != nil {
		return 0, fmt.Errorf("request timed out; ensure framebuffer-spy and xovi-message-broker are installed and xochitl is running")
	}
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	spyConfig, err := parseFramebufferSpyConfig(string(output))
	if err != nil {
		return 0, err
	}
	if spyConfig.width != Config.Width || spyConfig.height != Config.Height {
		return 0, fmt.Errorf("unexpected framebuffer dimensions %dx%d", spyConfig.width, spyConfig.height)
	}
	if spyConfig.pixelType != 2 || spyConfig.bytesPerLine != Config.StridePixels*Config.BytesPerPixel {
		return 0, fmt.Errorf("unexpected framebuffer format type=%d bytes-per-line=%d", spyConfig.pixelType, spyConfig.bytesPerLine)
	}
	if spyConfig.requiresReload {
		return 0, fmt.Errorf("framebuffer requires reload, which Paper Pure does not support")
	}
	return spyConfig.address, nil
}

func parseFramebufferSpyConfig(value string) (framebufferSpyConfig, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "NULL" {
		return framebufferSpyConfig{}, fmt.Errorf("framebuffer-spy has not detected the Paper Pure display")
	}
	parts := strings.Split(value, ",")
	if len(parts) != 6 {
		return framebufferSpyConfig{}, fmt.Errorf("invalid framebuffer-spy response %q", value)
	}

	address, err := strconv.ParseUint(parts[0], 0, 64)
	if err != nil || address == 0 {
		return framebufferSpyConfig{}, fmt.Errorf("invalid framebuffer address %q", parts[0])
	}
	values := make([]int, 4)
	for index := range values {
		parsed, parseErr := strconv.Atoi(parts[index+1])
		if parseErr != nil {
			return framebufferSpyConfig{}, fmt.Errorf("invalid framebuffer-spy value %q", parts[index+1])
		}
		values[index] = parsed
	}
	requiresReload := false
	switch parts[5] {
	case "0", "false":
	case "1", "true":
		requiresReload = true
	default:
		return framebufferSpyConfig{}, fmt.Errorf("invalid reload value %q", parts[5])
	}

	return framebufferSpyConfig{
		address:        int64(address),
		width:          values[0],
		height:         values[1],
		pixelType:      values[2],
		bytesPerLine:   values[3],
		requiresReload: requiresReload,
	}, nil
}
