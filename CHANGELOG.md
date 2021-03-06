# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.5.0] - 2020-10-07 

- Bump api/core/v2 version to 2.3.0 (adds validation for "prometheus_text" in 
  "event.check.output_format_metrics")

## [0.4.0] - 2020-09-16

- Write tags as JSON objects instead of an array of objects to improve Postgres 
  JSON SQL queries

## [0.3.0] - 2020-07-24 

- Update .goreleaser.yml to produce sha512 hashes for release artifacts 

## [0.2.0] - 2020-07-24

- Add support for configurable sslmodes (via `--sslmode` or `$TIMESCALEDB_SSLMODE`)
- Remove TravisCI in favor of Github Actions
- Update Sensu Go and SDK dependencies with the correct modules

## [0.1.0] - 2020-01-07

### Added
- Initial release
