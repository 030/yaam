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
	RetryDuration    = 5 * time.Second
	CannotReadErrMsg = "cannot read artifact from disk. Error: '%v'. Perhaps it resides in another repository?"
	WaitMsg          = "wait: '%v' before retrying"
)

func DownloadWithRetries(url string, auth ...string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if len(auth) > 0 {
		req.SetBasicAuth(auth[0], auth[1])
	}

	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 10 * time.Second
	retryClient.RetryWaitMax = 60 * time.Second
	standardClient := retryClient.StandardClient()

	/* #nosec */
	resp, err := standardClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Exists(f string) (int64, bool) {
	fi, err := os.Stat(f)
	if err != nil {
		return 0, false
	}

	return fi.Size(), true
}

func CreateIfDoesNotExistInvalidOrEmpty(url, f string, body io.ReadCloser, invalid bool) error {
	var written int64
	fileSize, fileExists := Exists(f)
	if !fileExists || fileSize == 0 || invalid {
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
	log.Infof("downloaded: '%s' to: '%s'. Wrote: '%d' bytes", url, f, written)

	return nil
}

func EmptyFile(f string) (emptyFile bool) {
	emptyFile = false
	fileSize, fileExists := Exists(f)
	if !fileExists {
		return emptyFile
	}

	if fileSize == 0 {
		log.Errorf("file: '%s' size is 0", f)
		log.Warnf(WaitMsg, RetryDuration)
		time.Sleep(RetryDuration)
		emptyFile = true
	}

	return emptyFile
}
