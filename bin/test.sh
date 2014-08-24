#!/bin/bash 
pushd .;
cd $GOPATH/src/github.com/jadekler/git-go-logictree/app/home;
go test;
popd;
