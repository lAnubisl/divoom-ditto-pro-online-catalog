# Divoom artwork formats 8, 9 and 12

This document covers all artwork formats currently accepted by `internal/animation`: still images (8), animations (9), and scrolling artwork (12).

Formats 8 and 12 were investigated on 2026-10-14 using the extracted Android app 3.8.40 and actual files returned by Christmas album 34 through `GetAlbumImageListV3`. Format 9 was already supported; the real reindeer sample below was checked again during that investigation.

## Shared encoding

All three formats are proprietary RGB containers, rather than PNG/JPEG/GIF files. Dispatch uses the downloaded file's first byte. The header is unencrypted; the remaining payload uses AES-128-CBC with key `78hrey23y28ogs89` and IV `1234567890123456`. The Android helper uses `AES/CBC/NoPadding`. No compression library or Android native code is required by the client.

Each 16x16 frame or tile contains 768 bytes: 256 row-major pixels, with three bytes per pixel in red, green, blue order. Pixels are opaque.

| Format | Unencrypted header | Decrypted payload | Client output |
| --- | --- | --- | --- |
| 8 (`0x08`) | One byte: format | Exactly 768 RGB bytes | One-frame 16x16 GIF |
| 9 (`0x09`) | Four bytes: format, metadata byte, two-byte big-endian frame delay in milliseconds | Consecutive 768-byte RGB frames, with limited trailing data tolerated | Multi-frame 16x16 GIF, or a still GIF if there is only one frame |
| 12 (`0x0c`) | Four bytes: format, scroll mode, two-byte big-endian step delay in milliseconds | Exactly 3,072 bytes: four consecutive 16x16 RGB tiles | 16x16 GIF generated from a moving viewport |

Gallery API `FileType`, the app's internal `PixelBean.type`, and the downloaded format byte are different namespaces. For example, the format-12 samples have API `FileType: 8`, but decode to internal `TYPE_SCROLL` (6), not `TYPE_PLANET_ANI` (8). API metadata must not be used as the format-dispatch byte.

## Format 8: still image

The complete file is exactly 769 bytes: one format byte followed by 768 encrypted RGB bytes. Decryption yields one 16x16 image. The GIF encoder assigns it a 10 ms delay; presentation time is controlled by playback, rather than that frame delay.

## Format 9: animation

The decoder skips the four-byte header, decrypts the remaining data, and splits it into 768-byte frames. The second header byte is not used or validated by the client; the checked 48-frame sample stores `0x30` there. The client derives the frame count from the decrypted payload length instead.

Bytes 2 and 3 store the frame delay in milliseconds. The GIF delay is rounded to the nearest centisecond using `(milliseconds + 5) / 10`, with a minimum of one centisecond (10 ms). The reindeer sample's header is `09 30 00 77`: 119 ms becomes 120 ms per GIF frame.

Between 1 and 1,024 complete frames are accepted. The encrypted payload must be nonempty and AES-block aligned. The encoder tolerates a remainder of at most 16 bytes after the complete frames and ignores it; it does not interpret or validate those bytes as PKCS padding.

## Format 12: scrolling artwork

Format 12 is not four ordinary animation frames. The app assembles four tiles into a larger canvas, then moves a 16x16 viewport across it.

| Mode | Canvas | Viewport movement |
| --- | --- | --- |
| 1 | Up to 64x16, tiles arranged horizontally | One pixel horizontally per frame, wrapping at the canvas edge |
| 2 | Up to 16x64, tiles arranged vertically | One pixel vertically per frame, wrapping at the canvas edge |
| 3 | 32x32, tiles arranged in a 2x2 grid | Clockwise around the four corners, one pixel per frame |

Modes 1 and 2 trim trailing completely black tiles after the last nonblack tile. An entirely black image retains all four tiles. Mode 3 always retains all four tiles. The output contains 16 frames per retained tile, up to 64 frames for a full loop.

The normal recovered Android renderer advances one pixel per frame and wraps. Its separate conversion routine advances two pixels and doubles the frame delay. The client follows the one-pixel renderer and uses the header delay with the same GIF rounding as format 9.

The enum also contains mode 0 (`ONE_ONE`), but the recovered scrolling renderer produces no frames for it. The client rejects mode 0 and all values outside 1?3. The complete file must be exactly 3,076 bytes.

## GIF output, caching and playback

All supported inputs are converted to GIF; raw PNG, JPEG and GIF downloads are not accepted by this decoder. Each frame has its own palette, preserving every RGB color in a 256-pixel image without quantization. Frames use no transparency, `DisposalNone`, and a loop count of zero (infinite looping for animations).

The source adapter stores converted GIFs permanently as `<GalleryId>.gif` and reuses cached files for preview and playback. Browser previews return `image/gif`. Device uploads also use `image/gif`: a one-frame GIF goes to `POST /api/images`, while multiple frames go to `POST /api/animations`. Playback duration and clock restoration remain controlled by the playback controller and rotation runner.

Decoding and automatic selection have different boundaries. `gallery.Item.Selectable` requires a positive ID, a nonempty file path, and size 0 or 1. It accepts API type 1 for categories and albums, plus type 0 for albums. API type 8 scroll entries are therefore excluded from automatic selection even though their format-12 files can be decoded for previews and manual sends.

Empty input, unsupported format bytes, truncated headers, invalid payload lengths, unsupported scroll modes, and invalid frame counts are rejected. Both encrypted artwork and generated GIFs have a 2 MiB size limit, defined in `internal/media/limits.go`.

## Actual files checked

| Gallery ID | Name | API type / size | File bytes | Header | Decoded result |
| --- | --- | --- | --- | --- | --- |
| 4203677 | Christmas Wonderland 2025+ Shoshkie | 0 / 1 | 769 | `08` | One frame; snowman image, 7 colors; 128-byte GIF |
| 4223443 | Christmas Wonderland2025+REINDEER | 1 / 1 | 36,868 | `09 30 00 77` | 48 frames; 10,511-byte GIF |
| 4212595 | Christmaswonderland2025+DarkNight | 8 / 1 | 3,076 | `0c 01 00 64` | Horizontal winter scene, 17 colors; 64 frames; 7,930-byte GIF |
| 4204635 | Christmas Wonderland 2025 + Christmas Icons | 8 / 1 | 3,076 | `0c 02 00 64` | Vertical Christmas symbols, 4 colors; 64 frames; 4,819-byte GIF |

Both scroll samples store a 100 ms step delay. Static decoded pixels were rendered and visually inspected. During the original format-8/12 verification, all three new samples returned HTTP 200 with `image/gif` through the local service's authenticated preview routes. The existing format-9 sample still produced a byte-identical 48-frame GIF. `go build ./cmd/server` and `go vet ./...` passed at that time; no tests were written and no artwork was sent to the physical display.

Scrolling playback timing has not been compared against the running mobile app or a physical display. Mode 3 follows the recovered algorithm but remains unverified against a real mode-3 artwork file.

Local diagnostic artifacts are under ignored `bin/`: `format-<ID>.artwork`, decrypted `.artwork.rgb`, enlarged tile-strip `.artwork.png`, converted `.gif`, and `format-research-summary.json`. The format-9 sample is stored as `christmas-reindeer-4223443.artwork` and `.gif`.
