#!/bin/bash

git clone -b master https://github.com/beyond-blockchain/libbbcsig.git libs
cd libs

if [ -z $1 ]; then
    echo "# Prepare for mac/linux"
    bash prepare.sh
elif [ $1 == "aws" ]; then
    echo "# Prepare for AWS Lambda"
    bash prepare.sh aws
fi

if [ -f "lib/libbbcsig.dylib" ]; then
    cp lib/libbbcsig.dylib ../
elif [ -f "lib/libbbcsig.so" ]; then
    cp lib/libbbcsig.so ../
fi
cp lib/libbbcsig.a ../
cp lib/libbbcsig.h ../

