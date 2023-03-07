package test

import (
	"fmt"
	"testing"

	"github.com/jeromelesaux/facedetection/facedetector"
)

func TestFindCarellFace(t *testing.T) {
	f2 := facedetector.NewFaceDetector("carell.png", "haarcascade_frontalface_default.xml")
	rects := f2.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f2.DrawOnImage()
}

func TestFindWomanFace(t *testing.T) {
	f := facedetector.NewFaceDetector("test.png", "haarcascade_frontalface_default.xml")
	rects := f.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f.DrawOnImage()

}

func TestMultipleFaces(t *testing.T) {
	f := facedetector.NewFaceDetector("trainingset.png", "haarcascade_frontalface_default.xml")
	rects := f.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f.DrawOnImage()
}

func TestGeorgeFace(t *testing.T) {
	f := facedetector.NewFaceDetector("test.png", "haarcascade_frontalface_default.xml")
	rects := f.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f.DrawOnImage()
}

func TestObamaFace(t *testing.T) {
	f := facedetector.NewFaceDetector("obama.jpg", "haarcascade_frontalface_default.xml")
	rects := f.GetFaces()
	for _, r := range rects {
		fmt.Println(r.ToString())
	}
	f.DrawOnImage()
}
