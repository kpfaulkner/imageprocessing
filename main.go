package main

import (
	"fmt"
	"github.com/anthonynsimon/bild/histogram"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"time"
	"github.com/pkg/profile"
)

type HSVType struct {
  H float64
  S float64
  V float64
}
var rgbToHSVLUT [256*256*256]HSVType

func init() {
  // rgbToHSL

	start := time.Now()
	generateLUT()
	end := time.Now()
	fmt.Printf("generateLUT took %d ms\n", end.Sub(start).Milliseconds())
}

func generateLUT() {

	var r uint16
	var g uint16
	var b uint16

	rgbToHSVLUT = [256*256*256]HSVType{}
  var index uint32
	// probably a bad idea :)
	for r = 0; r<=255;r++ {
		for g = 0; g<=255;g++ {
			for b = 0; b<=255;b++ {
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
		return v,v,v
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

func saturate(img *image.NRGBA, multiplier float64) error {

	c := 0


	start := time.Now()


	var r uint8
	var g uint8
	var b uint8

	/*
	bounds := img.Bounds()
	  dx := bounds.Dx()
	  dy := bounds.Dy()

		pix2 := img.Pix
	var x int
	var y int
  // usually takes 360ms +
	//col := color.RGBA{R: 100, G: 0, B: 0, A: 255}
	for x = 0; x < dx; x++ {
		for y = 0; y < dy; y++ {
			ydx := y*dx
			r := pix2[x+ydx]
      g := pix2[x+1+ydx]
      b := pix2[x+2+ydx]
      h,s,l := rgbToHSV8(r,g,b)
      //s *= multiplier
      r,g,b = hslToRGB(h,s,l)
			pix2[x+ydx] = r
			pix2[x+1+ydx] = g
			pix2[x+2+ydx] = b
      c++
		}
	}



  end := time.Now()
  fmt.Printf("first processed %d bytes in %d ms\n", c, end.Sub(start).Milliseconds())
*/
	start = time.Now()

	// BELOW IS WAY FASTER THAT ABOVE....  35ms-ish
	l := len(img.Pix)
	fmt.Printf("length %d\n", l)
	c = 0
	pix := img.Pix
	for i := 0; i < l; i += 4 {
	  r = pix[i]
    g = pix[i+1]
    b = pix[i+2]
    //h,s,ll := rgbToHSL8(r,g,b)
		//h,s,ll := rgbToHSL8(pix[i], pix[i+1], pix[i+2])
		//_,_,_ = rgbToHSV8(pix[i], pix[i+1], pix[i+2])
		index :=uint32( pix[i] << 16 + pix[i+1] << 8 + pix[i+2] )
		_= rgbToHSVLUT[index]
    //r,g,b = hslToRGB(hsv.H, hsv.S, hsv.V)
    pix[i] = r
    pix[i+1] = g
    pix[i+2] = b
    c++
	}
	fmt.Printf("second pix %d\n", img.Pix[0])

	end := time.Now()
	fmt.Printf("second processed %d bytes in %d ms\n", c, end.Sub(start).Milliseconds())

	return nil
}

func main() {

	fmt.Printf("so it begins...\n")
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	filePath := "../images"
	fileName := "test-highres.jpg"
	fullFileName := filepath.Join(filePath, fileName)

	start := time.Now()
	f, err := os.Open(fullFileName)
	if err != nil {
		log.Fatalf("BOOM %s", err.Error())
	}
	defer f.Close()
	end := time.Now()
	fmt.Printf("load image took %d ms\n", end.Sub(start).Milliseconds())

	start = time.Now()
	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatalf("BOOM %s", err.Error())
	}

	end = time.Now()
	fmt.Printf("decode image took %d ms\n", end.Sub(start).Milliseconds())


	origSize := img.Bounds()
	canvas := image.NewNRGBA(origSize)

	start = time.Now()
	draw.Draw(canvas, origSize, img, image.Point{0, 0}, draw.Src)
	end = time.Now()
	fmt.Printf("initial draw took %d ms\n", end.Sub(start).Milliseconds())

	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	p := profile.Start(profile.CPUProfile, profile.ProfilePath("."))

	start = time.Now()
	saturate(canvas, 0.6)
	end = time.Now()
	fmt.Printf("process took %d ms\n", end.Sub(start).Milliseconds())

	p.Stop()

	hist := histogram.NewRGBAHistogram(img)
	_ = hist.Image()
	//fmt.Printf("res is %s\n", result)
}
