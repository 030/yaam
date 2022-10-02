package artifact

import (
	"log"
	"testing"

	"github.com/030/yaam/internal/pkg/project"
	"github.com/030/yaam/internal/pkg/yaamtest"
	"github.com/tj/assert"
)

func init() {
	if err := yaamtest.Config(); err != nil {
		log.Fatal(err)
	}

	if err := project.Config(); err != nil {
		log.Fatal(err)
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
	expRepos := []string(nil)
	actRrepos, err := allowedRepos("world")

	assert.Equal(t, expRepos, actRrepos)
	assert.EqualError(t, err, "group: 'world' not found in config file")
}
