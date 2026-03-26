package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

// PreprocessResult คืนค่าภาพที่ผ่านการ preprocess แล้วพร้อม mimeType ที่อาจเปลี่ยนไป
type PreprocessResult struct {
	Data     []byte
	MIMEType string
}

// PreprocessIDCard ปรับภาพบัตรประชาชนให้ OCR อ่านได้ดีขึ้น
// คืนค่า PreprocessResult พร้อม mimeType ที่ถูกต้องหลัง encode
func PreprocessIDCard(imageBytes []byte, mimeType string) (*PreprocessResult, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("preprocess: อ่านภาพไม่ได้: %w", err)
	}

	bounds := img.Bounds()
	processed := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := img.At(x, y)
			r, g, b, a := px.RGBA()

			// Contrast enhancement factor 1.3
			const factor = 1.3
			r8 := uint8(clampF(float64(r>>8)*factor, 0, 255))
			g8 := uint8(clampF(float64(g>>8)*factor, 0, 255))
			b8 := uint8(clampF(float64(b>>8)*factor, 0, 255))

			processed.SetRGBA(x, y, color.RGBA{R: r8, G: g8, B: b8, A: uint8(a >> 8)})
		}
	}

	// เลือก encode format: JPEG สำหรับ jpg, PNG สำหรับทุกอย่างอื่น
	var buf bytes.Buffer
	outMIME := "image/png"
	switch mimeType {
	case "image/jpeg", "image/jpg":
		outMIME = "image/jpeg"
		err = jpeg.Encode(&buf, processed, &jpeg.Options{Quality: 92})
	default:
		// PNG, WebP, HEIC → encode เป็น PNG (Go standard library รองรับ)
		err = png.Encode(&buf, processed)
	}
	if err != nil {
		return nil, fmt.Errorf("preprocess: encode ไม่สำเร็จ: %w", err)
	}

	return &PreprocessResult{Data: buf.Bytes(), MIMEType: outMIME}, nil
}

// clampF จำกัดค่า float64 ให้อยู่ในช่วง [lo, hi]
func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
