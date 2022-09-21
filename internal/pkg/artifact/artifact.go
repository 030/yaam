package artifact

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/030/yaam/internal/pkg/file"
	"github.com/030/yaam/internal/pkg/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func createHomeAndReturnPath(requestURI string) (string, error) {
	h, err := project.RepositoriesHome()
	if err != nil {
		return "", err
	}

	artifactPath := filepath.Join(h, requestURI)
	artifactHome := filepath.Dir(artifactPath)
	if err := os.MkdirAll(artifactHome, os.ModePerm); err != nil {
		return "", err
	}

	return artifactPath, nil
}

func createIfDoesNotExist(path string, requestBody io.ReadCloser) error {
	if _, fileExists := file.Exists(path); !fileExists {
		dst, err := os.Create(filepath.Clean(path))
		if err != nil {
			return err
		}

		defer func() {
			if err := dst.Close(); err != nil {
				panic(err)
			}
		}()

		w, err := io.Copy(dst, requestBody)
		if err != nil {
			log.Error(err)
		}
		log.Debugf("file: '%s' created and it contains: '%d' bytes", path, w)
		if err := dst.Sync(); err != nil {
			return err
		}
	} else {
		log.Debugf("file: '%s' exists already", path)
	}
	return nil
}

func StoreOnDisk(requestURI string, requestBody io.ReadCloser) error {
	if err := validate(requestURI); err != nil {
		return err
	}

	path, err := createHomeAndReturnPath(requestURI)
	if err != nil {
		return err
	}

	if err := createIfDoesNotExist(path, requestBody); err != nil {
		return err
	}

	return nil
}

func FilepathOnDisk(url string) (string, error) {
	h, err := project.RepositoriesHome()
	if err != nil {
		return "", err
	}

	f := filepath.Join(h, url)
	log.Debugf("constructed filepath: '%s' after concatenating home: '%s' to url: '%s'", f, h, url)
	return f, nil
}

func ReadFromDisk(w http.ResponseWriter, reqURL string) error {
	f, err := FilepathOnDisk(reqURL)
	if err != nil {
		return err
	}

	log.Debugf("reading file: '%s' from disk...", f)
	b, err := os.ReadFile(filepath.Clean(f))
	if err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(b)); err != nil {
		return err
	}

	return nil
}

// ReadRepositoriesAndUrlsFromConfigFileAndCacheArtifact reads a repositories
// yaml file that contains repositories and their URLs. If a request is
// attempted to download a file, it will look up the name in the config file
// and find the public URLs so it can download the file from the public maven
// repository and cache it on disk.
func RepoInConfigFile(w http.ResponseWriter, urlString string) (PublicRepository, error) {
	yh, err := project.Home()
	if err != nil {
		return PublicRepository{}, err
	}

	viper.SetConfigName("caches")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(yh, "conf"))
	if err := viper.ReadInConfig(); err != nil {
		return PublicRepository{}, err
	}

	reposAndUrls := viper.GetStringMapString("mavenReposAndUrls")
	for repo, url := range reposAndUrls {

		// if err := pr.cache(n.RequestURL); err != nil {
		// 	return err
		// }
		// reqURLString := reqURL.String()
		log.Debugf("trying to cache artifact from: '%s'...", urlString)

		rr := repoRegex(repo)
		log.Debugf("repoRegex: '%s'", rr)
		pr := PublicRepository{Name: repo, Regex: rr, Url: url}
		riu, err := repoInUrl(rr, urlString)
		if err != nil {
			return PublicRepository{}, err
		}
		log.Debugf("repoInUrl: '%t'", riu)

		if riu {
			// if err := pr.createDirAndStoreOnDisk(rr, reqURLString); err != nil {
			// 	return err
			// }
			return pr, nil
		}
	}

	return PublicRepository{}, nil
}

type PublicRepository struct {
	Name, Regex, Url string
}

func repoRegex(repo string) string {
	return `^/(maven|npm)/` + repo + `/(.*)$`
}

func repoInUrl(repoRegex, url string) (bool, error) {
	log.Debugf("check whether url: '%s' contains repo according to regex: '%s'", url, repoRegex)
	match, err := regexp.MatchString(repoRegex, url)
	if err != nil {
		return false, err
	}
	log.Debugf("outcome regex check: '%t'", match)

	return match, nil
}

func DownloadUrl(publicRepoUrl, regex, url string) (string, error) {
	log.Debugf("check whether url: '%s' matches regex: '%s'. Params -> publicRepoUrl: '%s', regex: '%s' and url: '%s'", url, regex, publicRepoUrl, regex, url)
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(url)
	log.Debugf("number of matching elements: %d", len(match))
	if len(match) != 3 {
		return "", fmt.Errorf("should be 3! publicRepoUrl: '%s', regex: '%s', url: '%s'", publicRepoUrl, regex, url)
	}

	u := re.ReplaceAllString(url, publicRepoUrl+`$2`)

	return u, nil
}

func Dir(path string) error {
	h, err := project.RepositoriesHome()
	if err != nil {
		return err
	}

	dir := filepath.Join(h, filepath.Dir(path))

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return nil
}
