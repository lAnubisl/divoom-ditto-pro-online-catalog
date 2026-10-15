package animation

import "errors"

// scrollFrames renders format 12's four tile-major RGB images through a moving
// 16x16 viewport, following the Android app's one-pixel scrolling renderer.
func scrollFrames(tiles []byte, mode byte) ([]byte, error) {
	if len(tiles) != 4*bytesPerFrame {
		return nil, errors.New("invalid scrolling pixel data")
	}
	if mode < 1 || mode > 3 {
		return nil, errors.New("unsupported scrolling mode")
	}
	tileCount := 4
	if mode != 3 {
		// The app retains all four tiles for an entirely black image.
		for tile := 3; tile >= 0; tile-- {
			nonBlack := false
			for _, value := range tiles[tile*bytesPerFrame : (tile+1)*bytesPerFrame] {
				if value != 0 {
					nonBlack = true
					break
				}
			}
			if nonBlack {
				tileCount = tile + 1
				break
			}
		}
	}
	columns := 1
	if mode == 1 {
		columns = tileCount
	} else if mode == 3 {
		columns = 2
	}
	width := columns * frameWidth
	height := tileCount / columns * frameHeight
	canvas := make([]byte, tileCount*bytesPerFrame)
	for tile := 0; tile < tileCount; tile++ {
		for y := 0; y < frameHeight; y++ {
			dst := ((tile/columns*frameHeight+y)*width + tile%columns*frameWidth) * rgbChannels
			src := tile*bytesPerFrame + y*frameWidth*rgbChannels
			copy(canvas[dst:dst+frameWidth*rgbChannels], tiles[src:src+frameWidth*rgbChannels])
		}
	}
	frames := make([]byte, tileCount*frameWidth*bytesPerFrame)
	for step := 0; step < tileCount*frameWidth; step++ {
		x, y := 0, 0
		switch mode {
		case 1:
			x = step
		case 2:
			y = step
		case 3:
			// Clockwise around the four corners of a 32x32 canvas.
			switch {
			case step < 16:
				x = step
			case step < 32:
				x, y = 16, step-16
			case step < 48:
				x, y = 48-step, 16
			default:
				y = 64 - step
			}
		}
		for py := 0; py < frameHeight; py++ {
			for px := 0; px < frameWidth; px++ {
				src := (((y+py)%height)*width + (x+px)%width) * rgbChannels
				dst := step*bytesPerFrame + (py*frameWidth+px)*rgbChannels
				copy(frames[dst:dst+rgbChannels], canvas[src:src+rgbChannels])
			}
		}
	}
	return frames, nil
}
