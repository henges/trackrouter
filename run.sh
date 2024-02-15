#!/usr/bin/env bash

set -ex

readarray -t secrets <".secrets.env"
readarray -t settings <".env"
env "${secrets[@]}" "${settings[@]}" "go" "run" "github.com/henges/trackrouter/cmd/$1"
