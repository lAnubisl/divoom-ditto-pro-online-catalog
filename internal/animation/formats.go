package animation

import "errors"

// ErrUnsupportedFormat distinguishes unsupported encodings from damaged files
// and transient download failures.
var ErrUnsupportedFormat = errors.New("unsupported artwork format (expected 16x16 format 8, 9, or 12)")

// Binary artwork encodings, distinct from gallery API FileType values.
const (
	FormatStillImage byte = 8
	FormatAnimation  byte = 9
	FormatScroll     byte = 12
)

// SupportsFormat reports whether the decoder implements this artwork encoding.
func SupportsFormat(format byte) bool {
	switch format {
	case FormatStillImage, FormatAnimation, FormatScroll:
		return true
	default:
		return false
	}
}
