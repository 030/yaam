package artifact

import (
	"io"
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/030/yaam/internal/app/yaam/project"
	"github.com/030/yaam/internal/app/yaam/yaamtest"
	"github.com/stretchr/testify/assert"
)

func init() {
	if err := yaamtest.Config(); err != nil {
		log.Fatal(err)
	}

	if err := project.Config(); err != nil {
		log.Fatal(err)
	}
}

func TestStoreOnDisk(t *testing.T) {
	s := strings.NewReader("Hola mundo!")
	rc := io.NopCloser(s)

	err := StoreOnDisk(filepath.Join("/maven/releases/world", "hola.mundo"), rc)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
}

func TestStoreOnDiskFail(t *testing.T) {
	err := StoreOnDisk(filepath.Join("/maven/releases-not-allowed/world", "hola.mundo"), nil)

	assert.EqualError(t, err, "repository: 'releases-not-allowed' is not allowed. Allowed repos: '[releases]'")
}

const testUrl = "/hello/world"

func TestRepoInUrlTrue(t *testing.T) {
	match, err := repoInUrl("hello", testUrl)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, match)
}

func TestRepoInUrlFalse(t *testing.T) {
	match, err := repoInUrl("hello123", testUrl)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, false, match)
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
	expRepos := []string(nil)
	actRrepos, err := allowedRepos("world")

	assert.Equal(t, expRepos, actRrepos)
	assert.EqualError(t, err, "group: 'world' not found in config file")
}
