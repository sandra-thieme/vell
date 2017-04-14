package rpm

import (
	"github.com/rkcpi/vell/repos"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type yumRepository struct {
	store *yumRepoStore
	name  string
}

func NewRepository(store *yumRepoStore, name string) repos.AnyRepository {
	return &yumRepository{store, name}
}

func (r *yumRepository) path() string {
	return filepath.Join(r.store.base, r.name)
}

func (r *yumRepository) Add(filename string, f io.Reader) {
	log.Printf("Adding %s to repository %s", filename, r.path())
	destinationPath := filepath.Join(r.path(), filename)
	destination, err := os.Create(destinationPath)
	if err != nil {
		panic(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, f)
	if err != nil {
		panic(err)
	}
}

func (r *yumRepository) Update() error {
	path := r.path()
	log.Printf("Executing `createrepo --update %s`", path)
	cmd := exec.Command("createrepo", "--update", path)
	return cmd.Run()
}

func (r *yumRepository) repomdPath() string {
	return filepath.Join(r.path(), "repodata", "repomd.xml")
}

func (r *yumRepository) IsValid() bool {
	_, err := os.Stat(r.repomdPath())
	return err != nil
}

func (r *yumRepository) ListPackages() []repos.Package {
	files, _ := ioutil.ReadDir(r.path())
	packages := make([]repos.Package, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			p := repos.Package{
				Name:      file.Name(),
				Timestamp: file.ModTime().Format(time.RFC3339),
				Size:      file.Size(),
			}
			packages = append(packages, p)
		}
	}
	return packages
}
