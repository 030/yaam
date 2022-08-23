package project

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const hiddenFolderName = ".yaam"

func Home() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	yh := filepath.Join(h, hiddenFolderName)

	return yh, nil
}

func RepositoriesHome() (string, error) {
	h, err := Home()
	if err != nil {
		return "", err
	}

	rh := filepath.Join(h, "repositories")

	return rh, nil
}
