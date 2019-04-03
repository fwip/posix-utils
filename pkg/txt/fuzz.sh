#!/bin/bash

here="$(dirname "$0")"
cd "$here/go-fuzz-piece-table/" || exit

go-fuzz-build github.com/fwip/posix-utils/pkg/txt
go-fuzz -bin=txt-fuzz.zip -workdir=.

