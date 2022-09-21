package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/030/yaam/internal/pkg/project"
	"github.com/tj/assert"
)

const (
	allowedReposGeneric = `allowedRepos:
  - something`
	allowedReposMaven = `allowedRepos:
  - releases`
	allowedReposNpm = `allowedRepos:
  - 3rdparty-npm`
	gradleHomeDemoProject = "../../test/gradle/demo"
	npmHomeDemoProject    = "../../test/npm/demo"

	caches = `mavenReposAndUrls:
  3rdparty-maven: https://repo.maven.apache.org/maven2/
  3rdparty-maven-gradle-plugins: https://plugins.gradle.org/m2/
  3rdparty-maven-spring: https://repo.spring.io/release/
  3rdparty-npm: https://registry.npmjs.org/`
	groups = `groups:
  hello:
    - maven/releases
    - maven/3rdparty-maven
    - maven/3rdparty-maven-gradle-plugins
    - maven/3rdparty-maven-spring`
	cmdExitErrMsg        = "%v, err: '%v'"
	mavenReleasesRepo    = "maven/releases"
	testDir              = "/tmp/yaam"
	testDirGradle        = testDir + "/gradle"
	testDirNpm           = testDir + "/npm"
	testDirIso           = testDir + "/ubuntu.iso"
	testDirIsoDownloaded = testDir + "/downloaded-ubuntu.iso"
)

var (
	mavenRepos = []string{"maven/3rdparty-maven", "maven/3rdparty-maven-gradle-plugins", "maven/3rdparty-maven-spring", mavenReleasesRepo}
	npmrc      = `registry=` + project.Url + `/npm/3rdparty-npm/
	always-auth=true
	_auth=aGVsbG86d29ybGQ=`
)

func testConfigHelper() error {
	os.Setenv("YAAM_HOME", filepath.Join("/tmp", "yaam", "test"+time.Now().Format("20060102150405111")))
	dir := filepath.Join(os.Getenv("YAAM_HOME"), "conf")
	reposDir := filepath.Join(dir, "repositories")
	if err := os.MkdirAll(reposDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "caches.yaml"), []byte(caches), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "groups.yaml"), []byte(groups), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reposDir, "generic.yaml"), []byte(allowedReposGeneric), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reposDir, "maven.yaml"), []byte(allowedReposMaven), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reposDir, "npm.yaml"), []byte(allowedReposNpm), 0600); err != nil {
		return err
	}

	os.Setenv("YAAM_DEBUG", "true")
	os.Setenv("YAAM_USER", "hello")
	os.Setenv("YAAM_PASS", "world")

	return nil
}

func testNpmConfigHelper() (int, error) {
	os.Setenv("YAAM_PASS", "world")

	if err := os.RemoveAll(filepath.Join(npmHomeDemoProject, "node_modules")); err != nil {
		return 1, err
	}
	packageLockJson := filepath.Join(npmHomeDemoProject, "package-lock.json")
	if _, err := os.Stat(packageLockJson); err == nil {
		if err := os.Remove(packageLockJson); err != nil {
			return 1, err
		}
	}

	npmrcWithCacheLocation := npmrc + `
cache=` + testDirNpm + `/cache` + time.Now().Format("20060102150405111") + ``
	if err := os.WriteFile(filepath.Join(npmHomeDemoProject, ".npmrc"), []byte(npmrcWithCacheLocation), 0600); err != nil {
		return 1, err
	}

	cmd := exec.Command("bash", "-c", "npm cache clean --force && npm i")
	cmd.Dir = npmHomeDemoProject
	co, err := cmd.CombinedOutput()
	if err != nil {
		return cmd.ProcessState.ExitCode(), fmt.Errorf(cmdExitErrMsg, string(co), err)
	}

	return 0, nil
}

func init() {
	if err := testConfigHelper(); err != nil {
		panic(err)
	}

	go main()
}

func testMainGradleFile(content []byte, name string) error {
	if err := os.WriteFile(filepath.Join(gradleHomeDemoProject, name+".gradle"), content, 0600); err != nil {
		return err
	}
	return nil
}

func testMainGradleCleanBuildHelper(pass string) (int, error) {
	os.Setenv("YAAM_PASS", pass)

	os.Setenv("GRADLE_USER_HOME", testDirGradle+time.Now().Format("20060102150405111"))
	cmd := exec.Command("bash", "-c", "./gradlew clean build --no-daemon")
	cmd.Dir = gradleHomeDemoProject
	co, err := cmd.CombinedOutput()
	if err != nil {
		return cmd.ProcessState.ExitCode(), fmt.Errorf(cmdExitErrMsg, string(co), err)
	}

	return 0, nil
}

func testMainGradlePublishHelper(repo string) (int, error) {
	if err := testGradleBuildFileHelper(mavenRepos, repo); err != nil {
		return 1, err
	}
	if err := testGradleSettingsFileHelper(mavenRepos); err != nil {
		return 1, err
	}

	exitCode, err := testMainGradleCleanBuildHelper("world")
	if err != nil {
		return exitCode, err
	}

	cmd := exec.Command("bash", "-c", "./gradlew publish --no-daemon")
	cmd.Dir = gradleHomeDemoProject
	co, err := cmd.CombinedOutput()
	if err != nil {
		return cmd.ProcessState.ExitCode(), fmt.Errorf(cmdExitErrMsg, string(co), err)
	}

	return 0, nil
}

func testGradleMavenRepositoriesFileHelper(repos []string) string {
	var sb strings.Builder
	for _, repo := range repos {
		content := `
		maven {
			allowInsecureProtocol true
			url '` + project.Url + `/` + repo + `/'
			authentication {
				basic(BasicAuthentication)
			}
			credentials {
				username "hello"
				password "world"
			}
		}`
		sb.WriteString(content)
	}
	return sb.String()
}

func testGradlePublishingFileHelper(repo string) string {
	repos := []string{repo}
	content := `
publishing {
  publications {
    mavenJava(MavenPublication) {
      versionMapping {
        usage('java-api') {
          fromResolutionOf('runtimeClasspath')
        }
        usage('java-runtime') {
          fromResolutionResult()
        }
      }
    }
  }

  repositories {` +
		testGradleMavenRepositoriesFileHelper(repos) + `
  }
}
`
	return content
}

func testGradleBuildFileHelper(repos []string, repoPublish string) error {
	content := `
plugins {
  id 'org.springframework.boot' version '2.7.3'
  id 'io.spring.dependency-management' version '1.0.13.RELEASE'
  id 'java'
  id 'maven-publish'
}

group = 'com.example'
version = '0.0.1-SNAPSHOT'
sourceCompatibility = '17'

repositories {` +
		testGradleMavenRepositoriesFileHelper(repos) + `
}

` + testGradlePublishingFileHelper(repoPublish) + `

dependencies {
  implementation 'org.springframework.boot:spring-boot-starter'
  testImplementation 'org.springframework.boot:spring-boot-starter-test'
}

tasks.named('test') {
  useJUnitPlatform()
}
`

	if err := testMainGradleFile([]byte(content), "build"); err != nil {
		return err
	}
	return nil
}

func testGradleSettingsFileHelper(repos []string) error {
	content := `
pluginManagement {
  repositories {` +
		testGradleMavenRepositoriesFileHelper(repos) + `
  }
}

rootProject.name = 'demo'
`

	if err := testMainGradleFile([]byte(content), "settings"); err != nil {
		return err
	}
	return nil
}

func testStatusHelper(method, pass, uri string, body io.Reader, timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}
	req, err := http.NewRequest(method, project.Url+uri, body)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("hello", pass)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func testGenericArtifactHelper() error {
	if _, err := os.Stat(testDirIso); errors.Is(err, os.ErrNotExist) {
		resp, err := http.Get("https://releases.ubuntu.com/22.04.1/ubuntu-22.04.1-desktop-amd64.iso")
		if err != nil {
			return err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		f, err := os.Create(testDirIso)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return err
		}
		if err := f.Sync(); err != nil {
			return err
		}
	}

	return nil
}

func testGenericArtifactReqHelper(method, uri string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, project.Url+uri, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("hello", "world")

	client := &http.Client{
		Timeout: time.Second * 120,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func testGenericArtifactWriteOnDiskHelper(resp *http.Response) (int64, error) {
	f, err := os.Create(testDirIsoDownloaded)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	w, err := io.Copy(f, resp.Body)
	if err != nil {
		return 0, err
	}
	if err := f.Sync(); err != nil {
		return 0, err
	}
	return w, err
}

func TestGenericArtifact(t *testing.T) {
	// Upload
	if err := testGenericArtifactHelper(); err != nil {
		t.Error(err)
	}
	uri := "/generic/something/some.iso"
	b, err := os.ReadFile(testDirIso)
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(b)
	_, err = testGenericArtifactReqHelper("POST", uri, r)
	if err != nil {
		t.Error(err)
	}
	assert.NoError(t, err)

	// Download
	resp, err := testGenericArtifactReqHelper("GET", uri, nil)
	if err != nil {
		t.Error(err)
	}
	w, err := testGenericArtifactWriteOnDiskHelper(resp)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "written: 3826831360", fmt.Sprintf("written: %d", w))
}

func TestGenericArtifactFail(t *testing.T) {
	// Upload
	uri := "/generic/something2/some.iso"
	resp, err := testGenericArtifactReqHelper("POST", uri, nil)
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

	// Download
	resp, err = testGenericArtifactReqHelper("GET", uri, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 404, resp.StatusCode)
}

func TestStatus(t *testing.T) {
	b, err := testStatusHelper("GET", "world", "/status", nil, 10)
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, "ok", b)
}

func TestMainNpmBuild(t *testing.T) {
	exitCode, err := testNpmConfigHelper()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func TestMainGradleCleanBuild(t *testing.T) {
	if err := testGradleBuildFileHelper(mavenRepos, mavenReleasesRepo); err != nil {
		t.Error(err)
	}
	if err := testGradleSettingsFileHelper(mavenRepos); err != nil {
		t.Error(err)
	}

	exitCode, err := testMainGradleCleanBuildHelper("world")
	if err != nil {
		t.Error(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func TestMainGradleCleanBuildFail(t *testing.T) {
	if err := testConfigHelper(); err != nil {
		t.Error(err)
	}

	if err := testGradleBuildFileHelper(mavenRepos, mavenReleasesRepo); err != nil {
		t.Error(err)
	}
	if err := testGradleSettingsFileHelper(mavenRepos); err != nil {
		t.Error(err)
	}

	exitCode, err := testMainGradleCleanBuildHelper("incorrectPass")

	assert.Regexp(t, "was not found in any of the following sources", err)
	assert.Equal(t, 1, exitCode)
}

func TestMainGradleCleanBuildNonMavenFail(t *testing.T) {
	repos := []string{"3rdparty-maven", "3rdparty-maven-gradle-plugins", "3rdparty-maven-spring", "releases"}
	if err := testGradleBuildFileHelper(repos, "releases"); err != nil {
		t.Error(err)
	}
	if err := testGradleSettingsFileHelper(repos); err != nil {
		t.Error(err)
	}

	exitCode, err := testMainGradleCleanBuildHelper("world")

	assert.Regexp(t, "was not found in any of the following sources", err)
	assert.Equal(t, 1, exitCode)
}

func TestMainGradlePublish(t *testing.T) {
	exitCode, err := testMainGradlePublishHelper(mavenReleasesRepo)

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func TestMainGradlePublishFail(t *testing.T) {
	exitCode, err := testMainGradlePublishHelper("maven/releases-non-existent")

	assert.Regexp(t, `Could not PUT '`+project.Url+`/maven/releases-non-existent/com/example/demo/0.0.1-SNAPSHOT/maven-metadata.xml'. Received status code 500 from server: Internal Server Error`, err)
	assert.Equal(t, 1, exitCode)
}

func TestMainGradleCleanBuildGroup(t *testing.T) {
	repos := []string{"maven/groups/hello"}
	if err := testGradleBuildFileHelper(repos, mavenReleasesRepo); err != nil {
		t.Error(err)
	}
	if err := testGradleSettingsFileHelper(repos); err != nil {
		t.Error(err)
	}

	exitCode, err := testMainGradleCleanBuildHelper("world")

	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
}

func TestMainGradleCleanBuildGroupFail(t *testing.T) {
	repos := []string{"maven/groups/helloworld"}
	if err := testGradleBuildFileHelper(repos, mavenReleasesRepo); err != nil {
		t.Error(err)
	}
	if err := testGradleSettingsFileHelper(mavenRepos); err != nil {
		t.Error(err)
	}

	exitCode, err := testMainGradleCleanBuildHelper("world")

	assert.Regexp(t, `Could not GET '`+project.Url+`/maven/groups/helloworld/.*'. Received status code 500 from server: Internal Server Error`, err)
	assert.Equal(t, 1, exitCode)
}
