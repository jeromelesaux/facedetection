package test

import (
	"encoding/xml"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"
	"testing"
)

type FoundRect struct {
	X      int64
	Y      int64
	Width  int64
	Height int64
}

type FaceDetector struct {
	ClassifiedSize []int64
	Stages         []*Stages
	Width          int64
	Height         int64
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

type RootNode struct {
	Feature   *Feature `xml:"feature"`
	Threshold float64  `xml:"threshold"`
	LeftVal   float64  `xml:"left_val"`
	RightVal  float64  `xml:"right_val"`
}

type Feature struct {
	Tilted float64        `xml:"tiltded"`
	Rects  []*RectsValues `xml:"rects"`
}
type RectsValues struct {
	V      string `xml:"_"`
	X1     int64
	Y1     int64
	X2     int64
	Y2     int64
	Weight float64
}

func (r *RectsValues) Parse() {
	var err error
	s := strings.Split(r.V, " ")
	fmt.Println(s)
	r.X1, err = strconv.ParseInt(s[0], 10, 8)
	if err != nil {
		fmt.Println(err.Error())
	}
	r.Y1, err = strconv.ParseInt(s[1], 10, 8)
	if err != nil {
		fmt.Println(err.Error())
	}
	r.X2, err = strconv.ParseInt(s[2], 10, 8)
	if err != nil {
		fmt.Println(err.Error())
	}
	r.Y2, err = strconv.ParseInt(s[3], 10, 8)
	if err != nil {
		fmt.Println(err.Error())
	}
	r.Weight, err = strconv.ParseFloat(s[4], 4)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(r)
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

func TestUnmarshallingXmlFile(t *testing.T) {
	f, err := os.Open("haarcascade_frontalface_min.xml")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	o := &OpenCVStorage{}
	err = xml.NewDecoder(f).Decode(o)
	if err != nil {
		t.Error(err)
	}
	o.Haarcascade.Stages.Parse()
	fmt.Println(*o.Haarcascade.Stages.Stage[0].Trees.Trees[0].RootNode.Feature.Rects[0])
	fmt.Println(o.Haarcascade.Stages.Stage[0].Threshold)
}
