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

func replaceUrlPublicNpmWithYaamHost(f string) error {
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

func firstMatch(f, regex string) (string, error) {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(f)
	matchLength := len(match)
	log.Debugf("regex: '%s', match: '%v' matchLength: '%d' for file: '%s'", regex, match, matchLength, f)
	if matchLength <= 1 {
		return "", fmt.Errorf("no match was found for: '%s' with regex: '%s'", f, regex)
	}
	m := match[1]
	log.Debugf("firstMatch: '%s'", m)

	return m, nil
}

func pathTmp(f string) (string, error) {
	path, err := firstMatch(f, `^(/.*/[0-9a-z-\./@_]+)/-.*$`)
	if err != nil {
		return "", err
	}
	return filepath.Join(path + ".tmp"), nil
}

func versionShasum(f string) (string, error) {
	version, err := firstMatch(f, `-([0-9]+\.[0-9]+\.[0-9]+(-.+)?)\.tgz$`)
	if err != nil {
		return "", err
	}

	pt, err := pathTmp(f)
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(filepath.Clean(pt))
	if err != nil {
		return "", err
	}

	version = strings.Replace(version, ".", `\.`, -1)
	value := gjson.GetBytes(b, `versions.`+version+`.dist.shasum`)

	return value.String(), nil
}

func compareChecksumOnDiskWithExpectedSha(expChecksum, pathTmp string) (bool, error) {
	checksumValid := true
	f, err := os.Open(filepath.Clean(pathTmp))
	if err != nil {
		return checksumValid, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	/* #nosec */
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return checksumValid, err
	}
	fmt.Printf("%x", h.Sum(nil))
	checksum := fmt.Sprintf("%x", h.Sum(nil))
	if checksum != expChecksum {
		log.Errorf("file: '%s' checksum on disk: '%s' does not match expected checksum: '%s'", pathTmp, checksum, expChecksum)
		checksumValid = false
		time.Sleep(file.RetryDuration)
	}

	return checksumValid, nil
}

func checksum(f string) (bool, error) {
	checksumValid := true
	_, fileExists := file.Exists(f)
	if !fileExists && filepath.Ext(f) == ".tgz" {
		pt, err := pathTmp(f)
		if err != nil {
			return checksumValid, err
		}

		vs, err := versionShasum(f)
		if err != nil {
			return checksumValid, err
		}

		checksumValid, err := compareChecksumOnDiskWithExpectedSha(vs, pt)
		if err != nil {
			return checksumValid, err
		}
	}

	return checksumValid, nil
}

func (n Npm) downloadAgainIfInvalid(a artefact, resp *http.Response) error {
	checksumValid, err := checksum(a.path)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK || filepath.Ext(a.path) == ".tmp" {
		if err := file.CreateIfDoesNotExistInvalidOrEmpty(a.url, a.path, resp.Body, false); err != nil {
			return err
		}
	}

	if file.EmptyFile(a.path) || !checksumValid {
		if err := n.Preserve(); err != nil {
			return err
		}
	}

	if filepath.Ext(a.path) == ".tmp" {
		b, err := os.ReadFile(filepath.Clean(a.path))
		if err != nil {
			return err
		}

		if !json.Valid(b) {
			log.Errorf("json file: '%s' is invalid", a.path)
			if err := file.CreateIfDoesNotExistInvalidOrEmpty(a.url, a.path, resp.Body, true); err != nil {
				return err
			}
		}
	}

	return nil
}

func (n Npm) Preserve(urlStrings ...string) error {
	urlString := n.RequestURI
	if len(urlStrings) > 0 {
		urlString = urlStrings[0]
	}

	repoInConfigFile, err := artifact.RepoInConfigFile(n.ResponseWriter, urlString, "npm")
	if err != nil {
		return err
	}

	if !reflect.ValueOf(repoInConfigFile).IsZero() {
		h, err := project.RepositoriesHome()
		if err != nil {
			return err
		}
		dir := strings.Replace(urlString, "%2f", "/", -1)
		log.Debugf("extension found: '%s', file: '%s'", filepath.Ext(dir), dir)
		if filepath.Ext(dir) != ".tgz" {
			log.Debugf("file: '%s' does not have an extension", dir)
			dir = dir + ".tmp"
		}
		if err := artifact.Dir(dir); err != nil {
			return err
		}

		log.Debugf("downloadUrl before entering downloadUrl method: '%s', regex: '%s'", urlString, repoInConfigFile.Regex)
		du, err := artifact.DownloadUrl(repoInConfigFile.Url, repoInConfigFile.Regex, urlString)
		if err != nil {
			return err
		}
		completeFile := filepath.Join(h, dir)

		a := artefact{path: completeFile, url: du}
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

		if err := replaceUrlPublicNpmWithYaamHost(completeFile); err != nil {
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
