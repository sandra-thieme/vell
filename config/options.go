package config

import (
	"fmt"
	"os"
)

var httpPort string
var httpAddress string
var ReposPath string

var ListenAddress string

func init() {
	if httpPort = os.Getenv("VELL_HTTP_PORT"); httpPort == "" {
		httpPort = "8080"
	}

	httpAddress = os.Getenv("VELL_HTTP_ADDRESS")

	if ReposPath = os.Getenv("VELL_REPOS_PATH"); ReposPath == "" {
		ReposPath = "/var/lib/vell/repositories"
	}
	if ReposPath[len(ReposPath)-1:] != "/" {
		ReposPath += "/"
	}

	ListenAddress = fmt.Sprintf("%s:%s", httpAddress, httpPort)
}
