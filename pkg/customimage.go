package imageprocessing

type ImageType byte

const (
	RGB ImageType = iota
	HSV
)

// CustomImage just an uint8 array
// Quicker to deal with than Go's built in
// Image types etc
type CustomImage struct {
	Data []uint8

	// Width and Height are pixels
	Width int
	Height int

	// Stride is number of bytes across (4 bytes per pixel).. so its Width*4
	Stride int
	ImageType ImageType
}

func NewCustomImage(filename string) (CustomImage, error) {

}
