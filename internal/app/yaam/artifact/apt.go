package artifact

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/030/yaam/internal/app/yaam/file"
	log "github.com/sirupsen/logrus"
)

type Apt struct {
	ResponseWriter http.ResponseWriter
	RequestBody    io.ReadCloser
	RequestURI     string
}

func (a Apt) downloadAgainIfInvalid(atf artefact, resp *http.Response) error {
	log.Trace(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		if err := file.CreateIfDoesNotExistInvalidOrEmpty(atf.url, atf.path, resp.Body, false); err != nil {
			return err
		}
	}

	if file.EmptyFile(atf.path) {
		if err := a.Preserve(); err != nil {
			return err
		}
	}

	return nil
}

func (a Apt) Preserve(urlStrings ...string) error {
	urlString := a.RequestURI
	if len(urlStrings) > 0 {
		urlString = urlStrings[0]
	}
	log.Tracef("urlString: '%s'", urlString)

	repoInConfigFile, err := RepoInConfigFile(a.ResponseWriter, urlString, "apt")
	if err != nil {
		return err
	}

	if !reflect.ValueOf(repoInConfigFile).IsZero() {
		atf, err := maven(urlString, repoInConfigFile)
		if err != nil {
			return err
		}

		resp, err := file.DownloadWithRetries(atf.url, repoInConfigFile.User, repoInConfigFile.Pass)
		if err != nil {
			return err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		if err := a.downloadAgainIfInvalid(atf, resp); err != nil {
			return err
		}
	}

	return nil
}

func (a Apt) Read() error {
	if err := ReadFromDisk(a.ResponseWriter, a.RequestURI); err != nil {
		return fmt.Errorf(file.CannotReadErrMsg, err)
	}

	return nil
}
