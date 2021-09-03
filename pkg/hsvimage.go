package imageprocessing

import (
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"log"
	"os"
)

type ImageType byte

const (
	RGB ImageType = iota
	HSV
)

type HSVImage struct {
	Data []float64

	// Width and Height are pixels
	Width  int
	Height int

	// Stride is number of bytes across (4 bytes per pixel).. so its Width*4
	Stride    int
	ImageType ImageType
}

// will keep alpha in this image for later conversion back to RGB
func NewHSVImage(width int, height int) *HSVImage {
	img := HSVImage{}
	img.Width = width
	img.Height = height
	img.Data = make([]float64, width*height*4) // H,S,V,A <--- keeping alpha here for later conversion back to RGBA
	img.ImageType = HSV
	img.Stride = 4 // 4 bytes per pixel.
	return &img
}

func LoadImageAndConvertToHSV(filePath string) (*HSVImage, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatalf("BOOM %s", err.Error())
	}

	origSize := img.Bounds()
	canvas := image.NewNRGBA(origSize)

	draw.Draw(canvas, origSize, img, image.Point{0, 0}, draw.Src)
	//draw.Copy( canvas, image.Point{},img,origSize, draw.Src,nil)

	hsvImage, err := ConvertRGBToHSVImage(canvas)
	return hsvImage, err
}
