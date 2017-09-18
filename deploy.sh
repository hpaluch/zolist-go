#!/bin/bash

cd `dirname $0`
source _scripts/func.sh
gen_app_yaml
set -e
# check number of version (max 15)
max_versions=15
deployed_versions=$(gcloud app versions list | egrep -v '^SERVICE' | wc -l)

if [ $deployed_versions -ge $max_versions ]
then
	last_version=$(gcloud app versions list |
                       egrep -v '^SERVICE' | head -1 | awk '{print $2}')
	[ -n "$last_version" ]
	echo "Too many deployed versions $deployed_versions - deleteing last one: $last_version"
	set -x
	gcloud app versions delete $last_version
	set +x
fi
set -x
gcloud app deploy
# please ignore "skipped"/"copied" messages before deployment
# they have nothing common with Upload filter...
#gcloud app deploy --verbosity=info
gcloud app browse

