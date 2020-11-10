package main

import (
	"bytes"
	"fmt"
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

type com struct {
	com string
}

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
	fmt.Println(command)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	out, _ := cmd.Output()
	fmt.Println(string(out))
	fmt.Println()

	j := json.NewDecoder(strings.NewReader(req))
	j.Decode(&c)//idk why doest decode the request 
	fmt.Println(c, req)


}
func sendI(w http.ResponseWriter, r *http.Request) {
	n := screenshot.NumActiveDisplays()
	//this only take screenshots and send to the page
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, _ := screenshot.CaptureRect(bounds)
		buffer := new(bytes.Buffer)
		png.Encode(buffer, img)
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
