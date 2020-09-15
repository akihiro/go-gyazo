package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type UploadHandler struct {
	DataDir     string
	BaseURL     string
	MaxFileSize int64
}

func (h UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(h.MaxFileSize); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uploadedFile, handler, err := r.FormFile("imagedata")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer uploadedFile.Close()
	ext := filepath.Ext(handler.Filename)

	f, err := ioutil.TempFile(h.DataDir, ".tempfile-*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer f.Close()

	hgen := sha256.New()
	size, err := io.Copy(io.MultiWriter(f, hgen), uploadedFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	name := base64.RawURLEncoding.EncodeToString(hgen.Sum(nil))
	os.Rename(f.Name(), filepath.Join(h.DataDir, name+ext))
	log.Printf("size:%d name:%s", size, name)
	fmt.Fprintf(w, "%s/data/%s%s", h.BaseURL, name, ext)
}
