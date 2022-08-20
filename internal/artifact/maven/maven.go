package maven

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/030/yaam/internal/file"
	"github.com/030/yaam/internal/pkg/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func cache(repo, repoURL string, reqURL *url.URL) error {
	prh, err := project.RepositoriesHome()
	if err != nil {
		return err
	}

	repoRegex := `^\/` + repo + `\/(.*)$`
	reqURLString := reqURL.String()

	match, err := regexp.MatchString(repoRegex, reqURL.Path)
	if err != nil {
		return err
	}

	downloadURL := ""
	if match {
		re := regexp.MustCompile(repoRegex)
		downloadURL = re.ReplaceAllString(reqURLString, repoURL+`$1`)

		dir := filepath.Dir(reqURLString)
		completeDir := filepath.Join(prh, dir)
		if err := os.MkdirAll(completeDir, os.ModePerm); err != nil {
			return err
		}

		completeFile := filepath.Join(prh, reqURLString)
		log.Info(completeFile)
		if err := file.Download(completeFile, downloadURL); err != nil {
			return err
		}
		log.Infof("downloaded: '%s' to: '%s'", downloadURL, completeFile)
	}

	return nil
}

func Cache(w http.ResponseWriter, reqURL *url.URL) error {
	yh, err := project.Home()
	if err != nil {
		return err
	}

	viper.SetConfigName("repositories")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(yh, "conf"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	reposAndUrls := viper.GetStringMapString("mavenReposAndUrls")
	for repo, url := range reposAndUrls {
		if err := cache(repo, url, reqURL); err != nil {
			return err
		}
	}

	return nil
}
