# dragonfruit-cli

This is a command-line wrapper for Dragonfruit, our data-model and API prototyping tool.  

## Installation
### Mac
In a terminal, run `install.sh` from the root directory of this folder.

**Note:** This requires Homebrew (http://brew.sh).  If you don't have Homebrew, you can install it along with a bunch of other useful tools with the [IDEO Building Blocks project](https://github.com/ideo/building-blocks).  

### PC
See *Building from Source* for now.

### Linux 
See *Building from Source* for now.
 
#### 

## Command line options


## Configuration

## Building from source
![](http://i.imgur.com/GNywZoP.gif)
### Prerequisites:
You will need the following:
 * Mercurial/hg
 * CouchDB
 * Go

### Set up your paths:
 1. Set your `GOPATH` variable to the directory that you want your source files go to. (For example if you are working in `/tmp/dragonfruit` do `export GOPATH=/tmp/dragonfruit`.)
 2. Add `GOPATH/bin` to your `PATH` variable. `export PATH = $PATH:$GOPATH/bin`
 
### Install Godep and get the source code
 1. `go get github.com/tools/godep` (**Note:** Don't put `git://` or `http://` in front of `github.com` here.)
 
 
