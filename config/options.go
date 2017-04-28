package config

import (
	"fmt"
	"github.com/rkcpi/vell/repos"
	"github.com/rkcpi/vell/rpm"
	"os"
)

var (
	httpPort      = os.Getenv("VELL_HTTP_PORT")
	httpAddress   = os.Getenv("VELL_HTTP_ADDRESS")
	ReposPath     = os.Getenv("VELL_REPOS_PATH")
	RepoStore     repos.RepositoryStore
	ListenAddress string
)

func init() {
	if httpPort == "" {
		httpPort = "8080"
	}

	if ReposPath == "" {
		ReposPath = "/var/lib/vell/repositories"
	}
	RepoStore = rpm.NewRepositoryStore(ReposPath)

	ListenAddress = fmt.Sprintf("%s:%s", httpAddress, httpPort)
}
