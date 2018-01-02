#!/usr/bin/env bash
rm ./client

pwd="../"

cd ${pwd}/client
    go build -o ${pwd}/bin/client
cd ${pwd}/bin

./client