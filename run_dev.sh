#!/bin/bash

cd `dirname $0`
set -ex
dev_appserver.py app.yaml

