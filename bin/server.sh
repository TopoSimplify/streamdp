#!/usr/bin/env bash
pwd="../"
cd ${pwd}
    go build -o bin/server
cd bin/

./server