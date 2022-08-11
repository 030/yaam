package artifact

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/030/yaam/internal/pkg/project"
	log "github.com/sirupsen/logrus"
)

func ReadFromDisk(w http.ResponseWriter, r *http.Request) error {
	prh, err := project.RepositoriesHome()
	if err != nil {
		return err
	}

	completeFile := filepath.Join(prh, r.URL.String())
	b, err := os.ReadFile(filepath.Clean(completeFile))
	if err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(b)); err != nil {
		return err
	}

	return nil
}

func storeOnDisk(dir, artifactPathUri string, content io.ReadCloser) error {
	prh, err := project.RepositoriesHome()
	if err != nil {
		return err
	}

	artifactHome := filepath.Join(prh, dir)

	if err := os.MkdirAll(artifactHome, os.ModePerm); err != nil {
		return err
	}

	artifactPath := filepath.Join(prh, artifactPathUri)

	if _, err := os.Stat(artifactPath); errors.Is(err, os.ErrNotExist) {
		log.Info(artifactPath)
		out, err := os.Create(filepath.Clean(artifactPath))
		if err != nil {
			return err
		}
		defer func() {
			if err := out.Close(); err != nil {
				panic(err)
			}
		}()

		written, err := io.Copy(out, content)
		if err != nil {
			log.Error(err)
		}
		log.Info(written)
	} else {
		log.Infof("file: '%s' exists already", artifactPath)
	}

	return nil
}

func Publish(r *http.Request) error {
	ext := filepath.Ext(r.RequestURI)
	dir := filepath.Dir(r.RequestURI)
	if ext == "" || dir == "/" {
		return fmt.Errorf("should contain a dir and have an extension")
	}

	if err := storeOnDisk(dir, r.RequestURI, r.Body); err != nil {
		return err
	}

	return nil
}
