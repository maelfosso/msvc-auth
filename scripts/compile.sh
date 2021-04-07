#!/bin/bash

set -e

echo compiling

go build -o ./bin/start.bin *.go