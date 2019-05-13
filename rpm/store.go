package rpm

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rkcpi/vell/repos"
)

type yumRepoStore struct {
	base string
}

func NewRepositoryStore(base string) repos.RepositoryStore {
	return &yumRepoStore{base}
}

func (store *yumRepoStore) Get(name string) repos.AnyRepository {
	return NewRepository(store, name)
}

func (store *yumRepoStore) Initialize(name string) error {
	log.Printf("Initializing repository %s", name)
	path := store.ensureExists(name)
	log.Printf("Executing `createrepo --database %s`", path)
	cmd := exec.Command("createrepo", "--database", path)
	return cmd.Run()
}

func (store *yumRepoStore) ListRepositories() []repos.Repository {
	files, err := ioutil.ReadDir(store.base)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	reps := make([]repos.Repository, 0, len(files))
	for _, file := range files {
		repo := repos.Repository{Name: file.Name()}
		reps = append(reps, repo)
	}

	return reps
}

func (store *yumRepoStore) ensureExists(name string) string {
	path := filepath.Join(store.base, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Creating repository directory %s", path)
		os.MkdirAll(path, 0755)
	}
	return path
}
