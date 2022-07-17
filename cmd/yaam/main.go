package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/030/yaam/internal/file"
)

func handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if !(method == "GET" || method == "HEAD") {
		log.Errorf("only GETs and HEADs are supported. Method: '%s'", method)
		http.Error(w, "only GETs and HEADs supported", http.StatusNotFound)
		return
	}

	maven2centralRegex := `^\/maven2central\/(.*)$`

	maven2central, err := regexp.MatchString(maven2centralRegex, r.URL.Path)
	if err != nil {
		log.Error(err)
		http.Error(w, "maven2central regex", http.StatusInternalServerError)
		return
	}

	downloadURL := ""
	if maven2central {
		re := regexp.MustCompile(maven2centralRegex)
		downloadURL = re.ReplaceAllString(r.URL.String(), `https://repo1.maven.org/maven2/$1`)
	} else {
		log.Errorf("only maven2central is supported. URLPath: '%s' did not match regex: '%s'", r.URL.Path, maven2centralRegex)
		http.Error(w, "not maven2central", http.StatusInternalServerError)
		return
	}

	home := filepath.Join("/tmp", "yaam")
	dir := filepath.Dir(r.URL.String())
	completeDir := filepath.Join(home, dir)
	if err := os.MkdirAll(completeDir, os.ModePerm); err != nil {
		log.Errorf("cannot create dir. Error: '%v'", err)
		http.Error(w, "cannot create dir", http.StatusInternalServerError)
		return
	}

	completeFile := filepath.Join(home, r.URL.String())
	if err := file.Download(completeFile, downloadURL); err != nil {
		log.Errorf("cannot download file. Error: '%v'", err)
		http.Error(w, "cannot download file", http.StatusInternalServerError)
		return
	}
	log.Infof("downloaded: '%s' to: '%s'", downloadURL, completeFile)

	b, err := os.ReadFile(completeFile)
	if err != nil {
		log.Errorf("cannot read file. Error: '%v'", err)
		http.Error(w, "cannot read file", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(b))
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":25113", nil); err != nil {
		log.Fatal(err)
	}
}
