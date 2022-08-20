package main

import (
	"net/http"
	"strconv"

	"github.com/030/yaam/internal/api"
	"github.com/030/yaam/internal/artifact"
	"github.com/030/yaam/internal/artifact/maven"
	log "github.com/sirupsen/logrus"
)

const port = 25213

func httpInternalServerErrorReadTheLogs(w http.ResponseWriter) {
	http.Error(w, "check the server logs", http.StatusInternalServerError)
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	method := r.Method
	reqURL := r.URL

	if err := api.Validation(method, r, w); err != nil {
		log.Errorf("request is invalid. Error: '%v'", err)
		httpInternalServerErrorReadTheLogs(w)
		return
	}

	if method == "POST" {
		if err := artifact.Publish(r); err != nil {
			log.Errorf("publish of an artifact failed. Error: '%v'", err)
			httpInternalServerErrorReadTheLogs(w)
			return
		}
		return
	}

	if err := maven.Cache(w, reqURL); err != nil {
		log.Errorf("maven artifact caching failed. Error: '%v'", err)
		httpInternalServerErrorReadTheLogs(w)
		return
	}

	if err := artifact.ReadFromDisk(w, r); err != nil {
		log.Warnf("cannot read artifact from disk. Error: '%v'. Perhaps it resides in another repository?", err)
		http.Error(w, "check the server logs", http.StatusNotFound)
		return
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	http.HandleFunc("/", handler)

	log.Infof("Starting YAAM on localhost on port: '%d'...", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}
