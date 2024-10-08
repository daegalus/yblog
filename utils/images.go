package utils

import (
	"bytes"
	"image"
	"os"

	"github.com/gen2brain/avif"
	"github.com/gen2brain/jpegxl"
	"github.com/gen2brain/webp"
)

//TODO: AVIF/JXL have washed out colors from WEBP. Might be from the YcBcr to RGBA conversion that is happening.

func ImageToAVIF(original image.Image) ([]byte, error) {
	encodedImage := []byte{}
	buf := bytes.NewBuffer(encodedImage)
	options := avif.Options{
		Quality:           50,
		QualityAlpha:      50,
		Speed:             4,
		ChromaSubsampling: image.YCbCrSubsampleRatio420,
	}
	err := avif.Encode(buf, original, options)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ImageToJXL(original image.Image) ([]byte, error) {
	encodedImage := []byte{}
	buf := bytes.NewBuffer(encodedImage)
	options := jpegxl.Options{
		Quality: 75,
		Effort:  7,
	}
	err := jpegxl.Encode(buf, original, options)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ImageToWebP(original image.Image) ([]byte, error) {
	encodedImage := []byte{}
	buf := bytes.NewBuffer(encodedImage)
	options := webp.Options{
		Quality: 75,
		Method:  6,
	}
	err := webp.Encode(buf, original, options)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ImageFromWebP(filepath string) (image.Image, error) {
	webpImage, err := os.ReadFile(filepath)
	buf := bytes.NewBuffer(webpImage)
	if err != nil {
		return nil, err
	}

	imageData, err := webp.Decode(buf)
	if err != nil {
		return nil, err
	}

	return imageData, err
}
