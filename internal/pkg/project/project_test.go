package project

import (
	"regexp"
	"testing"

	"github.com/tj/assert"
)

func TestRepositoriesHome(t *testing.T) {
	h, err := RepositoriesHome()
	if err != nil {
		t.Error(err)
	}
	assert.Regexp(t, regexp.MustCompile("/home/[a-z]+/.yaam/repositories"), h)
}
