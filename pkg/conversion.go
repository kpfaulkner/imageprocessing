package imageprocessing

// RGBHSVConversion used for converting images (or really byte arrays) between
// RGB and HSV
type RGBHSVConversion struct {

}

// ConvertRGBToHSV converts image
func ConvertRGBToHSV( img CustomImage) error {

	return nil
}

func GenerateValueHistogram(img CustomImage) *Histogram {
	binCount := 256
	v := Histogram{make([]int, binCount)}

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			pos := y*img.Stride + x*4
			v.Bins[img.Data[pos+0]]++
		}
	}

	return &v
}

// Histogram holds a variable length slice of bins, which keeps track of sample counts.
type Histogram struct {
	Bins []int
}

// RGBAHistogram holds a sub-histogram per RGBA channel.
// Each channel histogram contains 256 bins (8-bit color depth per channel).
type RGBAHistogram struct {
	R Histogram
	G Histogram
	B Histogram
	A Histogram
}

