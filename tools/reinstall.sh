#!/bin/bash

chmod 755 ./tools/install-hooks
./tools/install-hooks

rm -rf ./tmp
mkdir ./tmp