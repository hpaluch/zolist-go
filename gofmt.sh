#!/bin/bash

# reformats my source
gofmt=/opt/gae/google-cloud-sdk/platform/google_appengine/goroot-1.8/bin/gofmt
set -x
$gofmt -d *.go 
echo -n "Reformat source [y/N]?"
read ans
case "$ans" in
  y*)
    $gofmt -w *.go
    ;;
  *)
    echo "Canceled"
   ;; 
esac

