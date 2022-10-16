# YAAM

[![CI](https://github.com/030/yaam/workflows/Go/badge.svg?event=push)](https://github.com/030/yaam/actions?query=workflow%3AGo)
[![GoDoc Widget]][godoc]
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/030/yaam)
[![Go Report Card](https://goreportcard.com/badge/github.com/030/yaam)](https://goreportcard.com/report/github.com/030/yaam)
[![StackOverflow SE Questions](https://img.shields.io/stackexchange/stackoverflow/t/yaam.svg?logo=stackoverflow)](https://stackoverflow.com/tags/yaam)
[![DevOps SE Questions](https://img.shields.io/stackexchange/devops/t/yaam.svg?logo=stackexchange)](https://devops.stackexchange.com/tags/yaam)
[![ServerFault SE Questions](https://img.shields.io/stackexchange/serverfault/t/yaam.svg?logo=serverfault)](https://serverfault.com/tags/yaam)
![Docker Pulls](https://img.shields.io/docker/pulls/utrecht/yaam.svg)
[![yaam on stackoverflow](https://img.shields.io/badge/stackoverflow-community-orange.svg?longCache=true&logo=stackoverflow)](https://stackoverflow.com/tags/yaam)
![Issues](https://img.shields.io/github/issues-raw/030/yaam.svg)
![Pull requests](https://img.shields.io/github/issues-pr-raw/030/yaam.svg)
![Total downloads](https://img.shields.io/github/downloads/030/yaam/total.svg)
![GitHub forks](https://img.shields.io/github/forks/030/yaam?label=fork&style=plastic)
![GitHub watchers](https://img.shields.io/github/watchers/030/yaam?style=plastic)
![GitHub stars](https://img.shields.io/github/stars/030/yaam?style=plastic)
![License](https://img.shields.io/github/license/030/yaam.svg)
![Repository Size](https://img.shields.io/github/repo-size/030/yaam.svg)
![Contributors](https://img.shields.io/github/contributors/030/yaam.svg)
![Commit activity](https://img.shields.io/github/commit-activity/m/030/yaam.svg)
![Last commit](https://img.shields.io/github/last-commit/030/yaam.svg)
![Release date](https://img.shields.io/github/release-date/030/yaam.svg)
![Latest Production Release Version](https://img.shields.io/github/release/030/yaam.svg)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=bugs)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=code_smells)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=coverage)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=ncloc)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=alert_status)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=security_rating)](https://sonarcloud.io/dashboard?id=030_yaam)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=030_yaam&metric=sqale_index)](https://sonarcloud.io/dashboard?id=030_yaam)
[![codecov](https://codecov.io/gh/030/yaam/branch/main/graph/badge.svg)](https://codecov.io/gh/030/yaam)
[![BCH compliance](https://bettercodehub.com/edge/badge/030/yaam?branch=main)](https://bettercodehub.com/results/030/yaam)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-web.svg)](https://golangci.com/r/github.com/030/yaam)
[![codebeat badge](https://codebeat.co/badges/af6b1a01-df2c-40e7-bfb1-13ec0bb90087)](https://codebeat.co/projects/github-com-030-yaam-main)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)

[godoc]: https://godoc.org/github.com/030/yaam
[godoc widget]: https://godoc.org/github.com/030/yaam?status.svg

Although there are many artifact managers, like Artifactory, Nexus3 and
Verdaccio, they are either monoliths, consume a lot of resources
(memory and CPU), lack Infrastructure as Code (IaC) or do not support all kind
of artifact types. Yet Another Artifact Manager (YAAM):

- is an artifact manager like Artifactory, Nexus3 or Verdaccio.
- enforces IaC.
- has no UI.
- does not have a database.
- scales horizontally.
- supports downloading and publication of Apt, Generic, Maven and NPM
  artefacts, preserves NPM and Maven packages from public repositories and
  unifies Maven repositories.

## Quickstart

Create a directory and change the permissions to ensure that YAAM can store
artifacts:

```bash
mkdir ~/.yaam/repositories
sudo chown 9999 -R ~/.yaam/repositories
```

Configure YAAM by creating a `~/.yaam/config.yml` with the following content:

```bash
---
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
      url: https://some-nexus/repository/some-repo/
      user: some-user
      pass: some-pass
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
    - 3rdparty-npm
```

Start YAAM:

```bash
docker run \
  -e YAAM_DEBUG=false \
  -e YAAM_USER=hello \
  -e YAAM_PASS=world \
  --rm \
  --name=yaam \
  -it \
  -v /home/${USER}/.yaam:/opt/yaam/.yaam \
  -p 25213:25213 utrecht/yaam:v0.5.0
```

Once YAAM has been started, configure a project to ensure that artifacts will
be downloaded from this artifact manager.

### Apt

sudo vim /etc/apt/auth.conf.d/hello.conf

```bash
machine http://localhost
login hello
password world
```

sudo vim /etc/apt/sources.list

```bash
deb http://localhost:25213/apt/3rdparty-ubuntu-nl-archive/ focal main restricted
```

Preserve the artifacts:

```bash
sudo apt-get update
```

### Generic

Upload:

```bash
curl -X POST -u hello:world \
http://yaam.some-domain/generic/something/world4.iso \
--data-binary @/home/${USER}/Downloads/ubuntu-22.04.1-desktop-amd64.iso
```

Troubleshooting:

```bash
413 Request Entity Too Large
```

add:

```bash
data:
  proxy-body-size: 5G
```

and restart the controller pod.

Verify in the `/etc/nginx/nginx.conf` file that the `client_max_body_size` has
been increased to 5G.

Download:

```bash
curl -u hello:world http://yaam.some-domain/generic/something/world6.iso \
-o /tmp/world6.iso
```

### Gradle

Adjust the `build.gradle` and/or `settings.gradle`:

```bash
repositories {
  maven {
    allowInsecureProtocol true
    url 'http://localhost:25213/maven/releases/'
    authentication {
      basic(BasicAuthentication)
    }
    credentials {
      username "hello"
      password "world"
    }
  }
  maven {
    allowInsecureProtocol true
    url 'http://localhost:25213/maven/3rdparty-maven/'
    authentication {
      basic(BasicAuthentication)
    }
    credentials {
      username "hello"
      password "world"
    }
  }
  maven {
    allowInsecureProtocol true
    url 'http://localhost:25213/maven/3rdparty-maven-gradle-plugins/'
    authentication {
      basic(BasicAuthentication)
    }
    credentials {
      username "hello"
      password "world"
    }
  }
}
```

Publish:

```bash
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

  repositories {
    maven {
        allowInsecureProtocol true
        url 'http://localhost:25213/maven/releases/'
        authentication {
          basic(BasicAuthentication)
        }
        credentials {
          username "hello"
          password "world"
        }
    }
  }
}
```

Preserve the artifacts:

```bash
./gradlew clean
```

or publish them:

```bash
./gradlew publish
```

### NPM

Create a `.npmrc` file in the directory of a particular NPM project:

```bash
registry=http://localhost:25213/npm/3rdparty-npm/
always-auth=true
_auth=aGVsbG86d29ybGQ=
cache=/tmp/some-yaam-repo/npm/cache20220914120431999
```

Note: the `_auth` key should be populated with the output of:
`echo -n 'someuser:somepass' | openssl base64`.

```bash
npm i -d
```

## Run

Next to docker, one could also use a binary or K8s or OpenShift to run YAAM:

- [Binary.](docs/start/BINARY.md)
- [K8s/OpenShift.](docs/start/K8SOPENSHIFT.md)

## Other

- [Background.](docs/other/BACKGROUND.md)
- [Maven unify.](docs/other/MAVEN.md)
