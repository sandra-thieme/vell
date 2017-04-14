package rpm

import (
	"github.com/rkcpi/vell/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var repo YumRepository
var path string

func setup() {
	name := "vell-repository"
	repo = YumRepository{name}
	config.ReposPath, _ = ioutil.TempDir("", "vell")
	path = filepath.Join(config.ReposPath, name)
}

func TestPath(t *testing.T) {

	if p := repo.path(); p != path {
		t.Errorf("Expected %s, but got %s", path, p)
	}
}

func TestEnsureExists(t *testing.T) {
	repo.ensureExists()
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("%s", err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		t.Errorf("%s", err)
	}
	if !fileInfo.IsDir() {
		t.Errorf("%s is not a directory", path)
	}
	if perm := fileInfo.Mode().Perm(); perm != 0755 {
		t.Errorf("%s has wrong permissions: %s (expected %s)", path, perm, 0755)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func shutdown() {
	defer os.RemoveAll(config.ReposPath)
}
