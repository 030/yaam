package project

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func Config() error {
	h, err := Home()
	if err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(h)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	log.Infof("config file used: '%s'", viper.ConfigFileUsed())

	return nil
}

func Home() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	h = filepath.Join(h, hiddenFolderName)

	if os.Getenv("YAAM_HOME") != "" {
		h = os.Getenv("YAAM_HOME")
	}
	log.Tracef("yaam home: '%s'", h)

	return h, nil
}

func RepositoriesHome() (string, error) {
	h, err := Home()
	if err != nil {
		return "", err
	}
	h = filepath.Join(h, "repositories")

	return h, nil
}
