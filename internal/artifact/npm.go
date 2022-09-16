package artifact

import (
	"crypto/sha1" // #nosec
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/030/yaam/internal/pkg/artifact"
	"github.com/030/yaam/internal/pkg/file"
	"github.com/030/yaam/internal/pkg/project"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Npm struct {
	ResponseWriter http.ResponseWriter
	RequestBody    io.ReadCloser
	RequestURI     string
}

func replaceUrlPublicNpmWithYaamHost(f, url string) error {
	if filepath.Ext(f) == ".tmp" {
		input, err := os.ReadFile(filepath.Clean(f))
		if err != nil {
			return err
		}

		host := os.Getenv("YAAM_HOST")
		if host == "" {
			host = project.HostAndPort
		}
		output := strings.Replace(string(input), "https://registry.npmjs.org", "http://"+host+"/npm/3rdparty-npm", -1)

		var re = regexp.MustCompile(`(/@[a-z]+)(/)`)
		s := re.ReplaceAllString(output, `$1%2f`)
		err = os.WriteFile(f, []byte(s), 0600)
		if err != nil {
			return err
		}

		b, err := os.ReadFile(filepath.Clean(f))
		if err != nil {
			return err
		}
		if !json.Valid(b) {
			return fmt.Errorf("json for file: '%s' is invalid", f)
		}
	}
	return nil
}

func npm(url string, repoInConfigFile artifact.PublicRepository) (artefact, error) {
	h, err := project.RepositoriesHome()
	if err != nil {
		return artefact{}, err
	}
	dir := strings.Replace(url, "%2f", "/", -1)
	log.Debugf("extension found: '%s', file: '%s'", filepath.Ext(dir), dir)
	if filepath.Ext(dir) != ".tgz" {
		log.Debugf("file: '%s' does not have an extension", dir)
		dir = dir + ".tmp"
	}
	if err := artifact.Dir(dir); err != nil {
		return artefact{}, err
	}

	log.Debugf("downloadUrl before entering downloadUrl method: '%s', regex: '%s'", url, repoInConfigFile.Regex)
	du, err := artifact.DownloadUrl(repoInConfigFile.Url, repoInConfigFile.Regex, url)
	if err != nil {
		return artefact{}, err
	}
	completeFile := filepath.Join(h, dir)
	log.Debugf("completeFile: '%s', downloadUrl: '%s'", completeFile, du)

	if err := replaceUrlPublicNpmWithYaamHost(completeFile, url); err != nil {
		return artefact{}, err
	}

	return artefact{path: completeFile, url: du}, err
}

func checksum(f string) (bool, error) {
	checksumValid := true
	_, fileExists := file.Exists(f)
	if !fileExists && filepath.Ext(f) == ".tgz" {
		re := regexp.MustCompile(`-([0-9]+\.[0-9]+\.[0-9]).tgz$`)
		match := re.FindStringSubmatch(f)
		log.Debugf("match version length: '%d' for file: '%s'", len(match), f)
		version := match[1]
		fmt.Println(version)

		re = regexp.MustCompile(`^(/.*/[0-9a-z-/@]+)/-.*$`)
		match = re.FindStringSubmatch(f)
		log.Debugf("match tmp dir length: '%d' for file: '%s'", len(match), f)
		blah := match[1]
		fmt.Println(blah)

		blahFile := filepath.Join(blah + ".tmp")
		b, err := os.ReadFile(filepath.Clean(blahFile))
		if err != nil {
			return checksumValid, err
		}

		version = strings.Replace(version, ".", `\.`, -1)
		value := gjson.GetBytes(b, `versions.`+version+`.dist.shasum`)
		println(value.String())

		f2, err := os.Open(filepath.Clean(blahFile))
		if err != nil {
			return checksumValid, err
		}
		defer func() {
			if err := f2.Close(); err != nil {
				panic(err)
			}
		}()

		/* #nosec */
		h := sha1.New()
		if _, err := io.Copy(h, f2); err != nil {
			return checksumValid, err
		}

		fmt.Printf("%x", h.Sum(nil))
		checksum := fmt.Sprintf("%x", h.Sum(nil))
		if checksum != value.String() {
			log.Errorf("file: '%s' checksum on disk: '%s' does not match expected checksum: '%s'", f, checksum, value.String())
			checksumValid = false
			log.Warnf(file.WaitMsg, file.RetryDuration)
			time.Sleep(file.RetryDuration)
		}
	}

	return checksumValid, nil
}

func (n Npm) downloadAgainIfInvalid(a artefact, resp *http.Response) error {
	checksumValid, err := checksum(a.path)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK || !checksumValid || filepath.Ext(a.path) == ".tmp" {
		if err := file.CreateIfDoesNotExistOrEmpty(a.url, a.path, resp.Body); err != nil {
			return err
		}
	}

	if file.EmptyFile(a.path) {
		if err := n.Preserve(); err != nil {
			return err
		}
	}

	return nil
}

func (n Npm) Preserve() error {
	repoInConfigFile, err := artifact.RepoInConfigFile(n.ResponseWriter, n.RequestURI)
	if err != nil {
		return err
	}

	if !reflect.ValueOf(repoInConfigFile).IsZero() {
		a, err := npm(n.RequestURI, repoInConfigFile)
		if err != nil {
			return err
		}

		resp, err := file.DownloadWithRetries(a.url)
		if err != nil {
			return err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		if err := n.downloadAgainIfInvalid(a, resp); err != nil {
			return err
		}
	}

	return nil
}

func (n Npm) Publish() error {
	if err := artifact.StoreOnDisk(n.RequestURI, n.RequestBody); err != nil {
		return err
	}

	return nil
}

func (n Npm) Read() error {
	reqUrlString := strings.Replace(n.RequestURI, "%2f", "/", -1)
	if filepath.Ext(reqUrlString) != ".tgz" {
		log.Debugf("file: '%s' does not have an extension", reqUrlString)
		reqUrlString = reqUrlString + ".tmp"
	}
	if err := artifact.ReadFromDisk(n.ResponseWriter, reqUrlString); err != nil {
		return fmt.Errorf(file.CannotReadErrMsg, err)
	}

	return nil
}
