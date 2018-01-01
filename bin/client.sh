#!/usr/bin/env bash
pwd="../"
cd ${pwd}/client
    go build -o ${pwd}/bin/client
cd ${pwd}/bin

./client