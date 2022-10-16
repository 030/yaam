package artifact

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/030/yaam/internal/app/yaam/yaamtest"
	"github.com/tj/assert"
)

func TestFirstMatch(t *testing.T) {
	exp := `world`
	act, err := firstMatch("helloworld", `(`+exp+`)`)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exp, act)
}

func TestFirstMatchFail(t *testing.T) {
	_, err := firstMatch("hello", "world")
	assert.EqualError(t, err, "no match was found for: 'hello' with regex: 'world'")
}

func TestPathTmp(t *testing.T) {
	exp := `/hello/world.tmp`
	act, err := pathTmp("/hello/world/-dfsf")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exp, act)
}

func TestPathTmpFail(t *testing.T) {
	_, err := pathTmp("helloworld")
	assert.EqualError(t, err, "no match was found for: 'helloworld' with regex: '^(/.*/[0-9a-z-\\./@_]+)/-.*$'")
}

func TestVersionShasum(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Error("err")
	}
	dirname := filepath.Dir(filename)

	exp := "7f4934d0f7ca8c56f95314939ddcd2dd91ce1d55"
	act, err := versionShasum(filepath.Join(dirname, "../../../..", yaamtest.Testdata, "npm/y18n/-/y18n-5.0.8.tgz"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exp, act)
}

func TestVersionShasumFail(t *testing.T) {
	_, err := versionShasum("does/not/exist")
	assert.EqualError(t, err, "no match was found for: 'does/not/exist' with regex: '-([0-9]+\\.[0-9]+\\.[0-9]+(-.+)?)\\.tgz$'")
}
