package main

import (
	"golang.org/x/sys/unix"
	"log"
	"os/exec"
	"github.com/rkcpi/vell/config"
)

func EnsureCreateRepo() {
	_, err := exec.LookPath("createrepo")
	if err != nil {
		log.Fatal("Command `createrepo' is not available")
	}
}

func EnsureWritableReposPath() {
	err := unix.Access(config.ReposPath, unix.W_OK)
	if err != nil {
		log.Fatalf("Repositories location `%s' is not writable: %s", config.ReposPath, err)
	}
}