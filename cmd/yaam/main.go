package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/030/yaam/internal/api"
	"github.com/030/yaam/internal/artifact"
	"github.com/030/yaam/internal/pkg/project"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const serverLogMsg = "check the server logs"

var Version string

func httpNotFoundReadTheLogs(w http.ResponseWriter, err error) {
	log.Error(err)
	http.Error(w, serverLogMsg, http.StatusNotFound)
}

func httpInternalServerErrorReadTheLogs(w http.ResponseWriter, err error) {
	log.Error(err)
	http.Error(w, serverLogMsg, http.StatusInternalServerError)
}

func mavenArtifact(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if err := api.Validation(r.Method, r, w); err != nil {
		httpInternalServerErrorReadTheLogs(w, err)
		return
	}

	m := artifact.Maven{RequestBody: r.Body, RequestURI: r.RequestURI, ResponseWriter: w}
	if r.Method == "PUT" {
		var p artifact.Publisher = m
		if err := p.Publish(); err != nil {
			httpInternalServerErrorReadTheLogs(w, err)
			return
		}
		return
	}

	var ap artifact.Preserver = m
	if err := ap.Preserve(); err != nil {
		httpNotFoundReadTheLogs(w, fmt.Errorf("maven artifact caching failed. Error: '%v'", err))
		return
	}

	var ar artifact.Reader = m
	if err := ar.Read(); err != nil {
		httpNotFoundReadTheLogs(w, fmt.Errorf("cannot read artifact from disk. Error: '%v'. Perhaps it resides in another repository?", err))
		return
	}
}

func mavenGroup(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if err := api.Validation(r.Method, r, w); err != nil {
		httpInternalServerErrorReadTheLogs(w, err)
		return
	}

	vars := mux.Vars(r)
	artifactURI := vars["artifact"]
	groupName := vars["name"]
	log.Debugf("Group: %v, Artifact: %v", groupName, artifactURI)
	var p artifact.Unifier = artifact.Maven{ResponseWriter: w, RequestURI: artifactURI}
	if err := p.Unify(groupName); err != nil {
		log.Error(fmt.Errorf("grouping failed. Error: '%v'", err))
		http.Error(w, serverLogMsg, http.StatusInternalServerError)
		return
	}
}

func genericArtifact(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if err := api.Validation(r.Method, r, w); err != nil {
		httpInternalServerErrorReadTheLogs(w, err)
		return
	}

	g := artifact.Generic{Request: r, RequestBody: r.Body, RequestURI: r.RequestURI, ResponseWriter: w}
	if r.Method == "POST" {
		var p artifact.Publisher = g
		if err := p.Publish(); err != nil {
			httpInternalServerErrorReadTheLogs(w, err)
			return
		}
		return
	}

	var ar artifact.Reader = g
	if err := ar.Read(); err != nil {
		httpNotFoundReadTheLogs(w, err)
		return
	}
}

func npmArtifact(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if err := api.Validation(r.Method, r, w); err != nil {
		httpInternalServerErrorReadTheLogs(w, err)
		return
	}

	n := artifact.Maven{RequestBody: r.Body, RequestURI: r.RequestURI, ResponseWriter: w}
	if r.Method == "POST" {
		var p artifact.Publisher = n
		if err := p.Publish(); err != nil {
			httpInternalServerErrorReadTheLogs(w, err)
			return
		}
		return
	}

	var ap artifact.Preserver = n
	if err := ap.Preserve(); err != nil {
		httpNotFoundReadTheLogs(w, err)
		return
	}

	var ar artifact.Reader = n
	if err := ar.Read(); err != nil {
		httpNotFoundReadTheLogs(w, err)
		return
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := io.WriteString(w, "ok"); err != nil {
		httpNotFoundReadTheLogs(w, err)
		return
	}
}

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	if os.Getenv("YAAM_DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	}

	r := mux.NewRouter()
	r.HandleFunc("/generic/{repo}/{artifact:.*}", genericArtifact)
	r.HandleFunc("/maven/groups/{name}/{artifact:.*}", mavenGroup)
	r.HandleFunc("/maven/{repo}/{artifact:.*}", mavenArtifact)
	r.HandleFunc("/npm/{repo}/{artifact:.*}", npmArtifact)
	r.HandleFunc("/status", status)

	srv := &http.Server{
		Addr: "0.0.0.0:" + project.PortString,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 120,
		ReadTimeout:  time.Second * 180,
		IdleTimeout:  time.Second * 240,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	log.Infof("Starting YAAM version: '%s' on localhost on port: '%d'...", Version, project.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
