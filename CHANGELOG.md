# Changelog
All notable changes to this project will be documented in this file.  
This project adheres to [Semantic Versioning](http://semver.org/).

## [0.11.0](https://github.com/shyiko/jabba/compare/0.10.1...0.11.0) - 2018-10-11

### Added
- Flags to filter remote versions by `--os` and `--arch` (e.g. `jabba ls-remote --os=windows`) ([#240](https://github.com/shyiko/jabba/issues/240))

### Fixed
- File tree normalization  
(Adopt OpenJDK 11 distributions contain meta file at the root level causing untgz --strip to fail) ([#311](https://github.com/shyiko/jabba/issues/311)). 

## [0.10.1](https://github.com/shyiko/jabba/compare/0.10.0...0.10.1) - 2018-05-07

### Fixed
- `jabba install <semver>` not checking whether JDK is already installed. 

## [0.10.0](https://github.com/shyiko/jabba/compare/0.9.6...0.10.0) - 2018-05-06

- [OpenJDK with Shenandoah GC](https://wiki.openjdk.java.net/display/shenandoah/Main) support ([#191](https://github.com/shyiko/jabba/issues/191))  
(e.g. `jabba install openjdk-shenandoah@1.9`).
- Ability to install JDK from `tar.xz` archives  
(e.g. `jabba install openjdk-shenandoah@1.9.0-220=tgx+file://$PWD/local-copy-of-openjdk-shenandoah-jdk9-b220-x86-release.tar.xz`).

## [0.9.6](https://github.com/shyiko/jabba/compare/0.9.5...0.9.6) - 2018-05-05

### Fixed
- Sporadic "open: permission denied" when installing from tgz/zip's ([#190](https://github.com/shyiko/jabba/issues/190)).  
Fix applied in 0.9.5 proved to be incomplete. 

## [0.9.5](https://github.com/shyiko/jabba/compare/0.9.4...0.9.5) - 2018-05-04

### Fixed
- Sporadic "open: permission denied" when installing from tgz/zip's ([#190](https://github.com/shyiko/jabba/issues/190)).

## [0.9.4](https://github.com/shyiko/jabba/compare/0.9.3...0.9.4) - 2018-05-01

### Fixed
- Installation from sources containing `Contents/Home` (macOS) (regression introduced in 0.9.3).

## [0.9.3](https://github.com/shyiko/jabba/compare/0.9.2...0.9.3) - 2018-04-30

### Fixed
- `Contents/Home` handling (macOS) ([#187](https://github.com/shyiko/jabba/issues/187)).

## [0.9.2](https://github.com/shyiko/jabba/compare/0.9.1...0.9.2) - 2017-11-18 

### Fixed
- `zip` & `tgz` stripping on Windows ([#116](https://github.com/shyiko/jabba/issues/116)).

## [0.9.1](https://github.com/shyiko/jabba/compare/0.9.0...0.9.1) - 2017-10-12 

### Fixed
- `tgz is not supported` when trying to install JDK from `tar.gz` on macOS & Windows. 

## [0.9.0](https://github.com/shyiko/jabba/compare/0.8.0...0.9.0) - 2017-09-19 

### Added
- Latest JDK / `default` alias (automatic) linking ([#6](https://github.com/shyiko/jabba/issues/6))

    ```sh
    $ ll ~/.jabba/jdk/
    lrwxrwxrwx  1 shyiko shyiko   30 Sep 19  2017  1.8 -> /home/shyiko/.jabba/jdk/1.8.144/
    drwxr-xr-x  8 shyiko shyiko 4096 Sep 19  2017  1.8.144/
    drwxr-xr-x  8 shyiko shyiko 4096 Sep 19  2017  1.8.141/
    lrwxrwxrwx  8 shyiko shyiko   30 Sep 19  2017  zulu@1.6 -> /home/shyiko/.jabba/jdk/zulu@1.6.97/
    drwxr-xr-x  8 shyiko shyiko 4096 Sep 19  2017  zulu@1.6.97/
    lrwxrwxrwx  1 shyiko shyiko   30 Sep 19  2017  default -> /home/shyiko/.jabba/jdk/1.8.144/
    ```

## [0.8.0](https://github.com/shyiko/jabba/compare/0.7.0...0.8.0) - 2017-09-19 

### Added
- [Adopt OpenJDK](https://adoptopenjdk.net/) support.
- `jabba ls <semver_range>` & `jabba ls-remote <semver_range>`.
- `jabba ls --latest=<major|minor|patch>` & `jabba ls-remote --latest=<major|minor|patch>`.

    ```sh
    $ jabba ls-remote "zulu@<1.9" --latest=minor
    zulu@1.8.144
    zulu@1.7.154
    zulu@1.6.97
    ```

- Ability to install JDK in a custom location (`jabba install -o /jdk/destination`)  
NOTE: any JDK installed in this way is considered to be unmanaged, i.e. not available to `jabba ls`, `jabba use`, etc. (unless `jabba link`ed).

### Changed
- semver library to [masterminds/semver](https://github.com/Masterminds/semver)  
(previously used library proved unreliable when given certain input (e.g. `>=1.6`)).

## [0.7.0](https://github.com/shyiko/jabba/compare/0.6.1...0.7.0) - 2017-05-12

### Added
* Ability to change the location of `~/.jabba` with `JABBA_HOME` env variable (e.g.
`curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | JABBA_HOME=/opt/jabba bash && . /opt/jabba/jabba.sh`)
* `--home` flag for `jabba which` (`jabba which --home <jdk_version>` returns `$JABBA_DIR/jdk/<jdk_version>/Contents/Home` on macOS and
`$JABBA_DIR/jdk/<jdk_version>` everywhere else)

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
