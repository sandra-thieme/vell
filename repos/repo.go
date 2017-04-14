package repos

import "io"

type Repository struct {
	Name string `json:"name"`
}

type Package struct {
	Name      string `json:"name"`
	Timestamp string `json:"lastUpdated"`
	Size      int64  `json:"size"`
}

type RepositoryStore interface {
	ListRepositories() []Repository
}

type AnyRepository interface {
	Initialize() error
	Add(filename string, f io.Reader)
	Update() error
	ListPackages() []Package
	IsValid() bool
}
