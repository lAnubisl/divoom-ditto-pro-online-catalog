// Package media defines size limits shared by downloads, decoding, and uploads.
package media

const (
	// MaxArtworkSize is the maximum encrypted gallery artwork size in bytes.
	MaxArtworkSize = 2 * 1024 * 1024 // 2 MiB
	// MaxAnimationSize is the device API's maximum GIF upload size in bytes.
	MaxAnimationSize = 2 * 1024 * 1024 // 2 MiB
)
