Sample application using bbclib-go
======

## install

### in the case of using "go get"

```
go get -u -d github.com/beyond-blockchain/bbclib-go
pushd ${GOPATH}/src/github.com/beyond-blockchain/bbclib-go
bash prepare.sh
popd
go install github.com/beyond-blockchain/bbclib-go
```

The important point is adding "-d" in the first go get. This means downloads repository only. Before installing, libbbcsig.a needs to be built.

### in the case of vendoring through "dep"

```
dep init
pushd vendor/github.com/beyond-blockchain/bbclib-go
bash prepare.sh
popd
go install github.com/beyond-blockchain/bbclib-go
```

The difference between "go get" and "dep" is just a directory including libraries.

### Utility script

In the src/main directory, there is "prepare_libbbcsig.h" for building libbbcsig.a in the case of both "go get" and "dep".


## Build app

In src/main directory,

```
go run main.go
```

or

```
go build main.go
./main
```

then you will see the string output of a sample transaction.

