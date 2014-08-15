#!/bin/bash 
pushd .;
cd $GOPATH/src/git-go-logictree/app/home;
go test;
popd;
