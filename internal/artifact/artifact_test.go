package artifact

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tj/assert"
)

const (
	groups = `groups:
  hello:
    - maven/releases
    - maven/3rdparty-maven
    - maven/3rdparty-maven-gradle-plugins
    - maven/3rdparty-maven-spring`
)

func init() {
	os.Setenv("YAAM_HOME", filepath.Join("/tmp", "yaam", "test"+time.Now().Format("20060102150405111")))
	dir := filepath.Join(os.Getenv("YAAM_HOME"), "conf")
	reposDir := filepath.Join(dir, "repositories")
	if err := os.MkdirAll(reposDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "groups.yaml"), []byte(groups), 0600); err != nil {
		panic(err)
	}
}

func TestAllowedRepos(t *testing.T) {
	expRepos := []string{"maven/releases", "maven/3rdparty-maven", "maven/3rdparty-maven-gradle-plugins", "maven/3rdparty-maven-spring"}
	actRrepos, err := allowedRepos("hello")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expRepos, actRrepos)
}

func TestAllowedReposFail(t *testing.T) {
	expRepos := []string{}
	actRrepos, err := allowedRepos("world")

	assert.Equal(t, expRepos, actRrepos)
	assert.EqualError(t, err, "group: 'world' not found in config file")
}
