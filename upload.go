package main

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var DataDir = "data/"

type UploadHandler struct {
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	uploadedFile, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer uploadedFile.Close()
	ext := filepath.Ext(handler.Filename)

	f, err := ioutil.TempFile(DataDir, ".tempfile-*")
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
	os.Rename(f.Name(), filepath.Join(DataDir, name+ext))
	log.Printf("size:%d name:%s", size, name)
}
