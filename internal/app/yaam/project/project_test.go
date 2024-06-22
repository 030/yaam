package project

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	h, err := Home()
	if err != nil {
		t.Error(err)
	}
	assert.Regexp(t, regexp.MustCompile("/home/[a-z]+/.yaam"), h)
}

func TestRepositoriesHome(t *testing.T) {
	h, err := RepositoriesHome()
	if err != nil {
		t.Error(err)
	}
	assert.Regexp(t, regexp.MustCompile("/home/[a-z]+/.yaam/repositories"), h)
}
