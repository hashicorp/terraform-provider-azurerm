#!/bin/bash

echo -n "Set \$TCPASS to: "
read -s INPUT
echo

export TCPASS="$INPUT"
