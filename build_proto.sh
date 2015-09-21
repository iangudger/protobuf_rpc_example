#!/bin/bash
rm -rf proto
mkdir proto
cd static/proto; find . -type f -name "*.proto" | xargs protoc --go_out=../../proto
