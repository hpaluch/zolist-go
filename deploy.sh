#!/bin/bash

cd `dirname $0`
source _scripts/func.sh
gen_app_yaml
set -ex
gcloud app deploy
# please ignore "skipped"/"copied" messages before deployment
# they have nothing common with Upload filter...
#gcloud app deploy --verbosity=info
gcloud app browse

