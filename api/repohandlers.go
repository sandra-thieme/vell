package api

import (
	"encoding/json"
	"fmt"
	"github.com/rkcpi/vell/config"
	"github.com/rkcpi/vell/repos"
	"io"
	"io/ioutil"
	"net/http"
)

// POST /repositories
func CreateRepo(w http.ResponseWriter, r *http.Request) *apiError {
	var repo repos.Repository
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return &apiError{err, "I/O error", http.StatusBadRequest}
	}
	if err := json.Unmarshal(body, &repo); err != nil {
		return &apiError{err, "Invalid input", http.StatusBadRequest}
	}
	if err := config.RepoStore.Initialize(repo.Name); err != nil {
		return &apiError{err, "Repo initialization failed", http.StatusInternalServerError}
	}

	w.Header().Set("Location", locationHeader(r, repo.Name))
	w.WriteHeader(http.StatusCreated)
	return nil
}

// GET /repositories
func ListRepos(w http.ResponseWriter, r *http.Request) *apiError {
	reps := config.RepoStore.ListRepositories()

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(reps); err != nil {
		panic(err)
	}

	return nil
}

func locationHeader(r *http.Request, name string) string {
	var protocol string
	if r.TLS == nil {
		protocol = "http"
	} else {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/repositories/%s", protocol, r.Host, name)
}
