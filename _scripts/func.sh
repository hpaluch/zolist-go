#!/bin/bash

# shell helper functions

gen_app_yaml () {
	# default Restaurant IDS:   Ldiak    Na Pude  Flamingo
	export REST_IDS=${REST_IDS:-18355040,16513797,16512711}
	{
		echo "# DO NOT EDIT - Generated at `date`"
		cat app.yaml.template
		echo "env_variables:"
		for i in ZOMATO_API_KEY REST_IDS
		do
			eval val="\$$i"
			[ -n "$val" ] || {
				echo "Mandatory variable '$val' undefined" >&2
				exit 1
			}
			echo "    $i: '$val'"
		done
	} > app.yaml

}

# point GOPATH to over src/github.com/hpaluch/zolist-go
export GOPATH=../../../../

