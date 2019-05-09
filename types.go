package merge_images

// image type
type ImgType int8

const (
	IMG_PNG ImgType = iota
	IMG_BMP
	IMG_JPG
	IMG_GIF
	IMG_UNKNOWN
)

// the orientation to merge images
type Orientation int8

const (
	HORIZONTAL Orientation = iota
	VERTICLE
)
