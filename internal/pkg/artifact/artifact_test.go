package artifact

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tj/assert"
)

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
