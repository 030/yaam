package file

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/030/yaam/internal/app/yaam/yaamtest"
	"github.com/tj/assert"
)

func TestExists(t *testing.T) {
	_, exists := Exists("does-not-exist")
	assert.Equal(t, false, exists)

}

func TestEmptyFile(t *testing.T) {
	empty := EmptyFile(filepath.Join("../..", yaamtest.TestdataCwd, "empty-file.txt"))
	assert.Equal(t, true, empty)
}

func TestEmptyFileFail(t *testing.T) {
	empty := EmptyFile("does-not-exist")
	assert.Equal(t, false, empty)
}

func TestCreateIfDoesNotExistInvalidOrEmpty(t *testing.T) {
	s := strings.NewReader("Hola mundo!")
	rc := io.NopCloser(s)

	err := CreateIfDoesNotExistInvalidOrEmpty("", "/tmp/yaam/testi.txt", rc, true)
	assert.NoError(t, err)
}

func TestCreateIfDoesNotExistInvalidOrEmptyFail(t *testing.T) {
	err := CreateIfDoesNotExistInvalidOrEmpty("", "/tmp2/file-does-not-exist", nil, true)
	assert.EqualError(t, err, "open /tmp2/file-does-not-exist: no such file or directory")
}
