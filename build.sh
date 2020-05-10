#!bin /bash

function buildAndRun () {
    printf "Building Main"

    go build ./main.go

    printf "Running Main"

    chmod 700 ./main

    ./main
}

buildAndRun
