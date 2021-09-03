package imageprocessing

import (
	"fmt"
	"github.com/pkg/profile"
	"image"
	"time"
)

type HSVType struct {
	H float64
	S float64
	V float64
}

type RGBType struct {
	R uint8
	G uint8
	B uint8
}

var rgbToHSVLUT [256 * 256 * 256]HSVType

//var hsvToRGBLUT map[HSVType]RGBType
var hsvToRGBLUT map[float64]RGBType

func init() {
	// rgbToHSL

	start := time.Now()
	generateLUT()
	end := time.Now()
	fmt.Printf("generateLUT took %d ms\n", end.Sub(start).Milliseconds())

	//hsvToRGBLUT = make(map[HSVType]RGBType)
	hsvToRGBLUT = make(map[float64]RGBType)
}

func generateLUT() {

	var r uint16
	var g uint16
	var b uint16

	rgbToHSVLUT = [256 * 256 * 256]HSVType{}
	var index uint32
	// probably a bad idea :)
	for r = 0; r <= 255; r++ {
		for g = 0; g <= 255; g++ {
			for b = 0; b <= 255; b++ {
				index = uint32(r<<16 + g<<8 + b)
				//fmt.Printf("index is %d\n", index)
				h, s, v := rgbToHSV8(uint8(r), uint8(g), uint8(b))
				rgbToHSVLUT[index] = HSVType{H: h, S: s, V: v}
			}
		}
	}
}

func rgbToHSV8(r, g, b uint8) (float64, float64, float64) {
	rr := float64(r) / 255.0
	gg := float64(g) / 255.0
	bb := float64(b) / 255.0

	max := rr
	if bb > max {
		max = bb
	}
	if gg > max {
		max = gg
	}

	min := rr
	if bb < min {
		min = bb
	}
	if gg < min {
		min = gg
	}

	//max := math.Max(rr, math.Max(gg, bb))
	//min := math.Min(rr, math.Min(gg, bb))

	l := (max + min) / 2

	if max == min {
		return 0, 0, l
	}

	var h, s float64
	d := max - min
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case rr:
		h = (gg - bb) / d
		if g < b {
			h += 6
		}
	case gg:
		h = (bb-rr)/d + 2
	case bb:
		h = (rr-gg)/d + 4
	}
	h /= 6

	return h, s, l
}

// rgbToHSL converts a color from RGB to HSL.
func rgbToHSL(r, g, b uint32) (float64, float64, float64) {
	rr := float64(r) / 255.0
	gg := float64(g) / 255.0
	bb := float64(b) / 255.0

	max := rr
	if bb > max {
		max = bb
	}
	if gg > max {
		max = gg
	}

	min := rr
	if bb < min {
		min = bb
	}
	if gg < min {
		min = gg
	}

	//max := math.Max(rr, math.Max(gg, bb))
	//min := math.Min(rr, math.Min(gg, bb))

	l := (max + min) / 2

	if max == min {
		return 0, 0, l
	}

	var h, s float64
	d := max - min
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case rr:
		h = (gg - bb) / d
		if g < b {
			h += 6
		}
	case gg:
		h = (bb-rr)/d + 2
	case bb:
		h = (rr-gg)/d + 4
	}
	h /= 6

	return h, s, l
}

// clamp rounds and clamps float64 value to fit into uint8.
func clamp(x float64) uint8 {
	v := int64(x + 0.5)
	if v > 255 {
		return 255
	}
	if v > 0 {
		return uint8(v)
	}
	return 0
}

// clamp rounds and clamps float64 value to fit into uint8.
func clamp32(x float64) uint32 {
	v := int64(x + 0.5)
	if v > 255 {
		return 255
	}
	if v > 0 {
		return uint32(v)
	}
	return 0
}

// hslToRGB converts a color from HSL to RGB.
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	var r, g, b float64
	if s == 0 {
		v := clamp(l * 255)
		return v, v, v
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = hueToRGB(p, q, h+1/3.0)
	g = hueToRGB(p, q, h)
	b = hueToRGB(p, q, h-1/3.0)

	return clamp(r * 255), clamp(g * 255), clamp(b * 255)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1/2.0 {
		return q
	}
	if t < 2/3.0 {
		return p + (q-p)*(2/3.0-t)*6
	}
	return p
}

func ConvertRGBToHSVImage(img *image.NRGBA) (*HSVImage, error) {
	origSize := img.Bounds()
	hsvImage := NewHSVImage(origSize.Dx(), origSize.Dy())

	l := len(img.Pix)
	origPix := img.Pix
	hsvPix := hsvImage.Data

	for i := 0; i < l; i += 4 {
		index := uint32(origPix[i]<<16 + origPix[i+1]<<8 + origPix[i+2])
		hsv := rgbToHSVLUT[index]
		hsvPix[i] = hsv.H
		hsvPix[i+1] = hsv.S
		hsvPix[i+2] = hsv.V
		hsvPix[i+3] = float64(origPix[i+3]) // alpha.
	}

	return hsvImage, nil
}

func ConvertHSVToRGBImage(img *HSVImage) (*image.NRGBA, error) {
	p := profile.Start(profile.CPUProfile, profile.ProfilePath("."))

	rgbImg := image.NewNRGBA(image.Rect(0, 0, img.Width, img.Height))

	l := len(img.Data)
	origPix := img.Data
	rgbPix := rgbImg.Pix

	var rgb RGBType
	var ok bool
	cached := 0
	nonCache := 0
	for i := 0; i < l; i += 4 {
		//hsv := HSVType{H:origPix[i], S:origPix[i+1], V:origPix[i+2]}
		index := origPix[i]*1000000000 + origPix[i+1]*10000000 + origPix[i+2]
		//index := origPix[i] << 31 + origPix[i+1]*10000000 + origPix[i+2]*100000
		if rgb, ok = hsvToRGBLUT[index]; !ok {
			r, g, b := hslToRGB(origPix[i], origPix[i+1], origPix[i+2])
			rgb.R = r
			rgb.G = g
			rgb.B = b
			hsvToRGBLUT[index] = rgb
			nonCache++
		} else {
			cached++
		}

		rgbPix[i] = rgb.R
		rgbPix[i+1] = rgb.G
		rgbPix[i+2] = rgb.B
		rgbPix[i+3] = uint8(origPix[i+3]) // alpha.
	}

	fmt.Printf("ConvertHSVToRGBImage cached %d : noncached %d\n", cached, nonCache)
	p.Stop()
	return rgbImg, nil
}
