package rpm

import (
	"github.com/rkcpi/vell/config"
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
	Name string
}

func NewRepository(name string) repos.AnyRepository {
	return &yumRepository{name}
}

func (r *yumRepository) ensureExists() string {
	path := r.path()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Creating repository directory %s", path)
		os.MkdirAll(path, 0755)
	}
	return path
}

func (r *yumRepository) path() string {
	return filepath.Join(config.ReposPath, r.Name)
}

func (r *yumRepository) Initialize() error {
	log.Printf("Initializing repository %s", r.Name)
	path := r.ensureExists()
	log.Printf("Executing `createrepo --database %s`", path)
	cmd := exec.Command("createrepo", "--database", path)
	return cmd.Run()
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
			p := repos.Package{file.Name(), file.ModTime().Format(time.RFC3339), file.Size()}
			packages = append(packages, p)
		}
	}
	return packages
}
