package facedetector

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var configXml sync.Once
var Config *ConfigHaarcascade = &ConfigHaarcascade{}

type FoundRect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (f *FoundRect) ToString() string {
	return "X:" + strconv.Itoa(f.X) + ",Y:" + strconv.Itoa(f.Y) + ",width:" + strconv.Itoa(f.Width) + ",height:" + strconv.Itoa(f.Height)
}

type ConfigHaarcascade struct {
	Stages []*Stage
}

type FaceDetector struct {
	ClassifiedSize []int
	Width          int
	Height         int
	FoundRects     []*FoundRect
	Image          image.Image
}

type OpenCVFile struct {
	OpenCVStorage *OpenCVStorage `xml:"opencv_storage"`
}

type OpenCVStorage struct {
	Haarcascade *HaarcascadeFrontalfaceDefault `xml:"haarcascade_frontalface_default"`
}

type HaarcascadeFrontalfaceDefault struct {
	Stages *Stages `xml:"stages"`
}

type Stages struct {
	Stage []*Stage `xml:"_"`
}

type Stage struct {
	Trees     *Trees  `xml:"trees"`
	Threshold float64 `xml:"stage_threshold"`
}

type Trees struct {
	Trees []*Tree `xml:"_"`
}

type Tree struct {
	RootNode *RootNode `xml:"_"`
}

type Feature struct {
	Tilted float64 `xml:"tiltded"`
	Rects  []*Rect `xml:"rects"`
}

type RectValue struct {
	X1     int
	Y1     int
	X2     int
	Y2     int
	Weight float64
}

type Rect struct {
	V     []string `xml:"_"`
	Rects []*RectValue
}

func (r *Rect) Parse() {
	var err error
	//fmt.Println(r.V)
	for _, in := range r.V {
		rect := &RectValue{}
		s := strings.Split(in, " ")
		rect.X1, err = strconv.Atoi(s[0])
		if err != nil {
			fmt.Println(err.Error())
		}
		rect.X2, err = strconv.Atoi(s[1])
		if err != nil {
			fmt.Println(err.Error())
		}
		rect.Y1, err = strconv.Atoi(s[2])
		if err != nil {
			fmt.Println(err.Error())
		}
		rect.Y2, err = strconv.Atoi(s[3])
		if err != nil {
			fmt.Println(err.Error())
		}
		rect.Weight, err = strconv.ParseFloat(s[4], 64)
		if err != nil {
			fmt.Println(err.Error())
		}
		r.Rects = append(r.Rects, rect)
	}
	//fmt.Println(r)
}

func (s *Stages) Parse() {
	for _, s := range s.Stage {
		s.Parse()
	}
}

func (s *Stage) Parse() {
	for _, t := range s.Trees.Trees {
		for _, r := range t.RootNode.Feature.Rects {
			r.Parse()
		}
	}
}

func (s *Stage) Pass(grayImage [][]float64, square [][]float64, i int, j int, scale float64) bool {

	sum := 0.
	for _, tree := range s.Trees.Trees {
		sum += tree.RootNode.GetVal(grayImage, square, i, j, scale)
	}

	return sum > s.Threshold
}

func (f *Feature) Add(rect *Rect) {
	f.Rects = append(f.Rects, rect)
}

func (r *RootNode) GetVal(grayImage [][]float64, square [][]float64, i int, j int, scale float64) float64 {
	w := int(24 * scale)
	h := int(24 * scale)
	var invArea = 1. / float64(w*h)

	totalx := grayImage[i+w][j+h] + grayImage[i][j] - grayImage[i][j+h] - grayImage[i+w][j]
	totalx2 := square[i+w][j+h] + square[i][j] - square[i][j+h] - square[i+w][j]

	//fmt.Printf("totalx %f totalx2 %f invarea:%f\n", totalx, totalx2, invArea)
	moy := totalx * invArea
	var vnorm = totalx2*invArea - moy*moy

	if vnorm > 1.0 {
		vnorm = math.Sqrt(vnorm)
	} else {
		vnorm = 1.0
	}
	var rectSum = 0.0

	f := r.Feature
	for _, rect := range f.Rects {
		for _, rectV := range rect.Rects {
			//fmt.Printf("%f %d %d %d %d %d\n", r.Threshold, scale, rect.X1, rect.Y1, rect.X2, rect.Y2)
			rx1 := i + int(scale*float64(rectV.X1))
			rx2 := i + int(scale*float64(rectV.X1+rectV.Y1))
			ry1 := j + int(scale*float64(rectV.X2))
			ry2 := j + int(scale*float64(rectV.X2+rectV.Y2))
			//fmt.Printf("%d %d %d %d\n", rectV.X1, rectV.X1, rectV.Y1, rectV.Y2)
			//fmt.Printf("%d %d %d %d\n", rx1, rx2, ry1, ry2)
			rectSum += (grayImage[rx2][ry2] - grayImage[rx1][ry2] - grayImage[rx2][ry1] + grayImage[rx1][ry1]) * rectV.Weight
		}
	}

	rectSum2 := rectSum * invArea
	//fmt.Printf("%f %f %f %f %f %f\n", rectSum2, r.Threshold, vnorm, rectSum, invArea, moy)
	if rectSum2 < (r.Threshold * vnorm) {
		return r.LeftVal
	} else {
		return r.RightVal
	}

}

type RootNode struct {
	Feature   *Feature `xml:"feature"`
	Threshold float64  `xml:"threshold"`
	LeftVal   float64  `xml:"left_val"`
	RightVal  float64  `xml:"right_val"`
}

func (f *FaceDetector) Equals(r1 *FoundRect, r2 *FoundRect) bool {
	distance := float64(r1.Width) * 0.2
	if float64(r2.X) <= float64(r1.X)+distance && float64(r2.X) >= float64(r1.X)-distance && float64(r2.Y) <= float64(r1.Y)+distance && float64(r2.Y) >= float64(r1.Y)-distance && float64(r2.Width) <= (float64(r1.Width)*1.2) && (float64(r2.Width)*1.2) >= float64(r1.Width) {
		return true
	}
	if r1.X >= r2.X && r1.X+r1.Width <= r2.X+r2.Width && r1.Y >= r2.Y && r1.Y+r1.Height <= r2.Y+r2.Height {
		return true
	}
	return false

}

func (f *FaceDetector) GetFaces() []*FoundRect {
	fmt.Println("GetFaces.")
	return f.merge(f.FoundRects, 1)
}

func (f *FaceDetector) merge(rects []*FoundRect, minNeighbors int64) []*FoundRect {
	retour := make([]*FoundRect, 0)
	ret := make([]int, len(rects))
	nbClasses := 0
	for i := 0; i < len(rects); i++ {
		found := false
		for j := 0; j < i; j++ {
			if f.Equals(rects[j], rects[i]) {
				found = true
				ret[i] = ret[j]
			}
		}
		if !found {
			ret[i] = nbClasses
			nbClasses++
		}
	}
	neighbors := make([]int, len(rects))
	rect := make([]*FoundRect, 0)
	for i := 0; i < len(rects); i++ {
		neighbors[i] = 0
		rect = append(rect, &FoundRect{X: 0, Y: 0, Width: 0, Height: 0})
	}
	for i := 0; i < len(rects); i++ {
		neighbors[ret[i]]++
		rect[ret[i]].X += rects[i].X
		rect[ret[i]].Y += rects[i].Y
		rect[ret[i]].Width += rects[i].Width
		rect[ret[i]].Height += rects[i].Height
	}

	for i := 0; i < nbClasses; i++ {
		n := neighbors[i]
		if int64(n) >= minNeighbors {
			r := &FoundRect{X: 0, Y: 0, Width: 0, Height: 0}
			r.X = (rect[i].X*2 + n) / (2 * n)
			r.Y = (rect[i].Y*2 + n) / (2 * n)
			r.Width = (rect[i].Width*2 + n) / (2 * n)
			r.Height = (rect[i].Height*2 + n) / (2 * n)
			retour = append(retour, r)
		}
	}
	return retour

}

func (face *FaceDetector) DrawImageInDirectory(directory string) []string {
	filesPath := make([]string, 0)
	for i, r := range face.GetFaces() {
		b := make([]byte, 16)
		rand.Read(b)
		id := fmt.Sprintf("%X", b)
		dstRect := image.Rect(r.X, r.Y, (r.X + r.Width), (r.Y + r.Height))
		//fmt.Println(dstRect)
		dst := image.NewRGBA(dstRect)
		draw.Draw(dst, dstRect, face.Image, image.Point{r.X, r.Y}, draw.Src)
		filename := directory + string(filepath.Separator) + "face_" + id + "_" + strconv.Itoa(r.X) + "_" + strconv.Itoa(r.Y) + "_" + strconv.Itoa(r.Width) + "_" + strconv.Itoa(r.Height) + strconv.Itoa(i) + ".png"
		fdst, _ := os.Create(filename)
		defer fdst.Close()
		png.Encode(fdst, dst)
		fmt.Printf("File %s saved as png.\n", filename)
		filesPath = append(filesPath, filename)
	}
	return filesPath
}

func (face *FaceDetector) DrawOnImage() {

	for i, r := range face.GetFaces() {
		b := make([]byte, 16)
		rand.Read(b)
		id := fmt.Sprintf("%X", b)
		dstRect := image.Rect(r.X, r.Y, (r.X + r.Width), (r.Y + r.Height))
		//fmt.Println(dstRect)
		dst := image.NewRGBA(dstRect)
		draw.Draw(dst, dstRect, face.Image, image.Point{r.X, r.Y}, draw.Src)
		filename := "face_" + id + "_" + strconv.Itoa(r.X) + "_" + strconv.Itoa(r.Y) + "_" + strconv.Itoa(r.Width) + "_" + strconv.Itoa(r.Height) + strconv.Itoa(i) + ".png"
		fdst, _ := os.Create(filename)
		defer fdst.Close()
		png.Encode(fdst, dst)
		fmt.Printf("File %s saved as png.\n", filename)

	}

}
func NewFaceDectectorFromImage(imgData image.Image) *FaceDetector {
	var err error
	face := &FaceDetector{}
	defer func() {
		if err != nil {
			panic(err)
			return
		}
	}()
	face.Image = imgData
	if err != nil {
		return face
	}

	configXml.Do(func() {

		fxml, err := os.Open("haarcascade_frontalface_default.xml")
		if err != nil {
			panic(err.Error())
			return
		}
		defer fxml.Close()
		o := &OpenCVStorage{}
		err = xml.NewDecoder(fxml).Decode(o)
		if err != nil {
			return
		}
		o.Haarcascade.Stages.Parse()

		for _, stage := range o.Haarcascade.Stages.Stage {
			Config.Stages = append(Config.Stages, stage)
		}
	})

	face.Width = face.Image.Bounds().Max.X - face.Image.Bounds().Min.X
	face.Height = face.Image.Bounds().Max.Y - face.Image.Bounds().Min.Y

	face.ClassifiedSize = make([]int, 2)
	face.ClassifiedSize[0] = 24
	face.ClassifiedSize[1] = 24
	var maxScale float64
	if (float64(face.Width) / float64(face.ClassifiedSize[0])) > (float64(face.Height) / float64(face.ClassifiedSize[1])) {
		maxScale = float64(face.Height) / float64(face.ClassifiedSize[1])
	} else {
		maxScale = float64(face.Width) / float64(face.ClassifiedSize[0])
	}

	grayImage := make([][]float64, face.Width)
	for i := range grayImage {
		grayImage[i] = make([]float64, face.Height)
	}
	img := make([][]float64, face.Width)
	for i := range img {
		img[i] = make([]float64, face.Height)
	}
	square := make([][]float64, face.Width)
	for i := range square {
		square[i] = make([]float64, face.Height)
	}

	fmt.Println("Copies done.")
	for i := 0; i < face.Width; i++ {
		var col = 0.
		var col2 = 0.

		for j := 0; j < face.Height; j++ {
			r, g, b, _ := face.Image.At(i, j).RGBA()
			value := (((float64(r) * 255 / 65535) * 30) + (59 * (float64(g) * 255 / 65535)) + (11 * (float64(b) * 255 / 65535))) / 100
			//fmt.Printf("value:%f\n", value)
			img[i][j] = value
			grayImage[i][j] = col + value
			if i > 0 {
				grayImage[i][j] = grayImage[i-1][j] + col + value
			}
			square[i][j] = col2 + value*value
			if i > 0 {
				square[i][j] = square[i-1][j] + col2 + value*value
			}
			col += value
			col2 += value * value
			//fmt.Printf("col %f col2 %f img:%f square:%f gray:%f\n", col, col2, img[i][j], square[i][j], grayImage[i][j])
		}
	}

	fmt.Println("Stages passing.")
	baseScale := 2.
	scaleInc := 1.25
	increment := 0.1
	//minNeighbors := 3
	for scale := baseScale; scale < maxScale; scale *= scaleInc {
		step := int(scale * 24 * increment)
		size := int(scale * 24)
		for i := 0; i < (face.Width - size); i += step {
			for j := 0; j < (face.Height - size); j += step {
				pass := true
				k := 0
				for _, stage := range Config.Stages {

					if !stage.Pass(grayImage, square, i, j, scale) {
						pass = false
						break
					}
					k++
				}

				if pass == true {
					fr := &FoundRect{X: i, Y: j, Width: size, Height: size}
					fmt.Println(fr)
					face.FoundRects = append(face.FoundRects, fr)
				}
			}
		}

	}
	return face

}

func NewFaceDetector(imagePath string) *FaceDetector {
	var err error

	defer func() {
		if err != nil {
			panic(err)
			return
		}
	}()
	f, err := os.Open(imagePath)
	if err != nil {
		return &FaceDetector{}
	}
	defer f.Close()
	imgData, _, err := image.Decode(f)
	if err != nil {
		return &FaceDetector{}
	}

	return NewFaceDectectorFromImage(imgData)
}
