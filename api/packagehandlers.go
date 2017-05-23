package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rkcpi/vell/config"
	"net/http"
	"errors"
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
	if !repo.IsValid() {
		return &apiError{errors.New("Repository does not exist"), "Invalid repository", http.StatusBadRequest}
	}
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		return &apiError{err, "I/O error", http.StatusBadRequest}
	}
	if len(r.MultipartForm.File) > 1 {
		return &apiError{errors.New("yolo"), "Invalid number of files", http.StatusBadRequest}
	}
	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			src, err := file.Open()
			if err != nil {
				return &apiError{err, "I/O error", http.StatusInternalServerError}
			}
			defer src.Close()

			err = repo.Add(file.Filename, src)
			if err != nil {
				return &apiError{err, "Error adding package to repository", http.StatusInternalServerError}
			}
		}
	}
	err = repo.Update()
	if err != nil {
		return &apiError{err, "Error updating the repository", http.StatusInternalServerError}
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}
