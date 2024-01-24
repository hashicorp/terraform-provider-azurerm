#!/bin/bash -e
cd "$(dirname $0)"
PATH=$HOME/go/bin:$PATH

if ! type -p goveralls; then
  echo go install github.com/mattn/goveralls
  go install github.com/mattn/goveralls
fi

echo plural...
go test -v -covermode=count -coverprofile=date.out .
go tool cover -func=date.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=date.out -service=travis-ci -repotoken $COVERALLS_TOKEN
