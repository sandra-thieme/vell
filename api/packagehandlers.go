package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rkcpi/vell/rpm"
	"net/http"
)

// GET /repositories/{name}/packages
func ListPackages(w http.ResponseWriter, r *http.Request) {
	repo := rpm.YumRepository{mux.Vars(r)["name"]}
	var packages []rpm.Package
	if repo.IsValid() {
		packages = repo.ListPackages()
	}
	if err := json.NewEncoder(w).Encode(packages); err != nil {
		fail(w, err)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// POST /repositories/{name}/packages
func AddRPM(w http.ResponseWriter, r *http.Request) {
	repo := rpm.YumRepository{mux.Vars(r)["name"]}
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		fail(w, err)
	}
	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			src, err := file.Open()
			if err != nil {
				fail(w, err)
			}
			defer src.Close()

			repo.Add(file.Filename, src)
		}
	}
	repo.Update()
	w.WriteHeader(http.StatusCreated)
}
