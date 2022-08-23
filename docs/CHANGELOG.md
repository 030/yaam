# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.2.1] - 2022-08-23

### Added

- debug logging option.
- line number to logging.
- resource requests and limits to k8s and openshift deployment.
- health and readiness check to k8s deployment.

### Changed

- empty YAAM version number in docker container to the real actual version.
- only show log if file is downloaded.
- volume templates to emptyDir in k8s and openshift deployment.
- the origin of various functions to pkg to ensure they can be imported as a
  library by external resources.

## [v0.2.0] - 2022-08-20

### Added

- capability to upload and cache artifacts.
- docker image.
- basic authentication.
- CI.
- k8s-openshift deployment.

### Changed

- Golang version to 1.19.
- port to 25213.

## [v0.1.0] - 2022-07-17

### Added

- Cache all Maven2 artifacts that are required by a gradle project locally.

[Unreleased]: https://github.com/030/yaam/compare/v0.2.1...HEAD
[v0.2.1]: https://github.com/030/yaam/compare/v0.2.0...0.2.1
[v0.2.0]: https://github.com/030/yaam/compare/v0.1.0...0.2.0
[v0.1.0]: https://github.com/030/yaam/releases/tag/v0.1.0
