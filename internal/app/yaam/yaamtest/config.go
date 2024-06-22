package yaamtest

import (
	"os"
	"path/filepath"
	"time"
)

const (
	Testdata              = "test/testdata"
	TestdataCwd           = "../../" + Testdata
	gradleHomeDemoProject = TestdataCwd + "/gradle/demo"
	npmHomeDemoProject    = TestdataCwd + "/npm/demo"

	conf = `---
caches:
  apt:
    3rdparty-ubuntu-nl-archive:
      url: http://nl.archive.ubuntu.com/ubuntu/
  maven:
    3rdparty-maven:
      url: https://repo.maven.apache.org/maven2/
    3rdparty-maven-gradle-plugins:
      url: https://plugins.gradle.org/m2/
    3rdparty-maven-spring:
      url: https://repo.spring.io/release/
    other-nexus-repo-releases:
      url: x
      user: y
      pass: z
  npm:
    3rdparty-npm:
      url: https://registry.npmjs.org/
groups:
  maven:
    hello:
      - maven/releases
      - maven/3rdparty-maven
      - maven/3rdparty-maven-gradle-plugins
      - maven/3rdparty-maven-spring
publications:
  generic:
    - something
  maven:
    - releases
  npm:
    - 3rdparty-npm`
	cmdExitErrMsg = "%v, err: '%v'"

	dir              = "/tmp/yaam"
	DirIso           = dir + "/ubuntu.iso"
	DirIsoDownloaded = dir + "/downloaded-ubuntu.iso"

	testDirGradle = dir + "/gradle"
	testDirNpm    = dir + "/npm"
)

func Config() error {
	m := make(map[string]string)
	m["HOME"] = filepath.Join(dir, "test"+time.Now().Format("20060102150405111"))
	m["DEBUG"] = "true"
	m["USER"] = "hello"
	m["PASS"] = "world"
	for k, v := range m {
		if err := os.Setenv("YAAM_"+k, v); err != nil {
			return err
		}
	}

	dir := os.Getenv("YAAM_HOME")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(conf), 0o600); err != nil {
		return err
	}

	return nil
}
