package file

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

const (
	RetryDuration    = 30 * time.Second
	CannotReadErrMsg = "cannot read artifact from disk. Error: '%v'. Perhaps it resides in another repository?"
	WaitMsg          = "wait: '%v' before retrying"
)

func DownloadWithRetries(url string) (*http.Response, error) {
	retryClient := retryablehttp.NewClient()

	retryClient.Logger = nil
	retryClient.RetryMax = 30
	retryClient.RetryWaitMin = 30 * time.Second
	retryClient.RetryWaitMax = 60 * time.Second
	standardClient := retryClient.StandardClient()
	log.Debugf("downloadURL: '%s'", url)

	/* #nosec */
	r, err := standardClient.Get(url)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func Exists(f string) (int64, bool) {
	fi, err := os.Stat(f)
	if err != nil {
		return 0, false
	}

	return fi.Size(), true
}

func CreateIfDoesNotExistOrEmpty(url, f string, body io.ReadCloser) error {
	var written int64
	fileSize, fileExists := Exists(f)
	if !fileExists || fileSize == 0 {
		dst, err := os.Create(filepath.Clean(f))
		if err != nil {
			return err
		}
		defer func() {
			if err := dst.Close(); err != nil {
				panic(err)
			}
		}()

		written, err = io.Copy(dst, body)
		if err != nil {
			return err
		}
		if err := dst.Sync(); err != nil {
			return err
		}
	}
	log.Debugf("downloaded: '%s' to: '%s'. Wrote: '%d' bytes", url, f, written)

	return nil
}

func EmptyFile(f string) (emptyFile bool) {
	fileSize, fileExists := Exists(f)
	if !fileExists {
		return false
	}

	if fileSize == 0 {
		log.Errorf("file: '%s' size is 0", f)
		log.Warnf(WaitMsg, RetryDuration)
		time.Sleep(RetryDuration)
		return true
	}

	return emptyFile
}
