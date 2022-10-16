package yaamtest

import (
	"testing"

	"github.com/tj/assert"
)

func TestGradleFile(t *testing.T) {
	err := gradleFile([]byte("hello"), "world")
	assert.EqualError(t, err, "open ../../test/testdata/gradle/demo/world.gradle: no such file or directory")
}
