#!/bin/bash
cd static/proto; protoc --go_out=../../proto */*.proto
