package api

import (
	"demotivator-generator/internal/domain"
	"net/http"
)

func GenerateDemotivator(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Add("Cache-Control", "s-max-age=60, stale-while-revalidate")

	demotivator := domain.Demotivator{}
	img := r.URL.Query().Get("img")
	demotivator.ImageURL = img
	demotivator.TextLine1 = r.URL.Query().Get("text1")
	demotivator.TextLine2 = r.URL.Query().Get("text2")

	err := demotivator.Generate(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
