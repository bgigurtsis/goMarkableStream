package stream

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/owulveryck/goMarkableStream/internal/remarkable"
)

// NewScreenshotHandler creates a new screenshot handler reading from file @pointerAddr
func NewScreenshotHandler(file io.ReaderAt, pointerAddr int64) *ScreenshotHandler {
	return &ScreenshotHandler{
		file:        file,
		pointerAddr: pointerAddr,
	}
}

// ScreenshotHandler is an http.Handler that serves PNG screenshots of the framebuffer
type ScreenshotHandler struct {
	file        io.ReaderAt
	pointerAddr int64
}

// ServeHTTP implements http.Handler
func (h *ScreenshotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imageDataPtr := rawFrameBuffer.Get().(*[]uint8)
	imageData := *imageDataPtr
	defer rawFrameBuffer.Put(imageDataPtr)

	_, err := h.file.ReadAt(imageData, h.pointerAddr)
	if err != nil {
		log.Printf("failed to read framebuffer: %v", err)
		http.Error(w, "failed to read framebuffer", http.StatusInternalServerError)
		return
	}

	width := remarkable.Config.Width
	height := remarkable.Config.Height

	// Create RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	stride := remarkable.Config.StridePixels
	if stride == 0 {
		stride = width
	}

	// Convert the runtime framebuffer format to RGBA.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dstIdx := (y*width + x) * 4
			if remarkable.Config.UseBGRA {
				srcIdx := (y*stride + x) * 4
				img.Pix[dstIdx+0] = imageData[srcIdx+2]
				img.Pix[dstIdx+1] = imageData[srcIdx+1]
				img.Pix[dstIdx+2] = imageData[srcIdx+0]
			} else {
				srcIdx := (y*stride + x) * 2
				pixel := uint16(imageData[srcIdx]) | uint16(imageData[srcIdx+1])<<8
				img.Pix[dstIdx+0] = uint8(((pixel >> 11) & 0x1f) * 255 / 31)
				img.Pix[dstIdx+1] = uint8(((pixel >> 5) & 0x3f) * 255 / 63)
				img.Pix[dstIdx+2] = uint8((pixel & 0x1f) * 255 / 31)
			}
			img.Pix[dstIdx+3] = 255 // A (fully opaque)
		}
	}

	// Generate filename with timestamp
	filename := "remarkable_" + time.Now().Format("20060102_150405") + ".png"

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Remarkable-Orientation", fmt.Sprintf("%d", remarkable.CurrentOrientation()))

	if err := png.Encode(w, img); err != nil {
		log.Printf("failed to encode PNG: %v", err)
	}
}
