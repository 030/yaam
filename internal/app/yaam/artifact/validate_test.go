package artifact

import (
	"path/filepath"
	"testing"

	"github.com/030/yaam/internal/app/yaam/project"
	"github.com/030/yaam/internal/app/yaam/yaamtest"
	log "github.com/sirupsen/logrus"
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

func TestValidate(t *testing.T) {
	///npm/3rdparty-npm/-/npm/v1/security/audits/quick
	expDir := filepath.Join("/maven", "releases", "world")
	err := validate(filepath.Join(expDir, "hola.mundo"))
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
}

func TestValidateFail(t *testing.T) {
	err := validate(filepath.Join("/something", "releases", "world"))
	assert.Regexp(t, err, "requestURI: '/something/releases/world' should start with a: '/' and contain an extension")
}
