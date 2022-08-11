package api

import (
	"fmt"
	"net/http"
	"os"
)

func basicAuth(r *http.Request) error {
	u, p, ok := r.BasicAuth()
	if ok {
		if !(u == os.Getenv("YAAM_USER") && p == os.Getenv("YAAM_PASS")) {
			return fmt.Errorf("auth failed")
		}
	} else {
		return fmt.Errorf("request is NOT using basic authentication")
	}
	return nil
}

func Validation(method string, r *http.Request, w http.ResponseWriter) error {
	if !(method == "POST" || method == "GET" || method == "HEAD") {
		return fmt.Errorf("only POSTs, GETs and HEADs are supported. Method: '%s'", method)
	}

	if err := basicAuth(r); err != nil {
		return fmt.Errorf("basic auth failed. Error: '%v'", err)
	}
	return nil
}
