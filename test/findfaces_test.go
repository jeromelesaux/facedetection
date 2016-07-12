package test

import (
	"facedetection/facedetector"
	"fmt"
	"testing"
)

func TestFindCarellFace(t *testing.T) {
	f2 := facedetector.NewFaceDetector("carell.png")
	rects := f2.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f2.DrawOnImage()
}

func TestFindWomanFace(t *testing.T) {
	f := facedetector.NewFaceDetector("test.png")
	rects := f.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f.DrawOnImage()

}

func TestMultipleFaces(t *testing.T) {
	f := facedetector.NewFaceDetector("trainingset.png")
	rects := f.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f.DrawOnImage()
}

func TestGeorgeFace(t *testing.T) {
	f := facedetector.NewFaceDetector("george-clooney.png")
	rects := f.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f.DrawOnImage()
}
