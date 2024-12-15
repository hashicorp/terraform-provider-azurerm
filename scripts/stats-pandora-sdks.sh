#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


SERVICES=$(ls internal/services/)

TOTAL=0
TOTAL_RESOURCES=$(find internal/services | egrep "_resource.go$" | wc -l | tr -d '[:space:]')
DONE=0
DONE_RESOURCES=0
PARTIAL=0

echo "## Service Packages"
echo
for s in $SERVICES; do
  (( TOTAL++ ))

  RESOURCES=$(find internal/services/$s | egrep "_resource.go" | wc -l | tr -d '[:space:]')
  RESOURCES_DONE=$(grep -rnw "internal/services/$s" -e 'go-azure-sdk' | egrep "_resource.go" | wc -l | tr -d '[:space:]')
  DONE_RESOURCES=$((DONE_RESOURCES + RESOURCES_DONE))

  if grep -rnw "internal/services/$s" -e 'azure-sdk-for-go' > /dev/null; then
      echo "- [ ] \`$s\` _($RESOURCES_DONE/$RESOURCES)_";

      if [ "$RESOURCES_DONE" == "0" ]; then
          (( PARTIAL++ ))
      fi
  else
    (( DONE++ ))

    echo "- [X] \`$s\` _($RESOURCES)_";
  fi

done

echo
echo "services: $DONE (+$PARTIAL partial) of $TOTAL "
echo "resources: $DONE_RESOURCES of $TOTAL_RESOURCES (this is an rough estimation)"

