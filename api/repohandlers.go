package api

import (
	"encoding/json"
	"fmt"
	"github.com/rkcpi/vell/repos"
	"github.com/rkcpi/vell/rpm"
	"io"
	"io/ioutil"
	"net/http"
)

// POST /repositories
func CreateRepo(w http.ResponseWriter, r *http.Request) {
	var repo repos.Repository
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &repo); err != nil {
		fail(w, err)
	}
	if err := rpm.NewRepository(repo.Name).Initialize(); err != nil {
		fail(w, err)
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Location", locationHeader(r, repo.Name))
	w.WriteHeader(http.StatusCreated)

}

// GET /repositories
func ListRepos(w http.ResponseWriter, r *http.Request) {
	reps := rpm.NewRepositoryStore().ListRepositories()

	if err := json.NewEncoder(w).Encode(reps); err != nil {
		fail(w, err)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func fail(w http.ResponseWriter, err error) {
	w.Header().Set("ContentType", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		panic(err)
	}
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
