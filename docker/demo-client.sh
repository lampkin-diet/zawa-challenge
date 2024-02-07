#!/bin/sh

set -e

# Steps to demo
# 1. Generate let's say 10 files with random content, by calling ""./client generate 10"
# 2. Upload that files to server, by calling "./client upload"
# 3. Check, that files were removed from storage.
# 4. Check that only file `root_hash` is present in storage.
# 5. Download file from server, by calling "./client download file1"
# 6. Check that file was added to storage.
# 7. Download file from server with corruption imitation, by calling "./client download file7 --corrupt"
# 8. Check that file was not added to storage.


CLIENT=/app/client
STORAGE=$STORAGE_PATH
RED='\033[0;31m'
NC='\033[0m' # No Color

function echo_color {
  echo -e "${RED}$1${NC}"
}

function show_files_in_storage {
  for f in $(ls -p ${STORAGE}); do
    echo_color "File in storage: $f"
  done
}

echo_color "0. Create a directory for storage"
mkdir -p ${STORAGE}


echo_color "1. Generate let's say 10 files with random content, by calling ./client generate 10"

${CLIENT} generate 10
FILES_COUNT=$(ls -p ${STORAGE} | wc -l)
if [ $FILES_COUNT -ne 10 ]; then
  echo_color "Files were not generated"
  exit 1
fi

show_files_in_storage

echo_color '2. Upload files to server, by calling "./client upload"'

${CLIENT} upload

echo_color '3. Check, that files were removed from storage.'

FILES_COUNT=$(ls -p ${STORAGE} | wc -l)
if [ $FILES_COUNT -ne 1 ]; then
  echo_color "Files were not removed from storage"
  exit 1
fi

show_files_in_storage

echo_color '# 4. Check that only file `root_hash` is present in storage.'

if ! test -f "${STORAGE}/root_hash"; then
  echo_color "root_hash file is not present in storage"
  exit 1
fi

echo_color '# 5. Download file from server, by calling "./client download file1"'

${CLIENT} download file1

echo_color '# 6. Check that file was added to storage.'

if ! test -f "${STORAGE}/file1"; then
  echo_color "file1 was not added to storage"
  exit 1
fi

show_files_in_storage

echo_color '# 7. Download file from server with corruption imitation, by calling "./client download file7 --corrupt"'

${CLIENT} download file7 --corrupt

echo_color '# 8. Check that file was not added to storage.'

if test -f "${STORAGE}/file7"; then
  echo_color "file7 was added to storage even proof is not valid"
  exit 1
fi

show_files_in_storage

echo_color "Demo is finished"
