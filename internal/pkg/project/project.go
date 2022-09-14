package project

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

const hiddenFolderName = ".yaam"

func Home() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	yh := filepath.Join(h, hiddenFolderName)

	if os.Getenv("YAAM_HOME") != "" {
		yh = os.Getenv("YAAM_HOME")
	}
	log.Debugf("yaam home: '%s'", yh)

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
