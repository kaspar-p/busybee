#!/bin/bash

rm -rf bin
mkdir bin
cd src
go build -o ../bin/busybee ./
cd ..