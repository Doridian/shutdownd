#!/bin/sh
RET="$(echo "shutdownd|trigger|$2" | nc -4u "$1" 10001 -w 1)"
if echo -n "$RET" | grep -q '^shutdownd|ok$'
then
	echo "OK"
	exit 0
else
	echo "Error: $RET"
	exit 1
fi
