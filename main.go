package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"
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

	// check if filetype is allowed
	if !CheckMimeType(header.Header) {
		log.Println("File type not allowed.")
		return
	}

	filename, err := SaveFile(file, strings.Split(header.Filename, ".")[1])

	fmt.Fprintf(w, "<a href='/img/"+filename+"'>"+filename+"</a>")
}

func CheckMimeType(header textproto.MIMEHeader) bool {
	mime := header["Content-Type"][0]
	allowed := [...]string{"image/jpg", "image/jpeg", "image/png", "image/webm", "image/gif"}

	for _, k := range allowed {
		if mime == k {
			return true
		}
	}
	return false
}

func SaveFile(src io.Reader, extension string) (string, error) {

	h := crc32.NewIEEE()
	dest, err := os.Create("tmpfile")
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	defer dest.Close()
	t := io.TeeReader(src, h)
	_, err = io.Copy(dest, t)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	filename := strconv.Itoa(int(h.Sum32())) + "." + extension
	err = os.Link("tmpfile", "html/img/"+filename)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return filename, nil
}
