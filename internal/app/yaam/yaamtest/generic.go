package yaamtest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/030/yaam/internal/app/yaam/project"
)

func GenericArtifact() error {
	if _, err := os.Stat(DirIso); errors.Is(err, os.ErrNotExist) {
		resp, err := http.Get("https://releases.ubuntu.com/22.04.1/ubuntu-22.04.1-desktop-amd64.iso")
		if err != nil {
			return err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		f, err := os.Create(DirIso)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return err
		}
		if err := f.Sync(); err != nil {
			return err
		}
	}

	return nil
}

func GenericArtifactReq(method, uri string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, project.Url+uri, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("hello", "world")

	client := &http.Client{
		Timeout: time.Second * 120,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GenericArtifactWriteOnDisk(resp *http.Response) (int64, error) {
	f, err := os.Create(DirIsoDownloaded)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	w, err := io.Copy(f, resp.Body)
	if err != nil {
		return 0, err
	}
	if err := f.Sync(); err != nil {
		return 0, err
	}
	return w, err
}
