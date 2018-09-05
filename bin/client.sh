#!/usr/bin/env bash
export GIN_MODE=release

rm ./client

pwd="../"

cd ${pwd}/client
    go build -o ${pwd}/bin/client
cd ${pwd}/bin

./client