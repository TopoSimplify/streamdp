#!/usr/bin/env bash
rm ./benchmark

pwd="../"
go build -o benchmark
./benchmark