package repos

import "io"

type Repository struct {
	Name string `json:"name"`
}

type Package struct {
	Name      string `json:"name"`
	Timestamp string `json:"lastUpdated"`
	Size      int64  `json:"size"`
	Version   string `json:"version"`
	Arch      string `json:"arch"`
}

type RepositoryStore interface {
	ListRepositories() []Repository

	Initialize(name string) error
	Get(name string) AnyRepository
}

type AnyRepository interface {
	Add(filename string, f io.Reader) error
	Update() error
	ListPackages() ([]Package, error)
	PackageWithNameAndVersion(packagename string, version string) (Package, error)
	IsValid() bool
}
