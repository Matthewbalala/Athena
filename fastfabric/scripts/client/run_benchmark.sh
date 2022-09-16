#!/usr/bin/env bash
offset=$1
threadCount=$2
target=$4
conflictPercentage=$5
accountCount=$6
totalThreads=$3
for i in $(seq $offset $(($threadCount-1)))
do
    node ./invoke2.js $i $(($accountCount/$totalThreads)) 1 $target $conflictPercentage&
done
