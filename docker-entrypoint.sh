#!/bin/bash

yell() { echo "$0: $*" >&2; }
die()  { yell "$*"; exit 111; }
try()  { echo "$ $@" 1>&2; "$@" || die "cannot $*"; }

conf="/dockerconfig.env"
yell "Looking for the optional $conf file..."

if [ -f $conf ]; then
	. $conf
    yell "Sourced $conf"
else
    yell "File $conf does not exist."
fi

yell "Running custom scripts in /docker-entrypoint.d/ ..."

# to prevent expansion to literal string `/docker-entrypoint.d/*` when there is nothing matching the glob
shopt -s nullglob

for file in /docker-entrypoint.d/*; do
    if [ -f "$file" ] && [ -x "$file" ]; then
        try "$file"
    else
        yell "Ignoring $file, it is either not a file or is not executable"
    fi
done

try exec "$@"