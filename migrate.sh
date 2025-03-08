#!/usr/bin/env sh

export CUR="github.com/carloscasalar/gin-starter"
export NEW="cca-fever-challenge"
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;

