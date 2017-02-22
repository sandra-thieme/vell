package rpm

import (
	"fmt"
	"github.com/rkcpi/vell/config"
	"log"
	"os"
	"os/exec"
)

type YumRepository struct {
	Name string `json:"name"`
}

func (r *YumRepository) ensureExists() string {
	path := r.path()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Creating repository directory %s", path)
		os.MkdirAll(path, 0755)
	}
	return path
}

func (r *YumRepository) path() string {
	return fmt.Sprintf("%s%s", config.ReposPath, r.Name)
}

func (r *YumRepository) initialize() error {
	log.Printf("Initializing repository %s", r.Name)
	path := r.ensureExists()
	log.Printf("Executing `createrepo --database %s`", path)
	cmd := exec.Command("createrepo", "--database", path)
	return cmd.Run()
}
