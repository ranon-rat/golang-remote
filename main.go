package main

import (
	"bytes"
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

func clearRequest(str string) string {
	str = strings.Replace(str, "{", "", -1)
	str = strings.Replace(str, "}", "", -1)
	str = strings.Replace(str, "\"", "", -1)
	return str

}
func madeMapMouse(str string) map[string]int {
	str = clearRequest(str)
	strV := strings.Split(str, ",")
	m := map[string]int{}
	//idk why doesnt made the json for that i do this
	for _, i := range strV {
		n := strings.Split(i, ":")

		l, _ := strconv.Atoi(n[1])

		m[n[0]] = l
	}
	return m
}
func bodyRequest(r *http.Request) string {
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	//first decode the request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	newStr := buf.String()
	//decode the request to a string
	return newStr

}
func readMouse(w http.ResponseWriter, r *http.Request) {
	newStr := bodyRequest(r)
	m := madeMapMouse(newStr)
	robotgo.MoveMouse(m["x"], m["y"])
	//is only for move the mouse, then i gona make something interesting
}
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
func main() {
	r := mux.NewRouter()

	fs := http.Dir("page/")
	r.Handle("/", http.StripPrefix("/", http.FileServer(fs)))
	r.HandleFunc("/image/{a}", sendI)
	r.HandleFunc("/mouse", readMouse)
	r.HandleFunc("/command", readCommand)
	http.ListenAndServe(":8090", r)
}
