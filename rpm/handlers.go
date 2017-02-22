package rpm

import (
	"encoding/json"
	"fmt"
	"github.com/rkcpi/vell/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateRepo(w http.ResponseWriter, r *http.Request) {
	var repo YumRepository
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &repo); err != nil {
		fail(w, err)
	}
	if err := repo.initialize(); err != nil {
		fail(w, err)
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Location", locationHeader(r, repo.Name))
	w.WriteHeader(http.StatusCreated)

}

func ListRepos(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(config.ReposPath)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	repos := make([]YumRepository, len(files), len(files))
	for i, file := range files {
		repos[i] = YumRepository{file.Name()}
	}
	if err := json.NewEncoder(w).Encode(repos); err != nil {
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
