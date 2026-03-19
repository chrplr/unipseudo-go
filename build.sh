#!/bin/bash

# Exit on error
set -e

echo "Building pseudoword-generator for local platform..."

# Compile the Go code
go build -o pseudoword-generator pseudoword_generator.go

echo "Build successful! Executable created: ./pseudoword-generator"
