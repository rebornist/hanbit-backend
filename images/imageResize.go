package images

import (
	"image"
	"os"

	"github.com/disintegration/imaging"
)

func imageResize(imagePath, imageType string) error {

	// open file
	src, err := imaging.Open(imagePath)
	if err != nil {
		return err
	}

	var divNum float64
	width, height, _ := imageScan(imagePath)
	if width >= height {
		divNum = float64(1980) / float64(width)
	} else {
		divNum = float64(1080) / float64(height)
	}
	width = int(float64(width) * divNum)

	src = imaging.Resize(src, width, 0, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err = imaging.Save(src, imagePath, imaging.JPEGQuality(90))
	if err != nil {
		return err
	}

	return nil
}

func imageScan(fp string) (width, height int, err error) {
	if reader, err := os.Open(fp); err == nil {
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			return 0, 0, err
		}
		return im.Width, im.Height, nil
	} else {
		return 0, 0, err
	}
}
