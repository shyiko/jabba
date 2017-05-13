# Changelog
All notable changes to this project will be documented in this file.  
This project adheres to [Semantic Versioning](http://semver.org/).

## [0.7.0](https://github.com/shyiko/jabba/compare/0.6.1...0.7.0) - 2017-05-12

### Added
* Ability to change the location of `~/.jabba` with `JABBA_HOME` env variable (e.g.
`curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | JABBA_HOME=/opt/jabba bash && . ~/.jabba/jabba.sh`)

### Changed
* `~` directory referencing inside shell integration scripts (path to home directory is now determined by `$PATH`).

### Fixed
* `jabba deactivate` escaping.
* `JAVA_HOME` propagation in Fish shell on CentOS 7.

## [0.6.1](https://github.com/shyiko/jabba/compare/0.6.0...0.6.1) - 2017-02-27

### Fixed
* `x509: certificate signed by unknown authority` while executing `jabba ls-remote` (macOS) ([#56](https://github.com/shyiko/jabba/issues/56)).

## [0.6.0](https://github.com/shyiko/jabba/compare/0.5.1...0.6.0) - 2016-12-10

### Added
* IBM SDK, Java Technology Edition support (e.g. `jabba install ibm@<version>`).

## [0.5.1](https://github.com/shyiko/jabba/compare/0.5.0...0.5.1) - 2016-11-28

### Fixed
* `link` command (it was missing `mkdir -p ...` call).

## [0.5.0](https://github.com/shyiko/jabba/compare/0.4.0...0.5.0) - 2016-11-11

### Added
* `.jabbarc` support. 

## [0.4.0](https://github.com/shyiko/jabba/compare/0.3.3...0.4.0) - 2016-07-22

### Added
* Windows 10 support.
* `link`/`unlink` commands.

### Fixed
* `bin+file://` handling (original file is now copied instead of being moved).

## [0.3.3](https://github.com/shyiko/jabba/compare/0.3.2...0.3.3) - 2016-03-30

### Fixed
* `dmg` handling when GNU tar is installed.

## [0.3.2](https://github.com/shyiko/jabba/compare/0.3.1...0.3.2) - 2016-03-26

### Fixed
* Zulu OpenJDK installation.

## [0.3.1](https://github.com/shyiko/jabba/compare/0.3.0...0.3.1) - 2016-03-26

### Fixed
* `current` (previously output was written to stderr instead of stdout).

## [0.3.0](https://github.com/shyiko/jabba/compare/0.2.0...0.3.0) - 2016-03-26

### Added
* Zulu OpenJDK support (e.g. `jabba install zulu@<version>`).
* Ability to install JDK from `zip` archives (in addition to already implemented `dmg`/`tar.gz`/`bin`).
* Support for custom registries (e.g. `JABBA_INDEX=https://github.com/shyiko/jabba/raw/master/index.json jabba install ...`). 

### Fixed
* `which <alias>`.

## [0.2.0](https://github.com/shyiko/jabba/compare/0.1.0...0.2.0) - 2016-03-24

### Added 
* `alias default`/`unalias default`, `which`, `deactivate` commands. 

## 0.1.0 - 2016-03-23
