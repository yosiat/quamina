#!/bin/sh

go test -benchmem -run="^$" -bench "^Benchmark"
