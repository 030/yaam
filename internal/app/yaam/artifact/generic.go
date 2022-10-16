package artifact

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/030/yaam/internal/app/yaam/file"
	log "github.com/sirupsen/logrus"
)

type Generic struct {
	Request     *http.Request
	RequestBody io.ReadCloser
	RequestURI  string

	ResponseWriter http.ResponseWriter
}

func (g Generic) Publish() error {
	if err := StoreOnDisk(g.RequestURI, g.RequestBody); err != nil {
		return err
	}

	return nil
}

func (g Generic) Read() error {
	f, err := FilepathOnDisk(g.RequestURI)
	if err != nil {
		return err
	}

	filename := filepath.Base(f)
	log.Debugf("determined filename: '%s' in file: '%s'", filename, f)

	if _, exists := file.Exists(f); !exists {
		return fmt.Errorf("file: '%s' not found on disk", f)
	}

	g.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	g.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(g.ResponseWriter, g.Request, f)

	return nil
}
