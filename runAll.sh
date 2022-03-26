#!/bin/bash
go build

rm -rf output
mkdir -p output

for file in input/*
do
    filename=${file##input/}
    filename_without_extension=${filename%%.in}
    ./hashcode-pizza $file output/$filename_without_extension.out;
done
