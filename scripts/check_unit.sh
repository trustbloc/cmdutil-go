#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

pwd=`pwd`
touch "$pwd"/coverage.out

amend_coverage_file () {
if [ -f profile.out ]; then
    cat profile.out | grep -v ".gen.go" >> "$pwd"/coverage.out
    rm profile.out
fi
}

# Running cmdutil-go unit tests
echo "cmdutil-go unit tests..."
PKGS=`go list github.com/trustbloc/cmdutil-go/... 2> /dev/null | \
                                                  grep -v /mocks`
go test $PKGS -count=1 -race -coverprofile=profile.out -covermode=atomic -timeout=10m
amend_coverage_file
echo "... done unit tests"

cd "$pwd"
