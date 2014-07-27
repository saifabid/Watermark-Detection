package main

import (
	"encoding/json"
	"fmt"
	"github.com/oliamb/cutter"
	"github.com/otiai10/gosseract"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	tempFileName      = "./temp/temp.jpeg"
	croppedFileName   = "./temp/crop.jpeg"
	sample            = "girl-on-bed.jpeg"
	jsonFileTop       = "badUsersTOP.txt"
	jsonFileBottom    = "badUsersBOTTOM.txt"
	emberAPI          = "http://api.emberchatapp.com/api/v1/persona?offset=0&max=501"
	badPicturesTop    []faulty
	badPicturesBottom []faulty
	lowPoint          float64
)

type faulty struct {
	Name         string
	Picture      string
	markLocation string
}
type description struct {
	Id           string   `json:"_id"`
	Rev          string   `json:"_rev" `
	Name         string   `json: "name"`
	Pictures     []string `json:"pictures"`
	Age          string   `json:"age"`
	Location     string   `json:"location"`
	Weight       int      `json:"weight"`
	Height       int      `json:"height"`
	Measurements string   `json:"measurements"`
	Type         string   `json:"type"`
	Active       bool     `json:"active"`
	Created      string   `json:"created"`
	Utc          int      `json:"utc"`
	ChatterId    string   `json:"chatterId"`
	Preset       string   `json:"string"`
	Lastactive   string   `lastactive:"lastactive"`
	Snap         bool     `json:"snap"`
}

func filterImg(fileName string) (clean bool) {

	out := gosseract.Must(map[string]string{"src": fileName})
	out = strings.TrimSpace(out)
	if len(out) > 3 {

		fmt.Println(out)

		clean = false

	} else {

		clean = true
	}

	return

}

func getUsers(api string) (users []description) {

	resp, err := http.Get(emberAPI)
	checkErr(err)
	body, _ := ioutil.ReadAll(resp.Body)

	notgood := json.Unmarshal(body, &users)
	checkErr(notgood)
	return
}

// Scans bottom half of picture
func getPictureBottom(picture string) {
	resp, err := http.Get(picture)
	checkErr(err)

	body, _ := ioutil.ReadAll(resp.Body)

	err = ioutil.WriteFile(tempFileName, body, 0644)
	if err != nil {
		fmt.Println("Error")
	}
	saif, _ := os.Open(tempFileName)
	imgQ, _, errQ := image.Decode(saif)
	checkErr(errQ)
	lowPoint = float64(imgQ.Bounds().Max.Y) / 2

	// cropping image
	f, err := os.Open(tempFileName)
	checkErr(err)
	img, _, err := image.Decode(f)
	checkErr(err)

	cImg, err := cutter.Crop(img, cutter.Config{
		Height:  1,                             // height in pixel or Y ratio(see Ratio Option below)
		Width:   1,                             // width in pixel or X ratio
		Mode:    cutter.TopLeft,                // Accepted Mode: TopLeft, Centered
		Anchor:  image.Point{0, int(lowPoint)}, // Position of the top left point
		Options: cutter.Ratio,                  // Accepted Option: Ratio
	})

	imgw, _ := os.Create(croppedFileName)
	jpeg.Encode(imgw, cImg, &jpeg.Options{jpeg.DefaultQuality})

	os.Remove(tempFileName)

	return
}

// Scans top half of picture
func getPictureTop(picture string) {
	resp, err := http.Get(picture)
	checkErr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(tempFileName, body, 0644)
	checkErr(err)
	f, err := os.Open(tempFileName)
	checkErr(err)
	img, _, err := image.Decode(f)
	checkErr(err)
	cImg, err := cutter.Crop(img, cutter.Config{
		Height:  1,                 // height in pixel or Y ratio(see Ratio Option below)
		Width:   2,                 // width in pixel or X ratio
		Mode:    cutter.TopLeft,    // Accepted Mode: TopLeft, Centered
		Anchor:  image.Point{0, 0}, // Position of the top left point
		Options: cutter.Ratio,      // Accepted Option: Ratio
	})
	checkErr(err)
	imgw, _ := os.Create(croppedFileName)
	jpeg.Encode(imgw, cImg, &jpeg.Options{jpeg.DefaultQuality})

	os.Remove(tempFileName)

	return
}
func faultsToFileTop() {

	m, _ := json.Marshal(&badPicturesTop)
	ioutil.WriteFile(jsonFileTop, m, 0644)

}
func faultsToFileBottom() {

	m, _ := json.Marshal(&badPicturesBottom)
	ioutil.WriteFile(jsonFileBottom, m, 0644)

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	users := getUsers(emberAPI)

	for _, user := range users {
		pic := user.Pictures[0]
		getPictureBottom(pic)
		isClean := filterImg(croppedFileName)
		if isClean == false {
			person := faulty{user.Name, pic, "Bottom"}
			badPicturesBottom = append(badPicturesBottom, person)

			fmt.Println(pic)
		}

		os.Remove(croppedFileName)

		getPictureTop(pic)
		isClean = filterImg(croppedFileName)

		if isClean == false {
			person := faulty{user.Name, pic, "Top"}
			badPicturesTop = append(badPicturesTop, person)

			fmt.Println(pic)
		}
		os.Remove(croppedFileName)

	}

	faultsToFileBottom()
	faultsToFileTop()

}
