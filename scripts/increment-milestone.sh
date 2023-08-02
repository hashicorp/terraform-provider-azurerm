#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


while getopts u:r:t: flag
do
  case "${flag}" in
    r) release=${OPTARG};;
    u) milestone_url=${OPTARG};;
    t) token=${OPTARG};;
  esac
done

echo "Getting current milestone number..."
milestones=$(curl -L \
-H "Accept: application/vnd.github+json" \
-H "Authorization: Bearer $token" \
-H "X-GitHub-Api-Version: 2022-11-28" \
"${milestone_url}?state=open&sort=due_on&direction=desc")

milestone_number=0
milestones_json=$(echo "$milestones" | jq -c -r '.[]')
for milestone in ${milestones_json[@]}; do
  if [[ $(echo $milestone | jq -r .title) == "$release" ]]; then
    milestone_number=$(echo $milestone | jq -r .number)
    break
  fi
done

if [[ $milestone_number != 0 ]]; then

  echo "Closing current milestone..."
  curl -L \
  -X PATCH \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $token" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "${milestone_url}/${milestone_number}" \
  -d '{"state":"closed"}'

  major=0
  weekly=0
  patch=0
  regex="v([0-9]+).([0-9]+).([0-9]+)"
  if [[ $release =~ $regex ]]; then
    major="${BASH_REMATCH[1]}"
    weekly="${BASH_REMATCH[2]}"
    patch="${BASH_REMATCH[3]}"
  fi

  if [[ $patch == 0 ]]; then
    if [[ $major != 0 ]]; then

      # Get next release version
      weekly=$((weekly + 1 ))
      new_milestone="v$major.$weekly.0"
      # Get next release due date
      date=$(date -d "next Thursday" +%Y-%m-%d)
      date+="T12:00:00Z"

      echo "Creating new milestone..."
      curl -L \
      -X POST \
      -H "Accept: application/vnd.github+json" \
      -H "Authorization: Bearer $token" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      "${milestone_url}" \
      -d "{\"title\":\"$new_milestone\", \"state\":\"open\", \"due_on\":\"$date\"}"
    else
      echo "Could not increment milestone"
      exit 1
    fi
  fi

else
    echo "Could not retrieve current milestone number to close"
    exit 1
fi