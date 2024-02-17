#!/usr/bin/env bash
readarray -t secrets < ".secrets.env"
readarray -t settings < ".env"
env "${secrets[@]}" "${settings[@]}" "go" "run" "github.com/henges/trackrouter/cmd/$1"
