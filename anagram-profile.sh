#!/bin/bash

function printLine() {
    printf "\n================= $1 ====================="
}

output="profile_all.txt"

if [ "$1" != "" ]; then
    output=$1
fi

mkdir -p profile
{
    printLine "BENCHMEM"
    go test ./... -bench=. -run=^$ -benchmem 

    printLine "MEM OUT"
    go test ./... -bench=. -run=^$ -memprofile=profile/mem.out
    go tool pprof -text -nodecount=10 src.test profile/mem.out

    printLine "CPU OUT"
    go test ./... -bench=. -run=^$ -cpuprofile=profile/cpu.out
    go tool pprof -text -nodecount=10 src.test profile/cpu.out

    printLine "BLOCK OUT"
    go test ./... -bench=. -run=^$ -blockprofile=profile/block.out
    go tool pprof -text -nodecount=10 src.test profile/block.out
}> profile/$output