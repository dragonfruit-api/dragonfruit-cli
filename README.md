# dragonfruit-cli

This is a command-line wrapper for Dragonfruit, a data-model and API prototyping tool, using CouchDB as a backend data store.

## Installation
### Mac
In a terminal, run `install.sh` from the root directory of this folder.

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
 * [Mercurial](https://www.mercurial-scm.org/)
 * [CouchDB](http://couchdb.apache.org/)
 * [Go](http://golang.org)

### Set up your paths:
 1. Set your `GOPATH` variable to the directory that you want your source files go to. (For example if you are working in `/tmp/dragonfruit` do `export GOPATH=/tmp/dragonfruit`.)
 2. Add `GOPATH/bin` to your `PATH` variable. `export PATH = $PATH:$GOPATH/bin`
 
### Install Govendor and get the source code
 1. Get `govendor` to handle dependencies: `go get -u github.com/kardianos/govendor` (**Note:** Don't put `git://` or `http://` in front of `github.com` here.)
 2. Get the source code: `go get github.com/dragonfruit-api/dragonfruit-cli` 
 3. Change to the `dragonfruit-cli` directory and then synchronize the dependencies with `govendor sync`.
 4. Install the binary: `go install`.  The binary will now be in your `$GOPATH/bin` folder and will be called `dragonfruit-cli`.
 5. You can now copy the binary somewhere.  If you want to rename it to `dragonfruit`, have a ball. You could also save it to `/usr/local/bin` or something.
 
## Initialize the database and copy the config files 
 1. Type `couchdb -b` to start CouchDB
 2. `curl -X PUT http://localhost:5984/swagger_docs`: This creates the swagger doc database
 3. Copy `/etc/dragonfruit.conf` to `/usr/local/etc/dragonfruit.conf`.
