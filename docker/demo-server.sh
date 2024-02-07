#!/bin/sh

set -e

# Steps to demo
# 1. Create a directory for storage
# 2. Start the server


STORAGE=$STORAGE_PATH

echo "1. Create a directory for storage"
mkdir -p ${STORAGE}

echo "2. Start the server"
/app/server
