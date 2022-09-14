package artifact

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tj/assert"
)

const (
	allowedReposMaven = `allowedRepos:
  - releases`
	allowedReposNpm = `allowedRepos:
  - 3rdparty-npm`
)

func ConfigHelper() error {
	os.Setenv("YAAM_HOME", filepath.Join("/tmp", "yaam", "test"+time.Now().Format("20060102150405111")))
	dir := filepath.Join(os.Getenv("YAAM_HOME"), "conf", "repositories")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "maven.yaml"), []byte(allowedReposMaven), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "npm.yaml"), []byte(allowedReposNpm), 0600); err != nil {
		return err
	}

	os.Setenv("YAAM_DEBUG", "true")
	os.Setenv("YAAM_USER", "hello")

	return nil
}

func init() {
	if err := ConfigHelper(); err != nil {
		panic(err)
	}
}

func TestValidate(t *testing.T) {
	expDir := filepath.Join("/maven", "releases", "world")
	err := validate(filepath.Join(expDir, "hola.mundo"))
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
}

func TestValidateFail(t *testing.T) {
	expErr := `Config File "something" Not Found in "\[/tmp/yaam/test[0-9]+/conf/repositories\]"`
	actErr := validate(filepath.Join("/something", "releases", "world", "hola.mundo"))

	assert.Regexp(t, expErr, actErr)
}
