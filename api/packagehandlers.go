package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rkcpi/vell/config"
	"net/http"
)

// GET /repositories/{name}/packages
func ListPackages(w http.ResponseWriter, r *http.Request) *apiError {
	repo := config.RepoStore.Get(mux.Vars(r)["name"])

	packages, err := repo.ListPackages()
	if err != nil {
		return &apiError{err, "Package listing error", http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(packages); err != nil {
		panic(err)
	}
	return nil
}

// POST /repositories/{name}/packages
func AddRPM(w http.ResponseWriter, r *http.Request) *apiError {
	repo := config.RepoStore.Get(mux.Vars(r)["name"])
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		return &apiError{err, "I/O error", http.StatusBadRequest}
	}
	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			src, err := file.Open()
			if err != nil {
				return &apiError{err, "I/O error", http.StatusInternalServerError}
			}
			defer src.Close()

			repo.Add(file.Filename, src)
		}
	}
	repo.Update()
	w.WriteHeader(http.StatusCreated)
	return nil
}
