#!/usr/bin/env bash
pwd="../"

rm ./server
cd ${pwd}
    go build -o bin/server
cd bin/

./server