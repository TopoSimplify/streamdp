#!/usr/bin/env bash
export GIN_MODE=release
pwd="../"

rm ./server
cd ${pwd}
    go build -o bin/server
cd bin/

./server