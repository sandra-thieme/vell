package config

import (
	"fmt"
	"os"
)

var (
	httpPort      = os.Getenv("VELL_HTTP_PORT")
	httpAddress   = os.Getenv("VELL_HTTP_ADDRESS")
	ReposPath     = os.Getenv("VELL_REPOS_PATH")
	ListenAddress string
)

func init() {
	if httpPort == "" {
		httpPort = "8080"
	}

	if ReposPath == "" {
		ReposPath = "/var/lib/vell/repositories"
	}
	if ReposPath[len(ReposPath)-1:] != "/" {
		ReposPath += "/"
	}

	ListenAddress = fmt.Sprintf("%s:%s", httpAddress, httpPort)
}
