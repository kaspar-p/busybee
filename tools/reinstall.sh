#!/bin/bash

# Remove and re-add tmp/ directory
rm -rf ./tmp
mkdir ./tmp

# Install pre-commit hooks
pre-commit install -f