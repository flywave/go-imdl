package imdl

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	webp "github.com/chai2010/webp"
)

func init() {
	image.RegisterFormat("webp", "RIFF????WEBPVP8", webp.Decode, webp.DecodeConfig)
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
}

type TextureFormat uint32

const (
	FormatJPG  TextureFormat = 0
	FormatPNG  TextureFormat = 1
	FormatWEBP TextureFormat = 2
)

func encodeImage(format TextureFormat, writer io.Writer, rgba image.Image) {
	if format == FormatJPG {
		jpeg.Encode(writer, rgba, nil)
	} else if format == FormatPNG {
		png.Encode(writer, rgba)
	} else if format == FormatWEBP {
		webp.Encode(writer, rgba, &webp.Options{Lossless: true})
	}
}

func decodeImage(format TextureFormat, reader io.Reader) image.Image {
	if format == FormatJPG {
		img, err := jpeg.Decode(reader)
		if err != nil {
			return nil
		}
		return img
	} else if format == FormatPNG {
		img, err := png.Decode(reader)
		if err != nil {
			return nil
		}
		return img
	} else if format == FormatWEBP {
		img, err := webp.Decode(reader)
		if err != nil {
			return nil
		}
		return img
	}
	return nil
}

func EncodeTexture(texture image.Image, format TextureFormat) []byte {
	writer := &bytes.Buffer{}
	encodeImage(format, writer, texture)
	return writer.Bytes()
}

func DecodeTexture(data []byte, format TextureFormat) image.Image {
	return decodeImage(format, bytes.NewBuffer(data))
}
