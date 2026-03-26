package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

// PreprocessIDCard ปรับภาพบัตรประชาชนให้ OCR อ่านได้ดีขึ้น (basic version)
// 1. Contrast enhancement
// 2. Brightness normalization
// 3. Sharpness enhancement
func PreprocessIDCard(imageBytes []byte, mimeType string) ([]byte, error) {
	// ใช้ fallback implementation ชั่วคราวจนกว่าจะมี OpenCV
	return PreprocessIDCardFallback(imageBytes, mimeType)
}

// PreprocessIDCardFallback ถ้า OpenCV ใช้ไม่ได้ จะทำ basic contrast enhancement แทน
func PreprocessIDCardFallback(imageBytes []byte, mimeType string) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("fallback: อ่านภาพไม่ได้: %w", err)
	}

	bounds := img.Bounds()
	processed := image.NewRGBA(bounds)

	// Simple contrast enhancement
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()

			// คูณด้วย factor 1.2 เพื่อเพิ่ม contrast
			factor := 1.2
			r8 := uint8(min(float64(r>>8)*factor, 255))
			g8 := uint8(min(float64(g>>8)*factor, 255))
			b8 := uint8(min(float64(b>>8)*factor, 255))

			processed.SetRGBA(x, y, color.RGBA{R: r8, G: g8, B: b8, A: uint8(a >> 8)})
		}
	}

	var buf bytes.Buffer
	switch mimeType {
	case "image/png":
		err = png.Encode(&buf, processed)
	case "image/jpeg":
		err = jpeg.Encode(&buf, processed, &jpeg.Options{Quality: 90})
	default:
		err = png.Encode(&buf, processed)
	}

	return buf.Bytes(), err
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
