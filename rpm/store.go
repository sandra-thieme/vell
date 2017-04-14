package rpm

import (
	"github.com/rkcpi/vell/config"
	"github.com/rkcpi/vell/repos"
	"io/ioutil"
	"log"
)

type yumRepoStore struct {
	base string
}

func NewRepositoryStore() repos.RepositoryStore {
	return &yumRepoStore{config.ReposPath}
}

func (store *yumRepoStore) ListRepositories() []repos.Repository {
	files, err := ioutil.ReadDir(config.ReposPath)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	reps := make([]repos.Repository, 0, len(files))
	for _, file := range files {
		repo := repos.Repository{file.Name()}
		reps = append(reps, repo)
	}

	return reps
}
