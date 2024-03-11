package api

import (
	"bytes"
	"demotivator-generator/internal/domain"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
)

func GenerateDemotivator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Add("Cache-Control", "s-max-age=60, stale-while-revalidate")

	imageBase64, err := getImageAsBase64(r)
	if err != nil {
		http.Error(w, "Could not read image file", http.StatusInternalServerError)
		return
	}

	err = domain.Demotivator{
		ImageBase64: imageBase64,
		TextLine1:   r.URL.Query().Get("text1"),
		TextLine2:   r.URL.Query().Get("text2"),
	}.Generate(w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getImageAsBase64(r *http.Request) (string, error) {
	file, err, closeFunc := parseMultipartFile(r)
	if err != nil {
		return "", err
	}
	defer closeFunc()

	imageBase64, err := imageToBase64(file)
	return imageBase64, err
}

func parseMultipartFile(r *http.Request) (multipart.File, error, func()) {
	err := r.ParseMultipartForm(10 << 20) // max size 10MB
	if err != nil {
		return nil, err, nil
	}

	file, _, err := r.FormFile("file")
	return file, err, func() {
		err = file.Close()
		if err != nil {
			return
		}
	}
}

func imageToBase64(file multipart.File) (string, error) {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, file)
	return base64.StdEncoding.EncodeToString(buf.Bytes()), err
}
