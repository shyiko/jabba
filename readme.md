# jabba

![jabba-the-hutt](https://cloud.githubusercontent.com/assets/370176/13943697/e6098ed0-efbb-11e5-9630-3ff0d0d0403d.jpg)

Java Version Manager inspired by [nvm](https://github.com/creationix/nvm) (Node.js). You might have come across similar beauty
in Go ([gvm](https://github.com/moovweb/gvm)) or Ruby ([rvm](https://rvm.io)).

Tested <sup>at this point **very** lightly</sup> on Mac OS X and Linux. Windows support is coming in [#1]().

## Installation

```sh
curl -o- https://github.com/shyiko/jabba/raw/master/install.sh | bash && . ~/.jabba/jabba.sh
```   

> Installer creates `~/.jabba` and adds initialization code to ~/.bashrc (and ~/.bash_profile, ~/.zshrc, ~/.profile) 
(provided they exist).

## Usage

```
# install particular version of jdk
jabba install 1.8 # "jabba use 1.8" will be called automatically  

# list all installed jdk's
jabba ls

# switch to a different version of jdk
jabba use 1.6.65

# list available jdk's
jabba ls-remote
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

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

By using this software you agree to [Oracle Binary Code License Agreement for the Java SE Platform Products and JavaFX](http://www.oracle.com/technetwork/java/javase/terms/license/index.html)
and [Oracle Technology Network Early Adopter Development License Agreement](http://www.oracle.com/technetwork/licenses/ea-license-152003.html) (in case of EA releases) 
(... and Apple's Software License Agreement in case of "Java for OS X"). 

This software is for educational purposes only.  
Use it at your own risk. 
