#!/bin/bash

set -xe
export GOROOT=/opt/gae/google-cloud-sdk/platform/google_appengine/goroot-1.6
go=$GOROOT/bin/goapp

cd `dirname $0`
cd ..
src=`pwd`
cd ../../../../
# must be absolute fo 'go get' to work!
export GOPATH=`pwd`

$go get golang.org/x/text/language 
$go get -u golang.org/x/text/cmd/gotext
#cd "$src"
#export PATH="$src/_scripts:$PATH"
#"$GOPATH/bin/gotext" extract

