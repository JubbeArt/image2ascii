package i2a

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"

	"github.com/aybabtme/rgbterm"
	"github.com/nfnt/resize"
)

type Options struct {
	Width  int
	Height int
}

func Image2String(image image.Image, options Options) string {
	scaledImage := scaleImage(image, options)
	scaledSize := scaledImage.Bounds()
	width := scaledSize.Dx()
	height := scaledSize.Dy()

	var buffer bytes.Buffer

	for i := 0; i < int(height); i++ {
		for j := 0; j < int(width); j++ {
			pixel := color.NRGBAModel.Convert(scaledImage.At(j, i))
			rawChar := pixelToASCII(pixel)
			buffer.WriteString(rawChar)
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// scaleImage resize the convert to expected size base on the convert options
func scaleImage(image image.Image, options Options) image.Image {
	img := image.Bounds()
	imgRatio := float64(img.Dx()) / float64(img.Dy())
	width, height := options.Width, options.Height

	if width == 0 && height == 0 {
		width = img.Dx()
		height = img.Dy()
	} else if width == 0 {
		width = int(float64(options.Height) * imgRatio)
	} else if height == 0 {
		height = int(float64(options.Width) / imgRatio)
	}

	charWidth := 0.5
	return resize.Resize(uint(width), uint(float64(height)*charWidth), image, resize.Lanczos3)
}

var pixels = []byte(" .,:;i1tfLCG08@")

// pixelToASCII converts a pixel to a ASCII char string
func pixelToASCII(pixel color.Color) string {
	r, g, b, a := pixel.RGBA()
	intensity := float64((r + g + b) * (a / 65535))
	maxIntensity := float64(65535 * 3)
	indexF := intensity / maxIntensity * float64(len(pixels)-1)
	index := int(math.Round(indexF))
	char := pixels[index]
	return rgbterm.FgString(string([]byte{char}), uint8(r), uint8(g), uint8(b))
}
