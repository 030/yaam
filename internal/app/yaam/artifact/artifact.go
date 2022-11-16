package artifact

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/030/yaam/internal/app/yaam/file"
	"github.com/030/yaam/internal/app/yaam/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type artefact struct {
	path, url string
}

func allowedRepos(name string) ([]string, error) {
	groups := viper.GetStringMapStringSlice("groups.maven")
	var repos []string
	if values, ok := groups[name]; ok {
		repos = values
	} else {
		return nil, fmt.Errorf("group: '%s' not found in config file", name)
	}

	return repos, nil
}

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
		log.Tracef("file: '%s' exists already", path)
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

	log.Tracef("reading file: '%s' from disk...", f)
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
func RepoInConfigFile(w http.ResponseWriter, urlString, artifactType string) (PublicRepository, error) {
	reposAndElements := viper.GetStringMap("caches." + artifactType)
	if len(reposAndElements) == 0 {
		return PublicRepository{}, fmt.Errorf("caches: '%s' not found in config file", artifactType)
	}

	for repo, elements := range reposAndElements {
		url, ok := elements.(map[string]interface{})["url"]
		if ok {
			log.Tracef("url: '%s'", url)
		}
		user, ok := elements.(map[string]interface{})["user"]
		if ok {
			log.Tracef("user: '%s'", user)
		}
		pass, ok := elements.(map[string]interface{})["pass"]
		if ok {
			log.Tracef("pass: **********")
		}

		log.Debugf("trying to cache artifact from: '%s'...", urlString)

		rr := repoRegex(repo, artifactType)
		log.Tracef("repoRegex: '%s'", rr)

		pr := PublicRepository{Name: repo, Regex: rr, Url: url.(string)}
		if user != nil && pass != nil {
			pr.User = user.(string)
			pr.Pass = pass.(string)
		}

		riu, err := repoInUrl(rr, urlString)
		if err != nil {
			return PublicRepository{}, err
		}
		log.Debugf("repoInUrl: '%t'", riu)

		if riu {
			return pr, nil
		}
	}

	return PublicRepository{}, nil
}

type PublicRepository struct {
	Name, Regex, Url, User, Pass string
}

func repoRegex(repo, repoType string) string {
	return `^/` + repoType + `/` + repo + `/(.*)$`
}

func repoInUrl(repoRegex, url string) (bool, error) {
	log.Tracef("check whether url: '%s' contains repo according to regex: '%s'", url, repoRegex)
	match, err := regexp.MatchString(repoRegex, url)
	if err != nil {
		return false, err
	}
	log.Tracef("outcome regex check: '%t'", match)

	return match, nil
}

func DownloadUrl(publicRepoUrl, regex, url string) (string, error) {
	log.Debugf("check whether url: '%s' matches regex: '%s'. Params -> publicRepoUrl: '%s', regex: '%s' and url: '%s'", url, regex, publicRepoUrl, regex, url)
	r := regexp.MustCompile(regex)
	match := r.FindStringSubmatch(url)
	log.Tracef("number of matching elements: %d. Content: '%v'", len(match), match)
	if len(match) != 2 {
		return "", fmt.Errorf("should be 3! publicRepoUrl: '%s', regex: '%s', url: '%s'", publicRepoUrl, regex, url)
	}

	u := r.ReplaceAllString(url, publicRepoUrl+`$1`)

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
