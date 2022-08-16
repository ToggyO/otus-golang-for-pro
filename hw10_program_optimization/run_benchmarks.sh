#!/usr/bin/env bash

isOld="true"
go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 10s -count 5 -args -old=true | tee old

isOld="false"
go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 10s -count 5 | tee new

benchstat old new > benchstat

echo "Finish"