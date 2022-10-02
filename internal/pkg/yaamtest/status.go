package yaamtest

import (
	"io"
	"net/http"
	"time"

	"github.com/030/yaam/internal/pkg/project"
)

func StatusHelper(method, uri string, body io.Reader, timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}
	req, err := http.NewRequest(method, project.Url+uri, body)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
