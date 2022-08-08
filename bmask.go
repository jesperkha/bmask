package bmask

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// Todo: add anti-aliasing

// Draw creates a canvas and draws a border of specified with around the source image.
// It then uses the canvas as a mask for the background image. The threshold is what
// determines if a pixel should be counted when checking for an edge. Default to 0.
func Draw(source string, background string, lineWidth int, threshold int) error {
	// Source image, the one a border is draw around and used as a mask
	img, err := openImage(source)
	if err != nil {
		return err
	}

	// The background image, used as a background for the final mask
	backgroundImg, err := openImage(background)
	if err != nil {
		return err
	}

	// Create new image canvas
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw border
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			neighbors := []uint32{
				alphaAt(img, x-1, y),
				alphaAt(img, x+1, y),
				alphaAt(img, x, y+1),
				alphaAt(img, x, y-1),
			}

			sum := uint32(0)
			ctr := 0
			for _, n := range neighbors {
				sum += n
				if n > uint32(threshold) {
					ctr++
				}
			}

			// No neighbors means its outside the image, 4 neighbors mean its inside.
			// Anything else is a border
			if ctr > 0 && ctr < 4 {
				for bx := x - lineWidth/2; bx < x+lineWidth/2; bx++ {
					for by := y - lineWidth/2; by < y+lineWidth/2; by++ {
						if bx < 0 || bx >= width || by < 0 || by >= height {
							continue
						}

						canvas.Set(bx, by, color.Black)
					}
				}
			}
		}
	}

	// Use canvas as mask over backround image
	// Todo: center background image
	outputImg := image.NewRGBA(image.Rect(0, 0, width, height))
	zero := image.Pt(0, 0)
	draw.DrawMask(outputImg, img.Bounds(), backgroundImg, zero, canvas, zero, draw.Over)

	return saveImage("out.png", outputImg)
}

// Todo: add support for jpg
func openImage(filename string) (img image.Image, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return img, err
	}

	defer file.Close()
	img, _, err = image.Decode(file)
	return img, err
}

func saveImage(filename string, img image.Image) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()
	return png.Encode(file, img)
}

func alphaAt(source image.Image, x, y int) uint32 {
	_, _, _, a := source.At(x, y).RGBA()
	return a
}
