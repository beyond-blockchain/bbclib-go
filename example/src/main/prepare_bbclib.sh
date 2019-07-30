#!/bin/bash

REPO=github.com/beyond-blockchain/bbclib-go

go get -u -d ${REPO}

if [ -d ./vendor/${REPO} ]; then
    WORKINGDIR=./vendor/${REPO}
elif [ -d ${GOPATH}/src/${REPO} ]; then
    WORKINGDIR=${GOPATH}/src/${REPO}
else
    echo "No bbclib-go repository found..."
    exit 1
fi

cd ${WORKINGDIR}

if [ $# -eq 1 ] && [ $1 = "aws" ]; then
  bash prepare.sh aws
else
  bash prepare.sh
fi

# go install ${REPO}