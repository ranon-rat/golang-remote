package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

//-------------------> this decode the body request into a string
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

func sendI(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request")
	upgrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	imageWebSocket(ws)

}
func reader(conn *websocket.Conn) {
	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("fuck")
			return
		}
		n := screenshot.NumActiveDisplays()
		for {
			for i := 0; i < n; i++ {

				bounds := screenshot.GetDisplayBounds(i)

				img, _ := screenshot.CaptureRect(bounds)
				buffer := new(bytes.Buffer)
				png.Encode(buffer, compress(img))
				encoded := base64.StdEncoding.EncodeToString(buffer.Bytes())

				//image encode the image and send the image
				if err := conn.WriteMessage(messageType, []byte(encoded)); err != nil {
					log.Println(err)
				}

			}
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("se ha conectado , por fin carajo >:(")
	reader(ws)
}
func imageWebSocket(ws *websocket.Conn) {

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println("fuck")
			return
		}
		log.Println(string(p), messageType)

		n := screenshot.NumActiveDisplays()

		for i := 0; i < n; i++ {

			bounds := screenshot.GetDisplayBounds(i)

			img, _ := screenshot.CaptureRect(bounds)
			// decode the image into a base64
			buffer := new(bytes.Buffer)
			png.Encode(buffer, compress(img))
			encoded := base64.StdEncoding.EncodeToString(buffer.Bytes())

			//image encode the image and send the image
			if err := ws.WriteMessage(messageType, []byte(encoded)); err != nil {
				log.Println(err)
			}

		}
	}
}
func setupRoutes() {

	http.Handle("/", http.FileServer(http.Dir("view/")))

	//routes
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/mouse", readMouse)
	http.HandleFunc("/typetext", typeSomething)
	http.HandleFunc("/command", readCommand)

}

//--------------------->  the main
func main() {

	setupRoutes()

	http.ListenAndServe(":8090", nil)
}
