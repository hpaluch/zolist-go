#!/bin/bash

# shell helper functions

check_zomato_key () {

	[ "x$ZOMATO_API_KEY" != "x" ] || {
		echo -e "ERROR:\tUndefined variable ZOMATO_API_KEY" >&2
		echo -e "\tGo to https://developers.zomato.com/api to get API key" >&2
		echo -e "\tand set shell variable like: export ZOMATO_API_KEY=my_key" >&2
		exit 1
	}


}

