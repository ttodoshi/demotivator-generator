package api

import (
	"bytes"
	"demotivator-generator/templates"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func GenerateDemotivator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	err = Demotivator{
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

	image, closeFunc := resizeImage(file)
	if closeFunc != nil {
		defer closeFunc()
	}

	imageBase64, err := imageToBase64(image)
	return imageBase64, err
}

func resizeImage(file multipart.File) (reader io.Reader, closeFunc func()) {
	reader = file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partWriter, err := writer.CreateFormFile("file", "file.png")
	if err != nil {
		return
	}
	_, err = io.Copy(partWriter, file)
	if err != nil {
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	err = writer.Close()
	if err != nil {
		return
	}

	request, err := http.NewRequest(http.MethodGet, os.Getenv("RESIZE_URL"), body)
	q := request.URL.Query()
	q.Add("width", "400")
	q.Add("height", "400")
	q.Add("save-proportions", "false")
	request.URL.RawQuery = q.Encode()

	if err != nil {
		return
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	return response.Body, func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}
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

func imageToBase64(file io.Reader) (string, error) {
	buf := bytes.Buffer{}
	_, err := io.Copy(&buf, file)
	return base64.StdEncoding.EncodeToString(buf.Bytes()), err
}

type Demotivator struct {
	ImageBase64 string
	TextLine1   string
	TextLine2   string
}

func (d Demotivator) Generate(resultWriter io.Writer) error {
	template, err := templates.GetTemplate()
	if err != nil {
		return err
	}

	err = template.Execute(resultWriter, d)
	if err != nil {
		return err
	}
	return nil
}
