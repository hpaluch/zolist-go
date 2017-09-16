#!/bin/bash

# reformats my source
gofmt=/opt/gae/google-cloud-sdk/platform/google_appengine/goroot-1.8/bin/gofmt
cd `dirname $0`
cd ..
$gofmt -d *.go 
echo -n "Reformat source [y/N]?"
read ans
case "$ans" in
  y*)
    $gofmt -w *.go
    echo "Reformatted"
    exit 0
    ;;
  *)
    echo "Canceled"
    exit 1
   ;; 
esac
