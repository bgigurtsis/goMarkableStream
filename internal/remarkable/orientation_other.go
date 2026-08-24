//go:build !linux || !arm64

package remarkable

// CurrentOrientation returns no correction on devices without the Paper Pure sensor profile.
func CurrentOrientation() int {
	return 0
}
