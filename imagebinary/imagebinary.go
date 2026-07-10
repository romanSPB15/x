package imagebinary

import (
	"encoding/binary"
	"fmt"
	"image"
	"math"
)

// Encode упаковывает байты в RGBA-изображение.
// Первые 4 байта хранят длину данных (big-endian).
func Encode(data []byte) *image.RGBA {
	pixels := float64(len(data)+4) / 4
	side := int(math.Ceil(math.Sqrt(pixels)))
	img := image.NewRGBA(image.Rect(0, 0, side, side))
	binary.BigEndian.PutUint32(img.Pix, uint32(len(data)))
	copy(img.Pix[4:], data)
	return img
}

// Decode извлекает байты из изображения.
// Поддерживает *image.RGBA и *image.NRGBA.
func Decode(img image.Image) ([]byte, error) {
	var pix []byte
	switch v := img.(type) {
	case *image.RGBA:
		pix = v.Pix
	case *image.NRGBA:
		pix = v.Pix
	default:
		return nil, fmt.Errorf("unsupported image type: %T", img)
	}
	if len(pix) < 4 {
		return nil, fmt.Errorf("image too small")
	}
	length := binary.BigEndian.Uint32(pix[:4])
	if length > uint32(len(pix)-4) {
		length = uint32(len(pix) - 4)
	}
	result := make([]byte, length)
	copy(result, pix[4:4+length])
	return result, nil
}
