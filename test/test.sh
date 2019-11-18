#!/usr/bin/env bash

# Script to run go test on multiple packages with a single
# coverage report for codecov.io.

# invoked from the github repo as ./test/test.sh
set -e
echo "" > coverage.txt
export GO111MODULE=on # for forks

echo "testing main"
#go build
cd cmd/fit
go test -v -coverprofile=profile.out -covermode=atomic
if [ -f profile.out ]; then
    cat profile.out >> ../../coverage.txt
    rm profile.out
fi
cd ../..

#for d in scm; do # bugapp bugs scm
for d in $(find ./* -maxdepth 0 -type d); do
    if ls $d/*_test.go &> /dev/null; then
        if [[ $d = *issues* ]] ; then
            continue
        fi
        echo "testing in $d"
        cd $d
        go test -v -coverprofile=profile.out -covermode=atomic
        if [ -f profile.out ]; then
            cat profile.out >> ../coverage.txt
            rm profile.out
        fi
        cd ..
    fi
done
