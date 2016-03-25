# jabba

![jabba-the-hutt](https://cloud.githubusercontent.com/assets/370176/13943697/e6098ed0-efbb-11e5-9630-3ff0d0d0403d.jpg)

Java Version Manager inspired by [nvm](https://github.com/creationix/nvm) (Node.js). You might have come across similar beauty
in Go ([gvm](https://github.com/moovweb/gvm)) or Ruby ([rvm](https://rvm.io)).

Tested <sup>at this point **very** lightly</sup> on Mac OS X and Linux. Windows support is coming in [#1](https://github.com/shyiko/jabba/issues/1).

It's written in [Go](https://golang.org/) to make maintenance easier (significantly shorter, easier to understand and less prone to errors 
compared to pure shell implementation). Plus it enables us to support Windows natively (no Cygwin) without rewriting 
the whole thing in PowerShell or whatever. 

The goal is to provide unified pain-free experience of installing (and switching between different versions of) JDK.

How does it compare to "apt-get install" (and alike)?  
Single command (`jabba install <version>`) works regardless of the OS (so you don't need to remember different ways to 
install JDK on Mac OS X, Arch and two versions of Debian). Package name no longer matters, as well as, whether LTS is over
or not.
And you can install ANY version, not just the latest stable one. Wanna try upcoming 1.9.0 release? No need to wait -
`jabba install 1.9.0-110`. How about specific build of 1.8? No problem - `jabba install 1.8.73`.

> `jabba` has single responsibility - managing different versions of JDK. Maven/Gradle/SBT/... are out of scope (for those use
[mvnw](https://github.com/shyiko/mvnw)/[gradlew](https://docs.gradle.org/current/userguide/gradle_wrapper.html)/[sbt-launcher](http://www.scala-sbt.org/0.13/docs/Manual-Installation.html)/...).
 
## Installation

> (use the same command to upgrade)

```sh
curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | bash && . ~/.jabba/jabba.sh
```   

> Installer creates `~/.jabba` and adds initialization code to ~/.bashrc (and ~/.bash_profile, ~/.zshrc, ~/.profile) 
(provided they exist).

## Usage

```sh
# install particular version of jdk
jabba install 1.8 # "jabba use 1.8" will be called automatically  

# list all installed jdk's
jabba ls

# switch to a different version of jdk
jabba use 1.6.65

# list available jdk's
jabba ls-remote

# set default java version on shell (available since 0.2.0)
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

Q: What if I already have `java` installed?

A: `jabba` feels perfectly fine in the environment where `java` has already been installed using some other means. You 
 can continue using system JDK and switch to `jabba`-provided one only when needed (`jabba use 1.6.65`).

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

By using this software you agree to [Oracle Binary Code License Agreement for the Java SE Platform Products and JavaFX](http://www.oracle.com/technetwork/java/javase/terms/license/index.html)
and [Oracle Technology Network Early Adopter Development License Agreement](http://www.oracle.com/technetwork/licenses/ea-license-152003.html) (in case of EA releases) 
(... and Apple's Software License Agreement in case of "Java for OS X"). 

This software is for educational purposes only.  
Use it at your own risk. 
