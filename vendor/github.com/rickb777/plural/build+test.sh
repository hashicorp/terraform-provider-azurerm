#!/bin/bash -ex
cd "$(dirname $0)"
go install tool
mage
