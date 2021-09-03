package main

import (
	"fmt"
	imageprocessing "github.com/kpfaulkner/imageprocessing/pkg"

	"log"
	"path/filepath"
	"time"
)

func main() {

	fmt.Printf("so it begins...\n")
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	filePath := "../images"
	fileName := "test-highres.jpg"
	fullFileName := filepath.Join(filePath, fileName)

	start := time.Now()
	hsvImage, err := imageprocessing.LoadImageAndConvertToHSV(fullFileName)
	if err != nil {
		log.Fatalf("BOOM unable to load image : %s", err.Error())
	}
	end := time.Now()
	fmt.Printf("load and convert took %d ms\n", end.Sub(start).Milliseconds())

	start = time.Now()
	imageprocessing.SaturateHSVImage(hsvImage, 0.6)
	//saturate(canvas, 0.6)
	end = time.Now()
	fmt.Printf("process took %d ms\n", end.Sub(start).Milliseconds())


	start = time.Now()
	_, err = imageprocessing.ConvertHSVToRGBImage(hsvImage)
	if err != nil {
		log.Fatalf("HSV to RGB BOOM %s", err.Error())
	}
	end = time.Now()
	fmt.Printf( "ConvertHSVToRGVImage took %d ms\n", end.Sub(start).Milliseconds())


}
