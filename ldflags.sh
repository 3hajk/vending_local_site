#!/usr/bin/env bash
set -e

VERSION_PKG=main
#get highest tag number, and add 1.0.0 if doesn't exist
BASE_VERSION=`git describe --abbrev=0 --tags 2>/dev/null`

if [[ BASE_VERSION == '' ]]
then
  BASE_VERSION='1.0.0'
fi

COMMIT_HASH=$(git show -s --format=%h)
COMMIT_STAMP=$(git show -s --format=%ct)
BUILD_USER=$(id -u -n)
BUILD_HOST=$(hostname)
BUILD_STAMP=$(date +%s)

vflag() {
    VFLAGS="$VFLAGS -X $VERSION_PKG.$1=$2"
}

vflag Number $BASE_VERSION
vflag CommitHash $COMMIT_HASH
vflag CommitStamp $COMMIT_STAMP
vflag BuildUser $BUILD_USER
vflag BuildHost $BUILD_HOST
vflag BuildStamp $BUILD_STAMP

echo "$VFLAGS"