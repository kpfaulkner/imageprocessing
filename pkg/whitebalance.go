package imageprocessing

import (
	"image"
	"sort"
)

// based off https://adadevelopment.github.io/gdal/white-balance-gdal.html
func WhiteBalanceRGBImage(img *image.NRGBA, percentForBalance float64) error {

	redBytes := getBytes(0, img)

	sort.Slice(redBytes, func(i int, j int) bool {
		return redBytes[i]<redBytes[j]
	})

	percentForBalance := 0.6
	perc05 := Percentile(redBytes, percentForBalance)

	return nil
}

func Percentile( sortedBytes []uint8, percentile float64 ) []uint8 {

}

// probably get rid of this...  to inefficient
func getBytes(offset int, img *image.NRGBA) []uint8 {
	pix := img.Pix
	l := len(pix)
	colourBand := make([]uint8, l/4)
	idx := 0
	for i:=0; i<l; i += 4 {
		v := pix[i+offset]
		colourBand[idx] = v
		idx++
	}

	return colourBand
}
