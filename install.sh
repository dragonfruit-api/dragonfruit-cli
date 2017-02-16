#!/bin/sh
echo "\033[31;1mHI!\033[0m"
echo "\033[1mThis script will install Dragonfruit\033[0m."
TDIR="._tmp"
mkdir -p $TDIR
HOMEDIR=$(pwd)
cd $HOMEDIR

echo "\033[31;1mFirst, we'll uninstall any old versions of Dragonfruit you might have.\033[0m."
brew uninstall dragonfruit

# set up dragonfruit
echo "\n\n"
cd $TDIR
cd "homebrew/dragonfruit"
brew install "dragonfruit.rb"

cd $HOMEDIR
rm -rf $TDIR