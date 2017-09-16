#!/bin/bash

cd `dirname $0`
source _scripts/func.sh
check_zomato_key
set -ex
dev_appserver.py app.yaml

