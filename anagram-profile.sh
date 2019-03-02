#!/bin/bash

function printLine() {
    printf "\n================= $1 ====================="
}

mkdir -p profile
{
    printLine "BENCHMEM"
    go test -bench=. -benchmem

    printLine "MEM OUT"
    go test -bench=. -memprofile=profile/mem.out
    go tool pprof -text -nodecount=10 anagram-find.test profile/mem.out

    printLine "CPU OUT"
    go test -bench=. -cpuprofile=profile/cpu.out
    go tool pprof -text -nodecount=10 anagram-find.test profile/cpu.out

    printLine "BLOCK OUT"
    go test -bench=. -blockprofile=profile/block.out
    go tool pprof -text -nodecount=10 anagram-find.test profile/block.out
}> profile/profile_all.txt