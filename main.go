package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/mux"
	"github.com/kbinani/screenshot"
)

type mouse struct {
	X     int             `json:"x"`
	Y     int             `json:"y"`
	Click clickAtrributes `json:"clickAttr"`
}
type clickAtrributes struct {
	Click bool   `json:"click"`
	Side  string `json:"side"`
}
type command struct {
	Command string `json:"command"`
}
type text struct {
	Text string `json:"word"`
}

func clearRequest(str string) string {
	str = strings.Replace(str, "{", "", -1)
	str = strings.Replace(str, "}", "", -1)
	str = strings.Replace(str, "\"", "", -1)
	return str

}

//---------------------> this was the tools for decode or compress
func compress(img image.Image) image.Image {
	var division float64 = 1.2
	upLeft := image.Point{0, 0}
	lowRight := image.Point{int(float64(img.Bounds().Max.X) / division), int(float64(img.Bounds().Max.Y) / division)}
	img2 := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for x := 0; x < int(float64(img.Bounds().Max.X)/division); x++ {
		for y := 0; y < int(float64(img.Bounds().Max.Y)/division); y++ {
			img2.Set(x, y, img.At(int(float64(x)*division), int(float64(y)*division)))
		}
	}
	return img2
}
func bodyRequest(r *http.Request) string {
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	newStr := buf.String()
	return newStr
}

//--------------------->  here is the remote control functions
func readMouse(w http.ResponseWriter, r *http.Request) {
	var mo mouse
	newStr := bodyRequest(r)
	json.Unmarshal([]byte(newStr), &mo)
	robotgo.MoveMouse(mo.X, mo.Y)
	if mo.Click.Click {
		robotgo.MouseClick(mo.Click.Side, mo.Click.Click)
	}

}
func readCommand(w http.ResponseWriter, r *http.Request) {
	mo := new(command)
	newStr := bodyRequest(r)
	json.Unmarshal([]byte(newStr), &mo)
	comm := strings.Split(mo.Command, " ")
	cmd := exec.Command(comm[0], comm[1:]...)
	cmd.Run()

}
func typeSomething(w http.ResponseWriter, r *http.Request) {
	var mo text
	newStr := bodyRequest(r)
	json.Unmarshal([]byte(newStr), &mo)
	robotgo.TypeStr(mo.Text)
}

//--------------------->  this is for get the image of the screen and send it

func readCommand(w http.ResponseWriter, r *http.Request) {

	req := bodyRequest(r)

	req = clearRequest(req)
	req1 := strings.Split(req, ":")
	command := strings.Split(req1[1], " ")
	cmd := exec.Command(command[0], command[1:]...)
	out, _ := cmd.CombinedOutput()
	fmt.Println(string(out))
	cmd.Run()

}
func compress(img image.Image) image.Image {
	var division float64 = 1.2
	upLeft := image.Point{0, 0}
	lowRight := image.Point{int(float64(img.Bounds().Max.X) / division), int(float64(img.Bounds().Max.Y) / division)}
	img2 := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for x := 0; x < int(float64(img.Bounds().Max.X)/division); x++ {
		for y := 0; y < int(float64(img.Bounds().Max.Y)/division); y++ {
			img2.Set(x, y, img.At(int(float64(x)*division), int(float64(y)*division)))
		}
	}
	return img2
}

func sendI(w http.ResponseWriter, r *http.Request) {

	n := screenshot.NumActiveDisplays()

	//this only take screenshots and send to the page
	for i := 0; i < n; i++ {

		bounds := screenshot.GetDisplayBounds(i)

		img, _ := screenshot.CaptureRect(bounds)
		buffer := new(bytes.Buffer)
		png.Encode(buffer, compress(img))
		//image encode the image and send the image
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		//image
		w.Write(buffer.Bytes())
	}
}

//--------------------->  the main
func main() {
	r := mux.NewRouter()
	fs := http.Dir("page/")
	r.Handle("/", http.FileServer(fs))
	r.HandleFunc("/image/{a}", sendI)
	r.HandleFunc("/mouse", readMouse)
	r.HandleFunc("/typetext", typeSomething)
	r.HandleFunc("/command", readCommand)
	http.ListenAndServe(":8090", r)
}
