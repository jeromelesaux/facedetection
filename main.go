package main

import (
	"facerecognition/facedetector"
	"fmt"
	"os"
)

// https://github.com/disintegration/imaging

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage: this image-to-detect-faces.png")
		return
	}
	//
	filename := os.Args[1]
	f := facedetector.NewFaceDetector(filename)
	f.DrawOnImage()

	//
	//var src image.Image
	//f, err := os.Open(filename)
	//if err != nil {
	//	fmt.Println("Error cannot open file " + filename + " with error : " + err.Error())
	//	panic(err.Error())
	//	return
	//}
	//
	//defer f.Close()
	//
	//point := strings.LastIndex(filename, ".")
	//extension := filename[point:]
	//switch {
	//case extension == ".png":
	//	src, err = png.Decode(f)
	//	if err != nil {
	//		fmt.Println("Error cannot decode file " + filename + " with error : " + err.Error())
	//		panic(err.Error())
	//		return
	//	}
	//case extension == ".jpg", extension == ".jpeg":
	//	src, err = jpeg.Decode(f)
	//	if err != nil {
	//		fmt.Println("Error cannot decode file " + filename + " with error : " + err.Error())
	//		panic(err.Error())
	//		return
	//	}
	//default:
	//	err := errors.New("not implemented.")
	//	panic(err)
	//	return
	//}
	//
	//src = imaging.AdjustContrast(src, 20)
	//src = imaging.AdjustBrightness(src, 40)
	//gray := grayscale.Convert(src, grayscale.ToGrayLuminance)
	//out, err := os.Create(filename + "-grayscale" + extension)
	//if err != nil {
	//	fmt.Println("Error cannot create file " + filename + "-grayscale" + extension + " with error : " + err.Error())
	//	panic(err.Error())
	//	return
	//}
	//defer out.Close()
	//switch {
	//case extension == ".png":
	//	png.Encode(out, gray)
	//case extension == ".jpeg", extension == ".jpg":
	//	jpeg.Encode(out, gray, &jpeg.Options{jpeg.DefaultQuality})
	//default:
	//	err := errors.New("not implemented.")
	//	panic(err)
	//	return
	//}
	return
}
