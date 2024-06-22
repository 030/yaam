package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/030/yaam/internal/app/yaam/project"
	"github.com/030/yaam/internal/app/yaam/yaamtest"
	"github.com/stretchr/testify/assert"
)

const (
	aptUri         = "/apt/3rdparty-ubuntu-nl-archive/some.iso"
	genericUri     = "/generic/something/some.iso"
	genericUriFail = "/generic/something2/some.iso"
)

func init() {
	if err := yaamtest.Config(); err != nil {
		log.Fatal(err)
	}

	go main()
}

func TestStatus(t *testing.T) {
	b, err := yaamtest.Status("GET", "/status", nil, 10)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, "ok", b)
}

// NPM: Preserve NPM artifacts by running `npm i` in a demo project.
func TestMainNpmBuild(t *testing.T) {
	exitCode, err := yaamtest.NpmConfig()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

// Apt.
func TestApt(t *testing.T) {
	resp, _ := yaamtest.GenericArtifactReq("GET", aptUri, nil)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, "check the server logs\n", bodyString)
}

/*
Maven: preserve artifacts by using the unify option that ensures that multiple
Maven repositories are grouped and only and endpoint has to be accessed and
subsequently publish a Maven artifact.
- GradleCleanBuildGroup
- GradleCleanBuildGroupFail
- GradlePublish
- GradlePublishFail.
*/
func TestGradleCleanBuildGroup(t *testing.T) {
	repos := []string{"maven/groups/hello"}
	if err := yaamtest.GradleBuildFile(repos, yaamtest.MavenReleasesRepo); err != nil {
		t.Error(err)
	}
	if err := yaamtest.GradleSettingsFile(repos); err != nil {
		t.Error(err)
	}

	exitCode, err := yaamtest.GradleCleanBuild("world")

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func TestGradleCleanBuildGroupFail(t *testing.T) {
	repos := []string{"maven/groups/helloworld"}
	if err := yaamtest.GradleBuildFile(repos, yaamtest.MavenReleasesRepo); err != nil {
		t.Error(err)
	}
	if err := yaamtest.GradleSettingsFile(yaamtest.MavenRepos); err != nil {
		t.Error(err)
	}

	exitCode, err := yaamtest.GradleCleanBuild("world")

	assert.Regexp(t, `Could not GET '`+project.Url+`/maven/groups/helloworld/.*'. Received status code 500 from server: Internal Server Error`, err)
	assert.Equal(t, 1, exitCode)
}

func TestGradlePublish(t *testing.T) {
	exitCode, err := yaamtest.GradlePublish(yaamtest.MavenReleasesRepo)

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func TestGradlePublishFail(t *testing.T) {
	exitCode, err := yaamtest.GradlePublish("maven/releases-non-existent")

	assert.Regexp(t, `Could not PUT '`+project.Url+`/maven/releases-non-existent/com/example/demo/0.0.1-SNAPSHOT/maven-metadata.xml'. Received status code 500 from server: Internal Server Error`, err)
	assert.Equal(t, 1, exitCode)
}

/*
Generic: upload an Ubuntu.iso and download it.
- Upload
- UploadFail
- Download
- DownloadFail.
*/
func TestGenericArtifactUpload(t *testing.T) {
	if err := yaamtest.GenericArtifact(); err != nil {
		t.Error(err)
	}

	b, err := os.ReadFile(yaamtest.DirIso)
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(b)
	_, err = yaamtest.GenericArtifactReq("POST", genericUri, r)
	if err != nil {
		t.Error(err)
	}
	assert.NoError(t, err)
}

func TestGenericArtifactUploadFail(t *testing.T) {
	resp, err := yaamtest.GenericArtifactReq("POST", genericUriFail, nil)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	bodyString := string(b)
	assert.Equal(t, "check the server logs\n", bodyString)
	assert.Equal(t, 500, resp.StatusCode)
	assert.NoError(t, err)
}

func TestGenericArtifactDownload(t *testing.T) {
	resp, err := yaamtest.GenericArtifactReq("GET", genericUri, nil)
	if err != nil {
		t.Error(err)
	}
	w, err := yaamtest.GenericArtifactWriteOnDisk(resp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "written: 3826831360", fmt.Sprintf("written: %d", w))
}

func TestGenericArtifactDownloadFail(t *testing.T) {
	resp, err := yaamtest.GenericArtifactReq("GET", genericUriFail, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 404, resp.StatusCode)
}
