#!/bin/bash

BUILD_TIME=$(date +%s)
echo "##teamcity[setParameter name='env.BUILD_START_TIME' value='$BUILD_TIME']"

GO_VERSION=$(goenv local)
goenv install -s "$GO_VERSION"
goenv rehash