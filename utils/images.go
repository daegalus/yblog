package utils

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"sync"

	"github.com/caarlos0/log"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/jpegxl"
	"github.com/gen2brain/webp"
	"github.com/spf13/afero"
)

//TODO: AVIF/JXL have washed out colors from WEBP. Might be from the YcBcr to RGBA conversion that is happening.

func ImageToAVIF(original image.Image) []byte {
	encodedImage := []byte{}
	buf := bytes.NewBuffer(encodedImage)
	options := avif.Options{
		Quality:           50,
		QualityAlpha:      50,
		Speed:             4,
		ChromaSubsampling: image.YCbCrSubsampleRatio420,
	}

	b := original.Bounds()
	m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), original, b.Min, draw.Src)
	err := avif.Encode(buf, m, options)

	if err != nil {
		log.WithError(err).Error("Failed to encode AVIF")
		return nil
	}
	return buf.Bytes()
}

func ImageToAVIFThread(original image.Image, wg *sync.WaitGroup, returnChan chan []byte) {
	defer wg.Done()
	encodedImage := []byte{}
	buf := bytes.NewBuffer(encodedImage)
	options := avif.Options{
		Quality:           50,
		QualityAlpha:      50,
		Speed:             4,
		ChromaSubsampling: image.YCbCrSubsampleRatio420,
	}

	b := original.Bounds()
	m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), original, b.Min, draw.Src)
	err := avif.Encode(buf, m, options)

	if err != nil {
		log.WithError(err).Error("Failed to encode AVIF")
	}

	returnChan <- buf.Bytes()
}

func ImageToJXL(original image.Image) []byte {
	encodedImage := []byte{}
	buf := bytes.NewBuffer(encodedImage)
	options := jpegxl.Options{
		Quality: 75,
		Effort:  7,
	}
	err := jpegxl.Encode(buf, original, options)
	if err != nil {
		log.WithError(err).Error("Failed to encode JXL")
		return nil
	}
	return buf.Bytes()
}

func ImageToWebP(original image.Image) []byte {
	encodedImage := []byte{}
	buf := bytes.NewBuffer(encodedImage)
	options := webp.Options{
		Quality: 75,
		Method:  6,
	}
	err := webp.Encode(buf, original, options)
	if err != nil {
		log.WithError(err).Error("Failed to encode WebP")
		return nil
	}
	return buf.Bytes()
}

func ImageToWebPThreaded(original image.Image, wg *sync.WaitGroup, returnChan chan []byte) {
	defer wg.Done()
	encodedImage := []byte{}
	buf := bytes.NewBuffer(encodedImage)
	options := webp.Options{
		Quality: 75,
		Method:  6,
	}
	err := webp.Encode(buf, original, options)
	if err != nil {
		log.WithError(err).Error("Failed to encode WebP")
	}
	returnChan <- buf.Bytes()
}

func ImageFromWebP(fs afero.Fs, filepath string) image.Image {
	webpImage, err := afero.ReadFile(fs, filepath)
	if err != nil {
		log.WithError(err).Error("Failed to read WebP file")
		return nil
	}

	imageData, err := webp.Decode(bytes.NewBuffer(webpImage))
	if err != nil {
		log.WithError(err).Error("Failed to decode WebP")
		return nil
	}

	return imageData
}

func ImageFromPNG(fs afero.Fs, filepath string) image.Image {
	pngImage, err := afero.ReadFile(fs, filepath)
	if err != nil {
		log.WithError(err).Error("Failed to read PNG file")
		return nil
	}

	imageData, err := png.Decode(bytes.NewBuffer(pngImage))
	if err != nil {
		log.WithError(err).Error("Failed to decode PNG")
		return nil
	}

	return imageData
}
