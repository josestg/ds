#!/usr/bin/env sh

set -e

mkdir -p build
cd build
cmake ..
cmake --build .

if [ -f ./test_runner ]; then
  ctest --output-on-failure
fi
