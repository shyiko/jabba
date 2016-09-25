# jabba ![Latest Version](https://img.shields.io/badge/latest-0.4.0-blue.svg) [![Build Status](https://travis-ci.org/shyiko/jabba.svg?branch=master)](https://travis-ci.org/shyiko/jabba)

![jabba-the-hutt](https://cloud.githubusercontent.com/assets/370176/13943697/e6098ed0-efbb-11e5-9630-3ff0d0d0403d.jpg)

Java Version Manager inspired by [nvm](https://github.com/creationix/nvm) (Node.js). You might have come across similar beauty
in Go ([gvm](https://github.com/moovweb/gvm)) or Ruby ([rvm](https://rvm.io)).

Supports installation of [Oracle JDK](http://www.oracle.com/technetwork/java/javase/archive-139210.html) (default), 
[Zulu OpenJDK](http://zulu.org/) (since 0.3.0) and from custom URLs.

It's written in [Go](https://golang.org/) to make maintenance easier (significantly shorter, easier to understand and less prone to errors 
compared to pure shell implementation). Plus it enables us to support Windows natively (no need for Cygwin) without rewriting 
the whole thing in PowerShell or whatever. 

The goal is to provide unified pain-free experience of **installing** (and **switching** between different versions of) JDK regardless of
the OS. 

> **jabba** has a single responsibility - managing different versions of JDK. For an easy way to install Scala/Kotlin/Groovy (+ a lot more) see [SDKMAN][0]. 
SBT/Maven/Gradle should <u>ideally</u> be "fixed in place" by [sbt-launcher][1]/[mvnw][2]/[gradlew][3].
 
[0]: http://sdkman.io/  
[1]: http://www.scala-sbt.org/0.13/docs/Manual-Installation.html
[2]: https://github.com/shyiko/mvnw
[3]: https://docs.gradle.org/current/userguide/gradle_wrapper.html
 
## Installation

> (use the same command to upgrade)

* Linux/Mac OS X

> (in bash/zsh/...)

```sh
curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | bash && . ~/.jabba/jabba.sh
```

> In [fish](https://fishshell.com/) command looks a little bit different - 
`curl -sL https://github.com/shyiko/jabba/raw/master/install.sh | bash; and . ~/.jabba/jabba.fish` 

> If you don't have `curl` installed - replace `curl -sL` with `wget -qO-`.

> If you are behind a proxy see -
[curl](https://curl.haxx.se/docs/manpage.html#ENVIRONMENT) / 
[wget](https://www.gnu.org/software/wget/manual/wget.html#Proxies) manpage. 
Usually simple `http_proxy=http://proxy-server:port https_proxy=http://proxy-server:port curl -sL ...` is enough. 

* Windows 10

> (in powershell)

```powershell
Invoke-Expression (wget https://github.com/shyiko/jabba/raw/master/install.ps1 -UseBasicParsing).Content
```

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
# this version will automatically be "jabba use"d every time you open up a new terminal
jabba alias default 1.6.65
```

> jsyk: **jabba** keeps everything under `~/.jabba` (on Linux/Mac OS X) / `%USERPROFILE%/.jabba` (on Windows).

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

**Q**: What if I already have `java` installed?

A: It's fine. You can switch between system JDK and `jabba`-provided one whenever you feel like it (`jabba use ...` / `jabba deactivate`). 
They are not gonna conflict with each other.

**Q**: How do I switch `java` globally?

A: **jabba** doesn't have this functionality built-in because the exact way varies greatly between the operation systems and usually 
involves elevated permissions. But. Here are the snippets that <u>should</u> work:    

* Windows

> (in powershell as administrator)

<pre style="word-wrap: break-word;">
# select jdk
jabba use ...

# modify global PATH & JAVA_HOME
$envRegKey = [Microsoft.Win32.Registry]::LocalMachine.OpenSubKey('SYSTEM\CurrentControlSet\Control\Session Manager\Environment', $true)
$envPath=$envRegKey.GetValue('Path', $null, "DoNotExpandEnvironmentNames").replace('%JAVA_HOME%\bin;', '')
setx JAVA_HOME "$(jabba which $(jabba current))" /m
setx PATH "%JAVA_HOME%\bin;$envPath" /m
</pre>

* Linux

> (tested on Debian/Ubuntu)

<pre style="word-wrap: break-word;">
# select jdk
jabba use ...

sudo update-alternatives --install /usr/bin/java java ${JAVA_HOME%*/}/bin/java 20000
sudo update-alternatives --install /usr/bin/javac javac ${JAVA_HOME%*/}/bin/javac 20000
</pre>

> To switch between multiple GLOBAL alternatives use `sudo update-alternatives --config java`.

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

By using this software you agree to [Oracle Binary Code License Agreement for the Java SE Platform Products and JavaFX](http://www.oracle.com/technetwork/java/javase/terms/license/index.html)
and [Oracle Technology Network Early Adopter Development License Agreement](http://www.oracle.com/technetwork/licenses/ea-license-152003.html) (in case of EA releases) 
(... and Apple's Software License Agreement in case of "Java for OS X"). 

This software is for educational purposes only.  
Use it at your own risk. 
