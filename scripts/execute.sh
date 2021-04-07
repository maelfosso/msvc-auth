#!/bin/bash

# Ensure no zombied/lingering pids are arround

pids=$(pgrep start.bin)
[[ -z "${pids}" ]] || kill -9 "${pids}"

# execute the binary
./bin/start.bin