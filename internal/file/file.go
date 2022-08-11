package file

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func Download(f string, url string) (err error) {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		out, err := os.Create(filepath.Clean(f))
		if err != nil {
			return err
		}
		defer func() {
			if err := out.Close(); err != nil {
				panic(err)
			}
		}()

		/* #nosec */
		r, err := http.Get(url)
		if err != nil {
			return err
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				panic(err)
			}
		}()

		if r.StatusCode != http.StatusOK {
			return fmt.Errorf("download failed: '%s'", r.Status)
		}

		w, err := io.Copy(out, r.Body)
		if err != nil {
			return err
		}
		log.Info(w)
	}

	return nil
}
