package main

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 300)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		return
	}
	// check if filetype is allowed
	err = CheckMIME(file)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//hackish workaround, CheckMIME already reads some byte,
	//so the file is not complete when written.
	file.Close()
	file, header, _ = r.FormFile("file")
	defer file.Close()
	//TODO: dont parse fileextension without verifying
	filename, err := SaveFile(file, filepath.Ext(header.Filename))

	fmt.Fprintf(w, "<a href='/img/"+filename+"'>"+filename+"</a>")
}

func CheckMIME(file io.Reader) error {
	b := make([]byte, 512)
	if _, err := file.Read(b); err != nil {
		return err
	}
	mime := http.DetectContentType(b)
	if !strings.HasPrefix(mime, "video") && !strings.HasPrefix(mime, "image") {
		return (errors.New("filetype " + mime + " not allowed"))
	}
	return nil
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
	filename := strconv.Itoa(int(h.Sum32())) + extension
	err = os.Link("tmpfile", "html/img/"+filename)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return filename, nil
}
