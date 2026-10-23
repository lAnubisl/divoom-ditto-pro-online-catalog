package animation

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/gif"

	"divoom-gallery-client/internal/media"
)

const (
	headerSize                 = 4
	delayOffset                = 2
	frameWidth                 = 16
	frameHeight                = 16
	rgbChannels                = 3
	pixelsPerFrame             = frameWidth * frameHeight
	bytesPerFrame              = pixelsPerFrame * rgbChannels
	minFrameCount              = 1
	maxFrameCount              = 1024
	millisecondsPerCentisecond = 10
	minGIFDelay                = 1 // Centiseconds.
	gifLoopForever             = 0
	opaqueAlpha                = 255
	artworkAESKey              = "78hrey23y28ogs89"
	artworkAESIV               = "1234567890123456"
)

// Decode converts Divoom still (8), animation (9), and scrolling (12) artwork
// into a 16x16 GIF. See docs/artwork-formats.md for the recovered layouts.
func Decode(data []byte) ([]byte, error) {
	if len(data) > media.MaxArtworkSize {
		return nil, errors.New("artwork exceeds size limit")
	}
	if len(data) == 0 {
		return nil, errors.New("empty artwork")
	}
	var raw []byte
	delay := minGIFDelay
	var err error
	switch data[0] {
	case FormatStillImage:
		if len(data) != 1+bytesPerFrame {
			return nil, errors.New("invalid still image length")
		}
		raw, err = decrypt(data[1:])
	case FormatAnimation, FormatScroll:
		if len(data) < headerSize {
			return nil, errors.New("truncated artwork header")
		}
		if data[0] == FormatScroll && len(data) != headerSize+4*bytesPerFrame {
			return nil, errors.New("invalid scrolling image length")
		}
		raw, err = decrypt(data[headerSize:])
		// Round milliseconds to the nearest GIF centisecond.
		delay = (int(binary.BigEndian.Uint16(data[delayOffset:headerSize])) + millisecondsPerCentisecond/2) / millisecondsPerCentisecond
		if delay < minGIFDelay {
			delay = minGIFDelay
		}
		if err == nil && data[0] == FormatScroll {
			raw, err = scrollFrames(raw, data[1])
		}
	default:
		return nil, ErrUnsupportedFormat
	}
	if err != nil {
		return nil, err
	}
	return encodeGIF(raw, delay)
}

func decrypt(encrypted []byte) ([]byte, error) {
	if len(encrypted) == 0 || len(encrypted)%aes.BlockSize != 0 {
		return nil, errors.New("invalid encrypted frame length")
	}
	block, err := aes.NewCipher([]byte(artworkAESKey))
	if err != nil {
		return nil, err
	}
	raw := make([]byte, len(encrypted))
	cipher.NewCBCDecrypter(block, []byte(artworkAESIV)).CryptBlocks(raw, encrypted)
	return raw, nil
}

func encodeGIF(raw []byte, delay int) ([]byte, error) {
	frameCount := len(raw) / bytesPerFrame
	if frameCount < minFrameCount || frameCount > maxFrameCount {
		return nil, errors.New("invalid animation frame count")
	}
	// APIxoo ignores trailing padding. Permit only a partial AES padding block.
	if len(raw)%bytesPerFrame > aes.BlockSize {
		return nil, errors.New("invalid frame data remainder")
	}
	animation := &gif.GIF{LoopCount: gifLoopForever}
	for f := 0; f < frameCount; f++ {
		palette := color.Palette{}
		indices := map[[rgbChannels]byte]uint8{}
		pixels := make([]byte, pixelsPerFrame)
		for p := 0; p < pixelsPerFrame; p++ {
			offset := f*bytesPerFrame + p*rgbChannels
			rgb := [rgbChannels]byte{raw[offset], raw[offset+1], raw[offset+2]}
			idx, exists := indices[rgb]
			if !exists {
				idx = uint8(len(palette))
				indices[rgb] = idx
				palette = append(palette, color.RGBA{rgb[0], rgb[1], rgb[2], opaqueAlpha})
			}
			pixels[p] = idx
		}
		frame := image.NewPaletted(image.Rect(0, 0, frameWidth, frameHeight), palette)
		copy(frame.Pix, pixels)
		animation.Image = append(animation.Image, frame)
		animation.Delay = append(animation.Delay, delay)
		animation.Disposal = append(animation.Disposal, gif.DisposalNone)
	}
	var output bytes.Buffer
	if err := gif.EncodeAll(&output, animation); err != nil {
		return nil, errors.New("GIF encoding failed")
	}
	if output.Len() > media.MaxAnimationSize {
		return nil, errors.New("GIF exceeds device upload limit")
	}
	return output.Bytes(), nil
}
