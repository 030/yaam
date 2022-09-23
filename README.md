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
- supports downloading and publication of Generic, Maven and NPM artefacts,
  preserves NPM and Maven packages from public repositories and unifies Maven
  repositories.

## Configuration

### General

- [Base.](docs/config/BASE.md)

### Artifact types

- [Generic.](docs/config/GENERIC.md)
- [Maven.](docs/config/MAVEN.md)
- [NPM.](docs/config/NPM.md)

## Start

- [Binary.](docs/start/BINARY.md)
- [Docker.](docs/start/DOCKER.md)
- [K8s/OpenShift.](docs/start/K8SOPENSHIFT.md)

## Capabilities

### Publish

- [Generic.](docs/publish/GENERIC.md)
- [Maven.](docs/publish/MAVEN.md)

### Preserve

- [Maven.](docs/preserve/MAVEN.md)
- [NPM.](docs/preserve/NPM.md)

### Unify

- [Maven.](docs/unify/MAVEN.md)

## Other

- [Background.](docs/other/BACKGROUND.md)
