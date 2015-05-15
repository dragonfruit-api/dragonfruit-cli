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
 2. `godep get github.com/ideo/dragonfruit-cli` This will install a binary called `dragonfruit-cli` into the `GOPATH/bin` directory
 3. You can now copy the binary somewhere.  If you want to rename it to `dragonfruit`, have a ball. You could also save it to `/usr/local/bin` or something.
 
## Initialize the database and copy the config files 
 1. Type `couchdb -b` to start CouchDB
 2. `curl -X PUT http://localhost:5984/swagger_docs`: This creates the swagger doc database
 3. Copy the `dragonfruit.conf` file in this folder to `/etc/`
