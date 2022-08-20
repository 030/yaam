# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- debug logging option.
- line number to logging.

### Changed

- empty YAAM version number in docker container to the real actual version.
- only show log if file is downloaded.

## [0.2.0] - 2022-08-20

### Added

- capability to upload and cache artifacts.
- docker image.
- basic authentication.
- CI.
- k8s-openshift deployment.

### Changed

- Golang version to 1.19.
- port to 25213.

## [0.1.0] - 2022-07-17

### Added

- Cache all Maven2 artifacts that are required by a gradle project locally.

[Unreleased]: https://github.com/030/yaam/compare/0.2.0...HEAD
[0.2.0]: https://github.com/030/yaam/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/030/yaam/releases/tag/0.1.0
