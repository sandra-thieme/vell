package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rkcpi/vell/config"
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

// GET /repositories/{name}/packages/{packagename}/version/{version}
func GetPackageWithNameAndVersion(w http.ResponseWriter, r *http.Request) *apiError {
	repo := config.RepoStore.Get(mux.Vars(r)["name"])
	packagename := mux.Vars(r)["packagename"]
	version := mux.Vars(r)["version"]

	pkg, err := repo.PackageWithNameAndVersion(packagename, version)
	if err != nil {
		return &apiError{err, "Package not found", http.StatusNotFound}
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(pkg); err != nil {
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

	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
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

	if err := repo.Update(); err != nil {
		return &apiError{err, "Error updating the repository", http.StatusInternalServerError}
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}
