package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tj/assert"
)

const (
	allowedReposMaven = `allowedRepos:
  - releases`
	allowedReposNpm = `allowedRepos:
  - 3rdparty-npm`
	gradleHomeDemoProject = "../../test/gradle/demo"
	npmHomeDemoProject    = "../../test/npm/demo"

	npmrc = `registry=http://localhost:25213/npm/3rdparty-npm/
always-auth=true
_auth=aGVsbG86d29ybGQ=`

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
	cmdExitErrMsg     = "%v, err: '%v'"
	mavenReleasesRepo = "maven/releases"
)

var (
	mavenRepos = []string{"maven/3rdparty-maven", "maven/3rdparty-maven-gradle-plugins", "maven/3rdparty-maven-spring", mavenReleasesRepo}
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
	if err := os.WriteFile(filepath.Join(reposDir, "maven.yaml"), []byte(allowedReposMaven), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reposDir, "npm.yaml"), []byte(allowedReposNpm), 0600); err != nil {
		return err
	}

	os.Setenv("YAAM_DEBUG", "true")
	os.Setenv("YAAM_USER", "hello")

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
cache=/tmp/yaam/test/npm/cache` + time.Now().Format("20060102150405111") + ``
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

	os.Setenv("GRADLE_USER_HOME", "/tmp/yaam/test/gradle"+time.Now().Format("20060102150405111"))
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
			url 'http://localhost:25213/` + repo + `/'
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

	assert.Regexp(t, `Could not PUT 'http://localhost:25213/maven/releases-non-existent/com/example/demo/0.0.1-SNAPSHOT/maven-metadata.xml'. Received status code 500 from server: Internal Server Error`, err)
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

	assert.Regexp(t, `Could not GET 'http://localhost:25213/maven/groups/helloworld/.*'. Received status code 500 from server: Internal Server Error`, err)
	assert.Equal(t, 1, exitCode)
}
