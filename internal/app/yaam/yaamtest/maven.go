package yaamtest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/030/yaam/internal/app/yaam/project"
)

var (
	MavenReleasesRepo = "maven/releases"
	MavenRepos        = []string{"maven/3rdparty-maven", "maven/3rdparty-maven-gradle-plugins", "maven/3rdparty-maven-spring", MavenReleasesRepo}
)

func gradleFile(content []byte, name string) error {
	if err := os.WriteFile(filepath.Join(gradleHomeDemoProject, name+".gradle"), content, 0o600); err != nil {
		return err
	}
	return nil
}

func GradleMavenRepositoriesFile(repos []string) string {
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

func GradlePublishingFile(repo string) string {
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
		GradleMavenRepositoriesFile(repos) + `
  }
}
`
	return content
}

func GradlePublish(repo string) (int, error) {
	if err := GradleBuildFile(MavenRepos, repo); err != nil {
		return 1, err
	}
	if err := GradleSettingsFile(MavenRepos); err != nil {
		return 1, err
	}

	exitCode, err := GradleCleanBuild("world")
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

func GradleSettingsFile(repos []string) error {
	content := `
pluginManagement {
  repositories {` +
		GradleMavenRepositoriesFile(repos) + `
  }
}

rootProject.name = 'demo'
`

	if err := gradleFile([]byte(content), "settings"); err != nil {
		return err
	}
	return nil
}

func GradleCleanBuild(pass string) (int, error) {
	if err := os.Setenv("YAAM_PASS", pass); err != nil {
		return 1, err
	}

	if err := os.Setenv("GRADLE_USER_HOME", testDirGradle+time.Now().Format("20060102150405111")); err != nil {
		return 1, err
	}
	cmd := exec.Command("bash", "-c", "./gradlew clean build --no-daemon")
	cmd.Dir = gradleHomeDemoProject
	co, err := cmd.CombinedOutput()
	if err != nil {
		return cmd.ProcessState.ExitCode(), fmt.Errorf(cmdExitErrMsg, string(co), err)
	}

	return 0, nil
}

func GradleBuildFile(repos []string, repoPublish string) error {
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
		GradleMavenRepositoriesFile(repos) + `
}

` + GradlePublishingFile(repoPublish) + `

dependencies {
  implementation 'org.springframework.boot:spring-boot-starter'
  testImplementation 'org.springframework.boot:spring-boot-starter-test'
}

tasks.named('test') {
  useJUnitPlatform()
}
`

	if err := gradleFile([]byte(content), "build"); err != nil {
		return err
	}
	return nil
}
