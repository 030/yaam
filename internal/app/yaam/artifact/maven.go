package artifact

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/030/yaam/internal/app/yaam/file"
	"github.com/030/yaam/internal/app/yaam/project"
	log "github.com/sirupsen/logrus"
)

type Maven struct {
	ResponseWriter http.ResponseWriter
	RequestBody    io.ReadCloser
	RequestURI     string
}

func maven(url string, repoInConfigFile PublicRepository) (artefact, error) {
	h, err := project.RepositoriesHome()
	if err != nil {
		return artefact{}, err
	}

	if err := Dir(url); err != nil {
		return artefact{}, err
	}

	du, err := DownloadUrl(repoInConfigFile.Url, repoInConfigFile.Regex, url)
	if err != nil {
		return artefact{}, err
	}

	completeFile := filepath.Join(h, url)
	log.Debugf("completeFile: '%s', downloadUrl: '%s'", completeFile, du)

	return artefact{path: completeFile, url: du}, nil
}

func (m Maven) downloadAgainIfInvalid(a artefact, resp *http.Response) error {
	log.Debug(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		if err := file.CreateIfDoesNotExistInvalidOrEmpty(a.url, a.path, resp.Body, false); err != nil {
			return err
		}
	}

	if file.EmptyFile(a.path) {
		if err := m.Preserve(); err != nil {
			return err
		}
	}

	return nil
}

func (m Maven) Preserve(urlStrings ...string) error {
	urlString := m.RequestURI
	if len(urlStrings) > 0 {
		urlString = urlStrings[0]
	}
	log.Debugf("urlString: '%s'", urlString)

	repoInConfigFile, err := RepoInConfigFile(m.ResponseWriter, urlString, "maven")
	if err != nil {
		return err
	}

	if !reflect.ValueOf(repoInConfigFile).IsZero() {
		a, err := maven(urlString, repoInConfigFile)
		if err != nil {
			return err
		}

		resp, err := file.DownloadWithRetries(a.url, repoInConfigFile.User, repoInConfigFile.Pass)
		if err != nil {
			return err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		if err := m.downloadAgainIfInvalid(a, resp); err != nil {
			return err
		}
	}

	return nil
}

func (m Maven) Publish() error {
	if err := StoreOnDisk(m.RequestURI, m.RequestBody); err != nil {
		return err
	}

	return nil
}

func (m Maven) Read() error {
	if err := ReadFromDisk(m.ResponseWriter, m.RequestURI); err != nil {
		return fmt.Errorf(file.CannotReadErrMsg, err)
	}

	return nil
}

func (m Maven) Unify(name string) error {
	repos, err := allowedRepos(name)
	if err != nil {
		return err
	}

	log.Debugf("repos: '%v'", repos)
	for _, repo := range repos {
		log.Debugf("repo: '%s'", repo)
		urlString := "/" + repo + "/" + m.RequestURI
		log.Debugf("urlString: '%s'", urlString)

		h, err := project.RepositoriesHome()
		if err != nil {
			return err
		}

		if err := m.Preserve(urlString); err != nil {
			log.Errorf("maven artifact caching failed. Error: '%v'", err)
		}

		if _, fileExists := file.Exists(filepath.Join(h, urlString)); fileExists {
			if err := ReadFromDisk(m.ResponseWriter, urlString); err != nil {
				log.Warnf(file.CannotReadErrMsg, err)
			}
			return nil
		}
	}

	return nil
}
