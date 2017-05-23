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

func (r *yumRepository) ensureExists() (string, error) {
	path := r.path()
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("Creating repository directory %s", path)
		err = os.MkdirAll(path, 0755)
	}
	return path, err
}

func (r *yumRepository) path() string {
	return filepath.Join(r.store.base, r.name)
}

func (r *yumRepository) Initialize() error {
	log.Printf("Initializing repository %s", r.name)
	path, err := r.ensureExists()
	if err != nil {
		return err
	}

	log.Printf("Executing `createrepo --database %s`", path)
	cmd := exec.Command("createrepo", "--database", path)
	return cmd.Run()
}

func (r *yumRepository) Add(filename string, f io.Reader) error {
	log.Printf("Adding %s to repository %s", filename, r.path())
	destinationPath := filepath.Join(r.path(), filename)
	destination, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, f)

	return err
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
	return err == nil
}

func (r *yumRepository) ListPackages() ([]repos.Package, error) {
	files, err := ioutil.ReadDir(r.path())
	if err != nil {
		return nil, err
	}

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
	return packages, nil
}
