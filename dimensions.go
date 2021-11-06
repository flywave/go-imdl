package imdl

import "math"

const maxTextureSize = 1024

type Dimensions struct {
	Width  uint32
	Height uint32
}

func ComputeDimensions(nEntries uint32, nRgbaPerEntry uint32, nExtraRgba uint32) *Dimensions {
	maxSize := maxTextureSize
	nRgba := nEntries*nRgbaPerEntry + nExtraRgba

	if nRgba < uint32(maxSize) {
		return &Dimensions{Width: nRgba, Height: 1}
	}

	width := uint32(math.Ceil(math.Sqrt(float64(nRgba))))

	remainder := width % nRgbaPerEntry
	if remainder != 0 {
		width += nRgbaPerEntry - remainder
	}

	height := uint32(math.Ceil(float64(nRgba) / float64(width)))

	return &Dimensions{width, height}
}
