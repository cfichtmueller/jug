# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2024-02-11

### Added

- Context provides access to HTTP request and response
- `ClientIP()` and `RemoteIP()` utility functions on Context

### Changed

- removed `String` designators in validation functions that validate string values

### Deprecated

- `RequireString...` validation methods

## [0.1.0] - 2023-09-27

### Added

- Engine
- Router
- Context
- Validatable interface
- Validator
- ResponseStatusError
- Map utilities

[unreleased]: https://github.com/cfichtmueller/jug/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/cfichtmueller/jug/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/cfichtmueller/jug/releases/tag/v0.1.0
