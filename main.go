package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"
)

func main() {
	log.Println("Start Uploader.")

	//TODO startup routine to create folders and stuff
	// html, html/img, html/index.html, html/img/index.html

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Root: " + cwd)
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/", fs)
	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	out, err := os.Create("html/img/" + header.Filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Fprintf(w, "File uploaded")
}

func CheckMimeType(header textproto.MIMEHeader) bool {
	for k, v := range header {
		log.Println(k + " - " + v)
	}
	return false
}

func GenerateFileName() string {

}
