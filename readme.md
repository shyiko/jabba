# jabba ![](https://img.shields.io/badge/latest-0.3.0-green.svg)

![jabba-the-hutt](https://cloud.githubusercontent.com/assets/370176/13943697/e6098ed0-efbb-11e5-9630-3ff0d0d0403d.jpg)

Java Version Manager inspired by [nvm](https://github.com/creationix/nvm) (Node.js). You might have come across similar beauty
in Go ([gvm](https://github.com/moovweb/gvm)) or Ruby ([rvm](https://rvm.io)).

Supports installation of [Oracle JDK](http://www.oracle.com/technetwork/java/javase/archive-139210.html) (default), 
[Zulu OpenJDK](http://zulu.org/) (since 0.3.0) and from custom URLs.

Tested on Mac OS X and Linux. Windows support is coming in [#1](https://github.com/shyiko/jabba/issues/1).

It's written in [Go](https://golang.org/) to make maintenance easier (significantly shorter, easier to understand and less prone to errors 
compared to pure shell implementation). Plus it enables us to support Windows natively (no Cygwin) without rewriting 
the whole thing in PowerShell or whatever. 

The goal is to provide unified pain-free experience of installing (and switching between different versions of) JDK.

> `jabba` has single responsibility - managing different versions of JDK. Maven/Gradle/SBT/... are out of scope (for those use
[mvnw](https://github.com/shyiko/mvnw)/[gradlew](https://docs.gradle.org/current/userguide/gradle_wrapper.html)/[sbt-launcher](http://www.scala-sbt.org/0.13/docs/Manual-Installation.html)/...).
 
## Installation

> (use the same command to upgrade)

```sh
curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | bash && . ~/.jabba/jabba.sh
```

> In [fish](https://fishshell.com/) command looks a little bit different - 
`curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | bash; and . ~/.jabba/jabba.fish` 

## Usage

```sh
# install Oracle JDK
jabba install 1.8 # "jabba use 1.8" will be called automatically  

# install Zulu OpenJDK (since 0.3.0)
jabba install zulu@1.8.72

# install from custom URL (supported qualifiers: zip (since 0.3.0), tgz, dmg, bin)
jabba install 1.8.0-custom=tgz+http://example.com/distribution.tar.gz
jabba install 1.8.0-custom=zip+file:///opt/distribution.zip

# list all installed JDK's
jabba ls

# switch to a different version of JDK
jabba use 1.6.65

# list available JDK's
jabba ls-remote

# set default java version on shell (since 0.2.0)
jabba alias default 1.6.65
```

For more information see `jabba --help`.  

## Development

> PREREQUISITE: [go1.6](https://github.com/moovweb/gvm)

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

**Q**: How does it compare to "apt-get install" (and alike)?

A: Single command (`jabba install <version>`) works regardless of the OS (so you don't need to remember different ways to 
   install JDK on Mac OS X, Arch and two versions of Debian). Package name no longer matters, as well as, whether LTS is over
   or not.
   And you can install ANY version, not just the latest stable one. Eager to try upcoming 1.9.0 release? No need to wait -
   `jabba install 1.9.0-110`. How about specific build of 1.8? No problem - `jabba install 1.8.73`.

**Q**: What if I already have `java` installed?

A: `jabba` feels perfectly fine in the environment where `java` has already been installed using some other means. You 
 can continue using system JDK and switch to `jabba`-provided one only when needed (`jabba use 1.6.65`).

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

By using this software you agree to [Oracle Binary Code License Agreement for the Java SE Platform Products and JavaFX](http://www.oracle.com/technetwork/java/javase/terms/license/index.html)
and [Oracle Technology Network Early Adopter Development License Agreement](http://www.oracle.com/technetwork/licenses/ea-license-152003.html) (in case of EA releases) 
(... and Apple's Software License Agreement in case of "Java for OS X"). 

This software is for educational purposes only.  
Use it at your own risk. 
