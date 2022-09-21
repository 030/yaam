package project

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

const (
	hiddenFolderName = ".yaam"
	Port             = 25213
	Host             = "localhost"
	Scheme           = "http"
)

var (
	PortString  = strconv.Itoa(Port)
	HostAndPort = Host + ":" + PortString
	Url         = Scheme + "://" + HostAndPort
)

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
