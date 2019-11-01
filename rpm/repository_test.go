package rpm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var repo *yumRepository
var path string

func setup() {
	name := "vell-repository"
	reposPath, _ := ioutil.TempDir("", "vell")
	store := &yumRepoStore{reposPath}
	repo = &yumRepository{store, name}

	path = filepath.Join(reposPath, name)
}

func TestPath(t *testing.T) {
	if p := repo.path(); p != path {
		t.Errorf("Expected %s, but got %s", path, p)
	}
}

func TestEnsureExists(t *testing.T) {
	repo.store.ensureExists(repo.name)
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
		t.Errorf("%s has wrong permissions: %s (expected %s)", path, perm, os.FileMode(0755))
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func shutdown() {
	defer os.RemoveAll(repo.store.base)
}
