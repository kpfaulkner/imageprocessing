package imageprocessing

import (
	"fmt"
	"time"
)

func WhiteBalanceHSVImage(img *HSVImage, multiplier float64) error {

	// quick hack LUT.. turns out 7 times slower... will come back to this.
	//satLUT := make(map[float64]float64)
	l := len(img.Data)
	fmt.Printf("length %d\n", l)
	pix := img.Data
	start := time.Now()
	//var newSat float64
	//var ok bool
	for i := 0; i < l; i += 4 {

		/* LUT is 7 times slower :(
		curSat := pix[i+1]
		if newSat,ok = satLUT[curSat]; !ok {
			newSat = curSat * multiplier
			satLUT[curSat] = newSat
		}
		pix[i+1] = newSat */
		
		pix[i+1] *= multiplier
	}
	end := time.Now()
	fmt.Printf("saturation processed in %d ms\n", end.Sub(start).Milliseconds())

	return nil
}

