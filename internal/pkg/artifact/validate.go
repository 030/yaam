package artifact

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/030/yaam/internal/pkg/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

func validate(requestURI string) error {
	log.Debugf("requestURI: '%s'", requestURI)
	dir := filepath.Dir(requestURI)
	ext := filepath.Ext(requestURI)
	if dir == "/" || ext == "" {
		return fmt.Errorf("requestURI: '%s' should start with a: '/' and contain an extension", requestURI)
	}

	regex := `^/([a-z]+)/([0-9a-z-]+)/`
	re := regexp.MustCompile(regex)
	repoTypeAndName := re.FindStringSubmatch(requestURI)
	if len(repoTypeAndName) <= 2 {
		return fmt.Errorf("no repo type or name detected. Verify whether the regex: '%s' matches the URL: '%s'", regex, requestURI)
	}
	repoType := repoTypeAndName[1]
	repoName := repoTypeAndName[2]

	if err := allowedRepo(repoName, repoType); err != nil {
		return err
	}

	return nil
}

func allowedRepo(name, repoType string) error {
	h, err := project.Home()
	if err != nil {
		return err
	}

	viper.SetConfigName(repoType)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(h, "conf", "repositories"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	repos := viper.GetStringSlice("allowedRepos")
	if !slices.Contains(repos, name) {
		return fmt.Errorf("repository: '%s' is not allowed. Allowed repos: '%v'", name, repos)
	}

	return nil
}
