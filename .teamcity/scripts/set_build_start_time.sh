#!/bin/bash

echo "##teamcity[setParameter name='env.BUILD_START_TIME' value='$(date +%s)']"