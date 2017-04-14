package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rkcpi/vell/config"
	"github.com/rkcpi/vell/repos"
	"net/http"
)

// GET /repositories/{name}/packages
func ListPackages(w http.ResponseWriter, r *http.Request) *apiError {
	repo := config.RepoStore.Get(mux.Vars(r)["name"])
	var packages []repos.Package
	if repo.IsValid() {
		packages = repo.ListPackages()
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
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
