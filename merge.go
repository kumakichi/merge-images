package merge_images

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/bmp"
)

var (
	backgroundColor color.Color
)

// MergeImage merges image files given by sourceImages argument one by one.
// The width/height of the merged image is the max width/height of image files in sourceImages.
func MergeImage(orientation Orientation, mergedImagePath string, sourceImagePaths ...string) (err error) {
	length := len(sourceImagePaths)
	widthSlice := make([]int, length)
	heightSlice := make([]int, length)

	for i, v := range sourceImagePaths {
		widthSlice[i], heightSlice[i], err = getImageDimension(v)
		if err != nil {
			return
		}
	}

	var x0, x1, y0, y1 int
	var dst *image.NRGBA

	if orientation == HORIZONTAL { // merge horizontally: y is uniform(lower is 0, higher is maxHeight)
		maxHeight := max(heightSlice...)
		totalWidth := sum(widthSlice...)
		dst = image.NewNRGBA(image.Rect(0, 0, totalWidth, maxHeight))

		if isBackgroundColorSetted() {
			for x := 0; x < totalWidth; x++ {
				for y := 0; y < maxHeight; y++ {
					dst.Set(x, y, backgroundColor)
				}
			}
		}

		for i, v := range sourceImagePaths {
			src, err := getImageFromPath(v)
			if err != nil {
				return err
			}

			if i == 0 {
				x0 = 0
				x1 = widthSlice[0]
			} else {
				x0 = sum(widthSlice[:i]...)
				x1 = x0 + widthSlice[i]
			}
			draw.Draw(dst, image.Rect(x0, 0, x1, maxHeight), src, src.Bounds().Min, draw.Src)
		}
	} else { // merge vertically: x is uniform(lower is 0, higher is maxWidth)
		maxWidth := max(widthSlice...)
		totalHeight := sum(heightSlice...)
		dst = image.NewNRGBA(image.Rect(0, 0, maxWidth, totalHeight))

		if isBackgroundColorSetted() {
			for x := 0; x < maxWidth; x++ {
				for y := 0; y < totalHeight; y++ {
					dst.Set(x, y, backgroundColor)
				}
			}
		}

		for i, v := range sourceImagePaths {
			src, err := getImageFromPath(v)
			if err != nil {
				return err
			}

			if i == 0 {
				y0 = 0
				y1 = heightSlice[0]
			} else {
				y0 = sum(heightSlice[:i]...)
				y1 = y0 + heightSlice[i]
			}
			draw.Draw(dst, image.Rect(0, y0, maxWidth, y1), src, src.Bounds().Min, draw.Src)
		}
	}

	outfile, err := os.Create(mergedImagePath)
	if err != nil {
		return err
	}

	err = jpeg.Encode(outfile, dst, nil)
	if err != nil {
		return
	}

	err = outfile.Close()
	return
}

// Set background color
func SetBackgroundColor(c color.Color) {
	backgroundColor = c
}

// Unset background color
func UnsetBackgroundColor(c color.Color) {
	backgroundColor = nil
}

func isBackgroundColorSetted() bool {
	return backgroundColor != nil
}

func getImageFromPath(imgPath string) (img image.Image, err error) {
	fp, err := os.Open(imgPath)
	if err != nil {
		return
	}

	switch getFormat(fp) {
	case IMG_JPG:
		img, err = jpeg.Decode(fp)
	case IMG_PNG:
		img, err = png.Decode(fp)
	case IMG_GIF:
		img, err = gif.Decode(fp)
	case IMG_BMP:
		img, err = bmp.Decode(fp)
	default:
		err = errors.New("unsupported image type")
		return
	}

	err = fp.Close()
	return
}

func getImageDimension(imgPath string) (width, height int, err error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return
	}

	conf, _, err := image.DecodeConfig(file)
	if err != nil {
		return
	}

	err = file.Close()
	if err != nil {
		return
	}

	width = conf.Width
	height = conf.Height
	return
}

func max(slice ...int) int {
	max := slice[0]

	for _, v := range slice {
		if v > max {
			max = v
		}
	}

	return max
}

func sum(i ...int) int {
	sum := 0
	for _, v := range i {
		sum += v
	}
	return sum
}

func getFormat(file *os.File) ImgType {
	bytes := make([]byte, 4)
	n, _ := file.ReadAt(bytes, 0)
	if n < 4 {
		return IMG_UNKNOWN
	}
	if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		return IMG_PNG
	}
	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		return IMG_JPG
	}
	if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
		return IMG_GIF
	}
	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		return IMG_BMP
	}
	return IMG_UNKNOWN
}
