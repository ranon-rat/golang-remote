package main

import (
	"bytes"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kbinani/screenshot"
)

func sendI(w http.ResponseWriter, r *http.Request) {
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, _ := screenshot.CaptureRect(bounds)
		buffer := new(bytes.Buffer)
		png.Encode(buffer, img)
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		w.Write(buffer.Bytes())
	}
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/image/{a}", sendI)
	r.Handle("/", http.FileServer(http.Dir("./page")))
	http.ListenAndServe(":8090", r)
}
