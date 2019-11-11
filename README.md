bbclib-go
====
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://travis-ci.org/beyond-blockchain/bbclib-go.svg?branch=develop)](https://travis-ci.org/beyond-blockchain/bbclib-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/beyond-blockchain/bbclib-go)](https://goreportcard.com/report/github.com/beyond-blockchain/bbclib-go)
[![Coverage Status](https://coveralls.io/repos/github/beyond-blockchain/bbclib-go/badge.svg?branch=develop)](https://coveralls.io/github/beyond-blockchain/bbclib-go?branch=develop)
[![Maintainability](https://api.codeclimate.com/v1/badges/0c523f5a3d71b77aad46/maintainability)](https://codeclimate.com/github/beyond-blockchain/bbclib-go/maintainability)

Golang implementation of bbc1.core.bbclib and bbc1.core.libs modules in https://github.com/beyond-blockchain/bbc1.
This reposigory is originally from https://github.com/quvox/bbclib-go


### Features
* Support most of the features of py-bbclib in https://github.com/beyond-blockchain/py-bbclib
    * BBc-1 version 1.6
    * transaction header version 1 and 2.
* Go v1.12 or later (need go mod)

## Usage

import "github.com/beyond-blockchain/bbclib-go"

An example source code is in example/.


## Install (step by step)

```bash
go get -u github.com/beyond-blockchain/bbclib-go
```

NOTE: [example/](./example) directory includes a sample code for this module.
