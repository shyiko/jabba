# jabba ![Latest Version](https://img.shields.io/badge/latest-0.11.2-blue.svg) [![Build Status](https://travis-ci.org/shyiko/jabba.svg?branch=master)](https://travis-ci.org/shyiko/jabba)

Java Version Manager inspired by [nvm](https://github.com/creationix/nvm) (Node.js). Written in Go.

The goal is to provide unified pain-free experience of **installing** (and **switching** between different versions of) JDK regardless of
the OS (macOS, Linux x86/x86_64/ARMv7+, Windows x86_64). 

`jabba install`
- [Oracle JDK](http://www.oracle.com/technetwork/java/javase/archive-139210.html) (latest-version only)
- [Oracle Server JRE](http://www.oracle.com/technetwork/java/javase/downloads/server-jre8-downloads-2133154.html) (latest-version only), 
- [Adopt OpenJDK](https://adoptopenjdk.net/) <sup>(jabba >=0.8.0 is required)</sup> 
  - Hotspot 
  - [Eclipse OpenJ9](https://www.eclipse.org/openj9/oj9_faq.html)
- [Zulu OpenJDK](http://zulu.org/) <sup>(jabba >=0.3.0 is required)</sup>
- [IBM SDK, Java Technology Edition](https://developer.ibm.com/javasdk/) <sup>(jabba >=0.6.0 is required)</sup> 
- [GraalVM CE](https://www.graalvm.org/)
- [OpenJDK](http://openjdk.java.net/)
- [OpenJDK Reference Implementation](http://openjdk.java.net/)
- [OpenJDK with Shenandoah GC](https://wiki.openjdk.java.net/display/shenandoah/Main) <sup>(jabba >=0.10.0 is required)
- [Liberica JDK](https://bell-sw.com/)
- [Amazon Corretto](https://aws.amazon.com/corretto/)
</sup>

... and from custom URLs.

## Installation

#### macOS / Linux

> (in bash/zsh/...)

```sh
curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | bash && . ~/.jabba/jabba.sh
```

> (use the same command to upgrade)

The script modifies common shell rc files by default. To skip these provide the `--skip-rc` flag to `install.sh` like so:

```sh
curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | bash -s -- --skip-rc && . ~/.jabba/jabba.sh
```

Make sure to source `jabba.sh` in your environment if you skip it:

```sh
[ -s "$JABBA_HOME/jabba.sh" ] && source "$JABBA_HOME/jabba.sh"
```

> In [fish](https://fishshell.com/) command looks a little bit different - 
`curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | bash; and . ~/.jabba/jabba.fish` 

> If you don't have `curl` installed - replace `curl -sL` with `wget -qO-`.

> If you are behind a proxy see -
[curl](https://curl.haxx.se/docs/manpage.html#ENVIRONMENT) / 
[wget](https://www.gnu.org/software/wget/manual/wget.html#Proxies) manpage. 
Usually simple `http_proxy=http://proxy-server:port https_proxy=http://proxy-server:port curl -sL ...` is enough. 

**NOTE**: The brew package is currently broken. We are working on a fix.

#### Docker

While you can use the same snippet as above, chances are you don't want jabba binary & shell 
integration script(s) to be included in the final Docker image, all you want is a JDK. Here is the `Dockerfile` showing how this can be done:

```dockerfile
FROM buildpack-deps:jessie-curl

RUN curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | \
    JABBA_COMMAND="install 1.15.0 -o /jdk" bash

ENV JAVA_HOME /jdk
ENV PATH $JAVA_HOME/bin:$PATH
```

> (when `JABBA_COMMAND` env variable is set `install.sh` downloads `jabba` binary, executes specified command and then deletes the binary)

```sh
$ docker build -t <image_name>:<image_tag> .
$ docker run -it --rm <image_name>:<image_tag> java -version

java version "1.15.0....
```

#### Windows 10

> (in powershell)

```powershell
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Invoke-Expression (
  Invoke-WebRequest https://github.com/shyiko/jabba/raw/master/install.ps1 -UseBasicParsing
).Content
```

> (use the same command to upgrade)

## Usage

```sh
# list available JDK's
jabba ls-remote
# you can use any valid semver range to narrow down the list
jabba ls-remote zulu@~1.8.60
jabba ls-remote "*@>=1.6.45 <1.9" --latest=minor

# install Oracle JDK
jabba install 1.15.0
# install Oracle Server JRE
jabba install sjre@1.8  
# install Adopt OpenJDK (Hotspot)
jabba install adopt@1.8-0
# install Adopt OpenJDK (Eclipse OpenJ9)
jabba install adopt-openj9@1.9-0
# install Zulu OpenJDK
jabba install zulu@1.8
jabba install zulu@~1.8.144 # same as "zulu@>=1.8.144 <1.9" 
# install IBM SDK, Java Technology Edition
jabba install ibm@1.8
# install GraalVM CE
jabba install graalvm@1.0-0
# install OpenJDK
jabba install openjdk@1.10-0
# install OpenJDK with Shenandoah GC
jabba install openjdk-shenandoah@1.10-0

# install from custom URL
# (supported qualifiers: zip (since 0.3.0), tgz, tgx (since 0.10.0), dmg, bin, exe)
jabba install 1.8.0-custom=tgz+http://example.com/distribution.tar.gz
jabba install 1.8.0-custom=tgx+http://example.com/distribution.tar.xz
jabba install 1.8.0-custom=zip+file:///opt/distribution.zip

# uninstall JDK
jabba uninstall zulu@1.6.77

# link system JDK
jabba link system@1.8.72 /Library/Java/JavaVirtualMachines/jdk1.8.0_72.jdk

# list all installed JDK's
jabba ls

# switch to a different version of JDK (it must be already `install`ed)
jabba use adopt@1.8
jabba use zulu@~1.6.97

echo "1.8" > .jabbarc
# switch to the JDK specified in .jabbarc (since 0.5.0)
jabba use

# set default java version on shell (since 0.2.0)
# this version will automatically be "jabba use"d every time you open up a new terminal
jabba alias default 1.8
```

> `.jabbarc` has to be a valid YAML file. JDK version can be specified as `jdk: 1.8` or simply as `1.8` 
(same as `~1.8`, `1.8.x` `">=1.8.0 <1.9.0"` (mind the quotes)).

> jsyk: **jabba** keeps everything under `~/.jabba` (on Linux/Mac OS X) / `%USERPROFILE%/.jabba` (on Windows). If at any point of time you decide to uninstall **jabba** - just remove this directory. 

For more information see `jabba --help`.  

## Development

> PREREQUISITE: [go1.8](https://github.com/moovweb/gvm)

```sh
git clone https://github.com/shyiko/jabba $GOPATH/src/github.com/shyiko/jabba 
cd $GOPATH/src/github.com/shyiko/jabba 
make fetch

go run jabba.go

# to test a change
make test # or "test-coverage" if you want to get a coverage breakdown

# to make a build
make build # or "build-release" (latter is cross-compiling jabba to different OSs/ARCHs)   
```

## FAQ

**Q**: What if I already have `java` installed?

A: It's fine. You can switch between system JDK and `jabba`-provided one whenever you feel like it (`jabba use ...` / `jabba deactivate`). 
They are not gonna conflict with each other.

**Q**: How do I switch `java` globally?

A: **jabba** doesn't have this functionality built-in because the exact way varies greatly between the operation systems and usually 
involves elevated permissions. But. Here are the snippets that <u>should</u> work:    

* Windows

> (in powershell as administrator)

```
# select jdk
jabba use ...

# modify global PATH & JAVA_HOME
$envRegKey = [Microsoft.Win32.Registry]::LocalMachine.OpenSubKey('SYSTEM\CurrentControlSet\Control\Session Manager\Environment', $true)
$envPath=$envRegKey.GetValue('Path', $null, "DoNotExpandEnvironmentNames").replace('%JAVA_HOME%\bin;', '')
[Environment]::SetEnvironmentVariable('JAVA_HOME', "$(jabba which $(jabba current))", 'Machine')
[Environment]::SetEnvironmentVariable('PATH', "%JAVA_HOME%\bin;$envPath", 'Machine')
```

* Linux

> (tested on Debian/Ubuntu)

```
# select jdk
jabba use ...

sudo update-alternatives --install /usr/bin/java java ${JAVA_HOME%*/}/bin/java 20000
sudo update-alternatives --install /usr/bin/javac javac ${JAVA_HOME%*/}/bin/javac 20000
```

> To switch between multiple GLOBAL alternatives use `sudo update-alternatives --config java`.

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

By using this software you agree to
- [Oracle Binary Code License Agreement for the Java SE Platform Products and JavaFX](http://www.oracle.com/technetwork/java/javase/terms/license/index.html)
- [Oracle Technology Network Early Adopter Development License Agreement](http://www.oracle.com/technetwork/licenses/ea-license-noexhibits-1938914.html) in case of EA releases
- Apple's Software License Agreement in case of "Java for OS X"
- [International License Agreement for Non-Warranted Programs](http://www14.software.ibm.com/cgi-bin/weblap/lap.pl?la_formnum=&li_formnum=L-PMAA-A3Z8P2&l=en) in case of IBM SDK, Java Technology Edition.

This software is for educational purposes only.  
Use it at your own risk. 
