#!/bin/bash

REPO=github.com/beyond-blockchain/bbclib-go
VERSION=v1.4.2
MODULE_DIR="${GOPATH}/pkg/mod/${REPO}@${VERSION}"

go get -u -d ${REPO}

if [ -d ${MODULE_DIR} ]; then
    echo "Install libbbcsig into ${MODULE_DIR}"
    WORKINGDIR=${MODULE_DIR}
    chmod 744 ${GOPATH}/pkg/mod/github.com/beyond-blockchain ${MODULE_DIR}
elif [ -d ./vendor/${REPO} ]; then
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
