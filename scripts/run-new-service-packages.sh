#!/usr/bin/env bash

# This is a weird script, in short Github will only run (new/changed) Github Actions
# in a PR when that PR originates from the main repository itself (e.g. not forks)
#
# As such whilst we generate the new Github Action steps as a part of PR's, they can't
# always be run - so this script instead takes a list of service packages from the `main`
# branch (stored in `./scripts/.service-packages`) - lists the service packages available
# in this branch and then runs the acceptance tests for each package

# ./scripts/.service-packages

function getAvailableServicePackages {
  find ./internal/services -maxdepth 1 -type d | xargs basename | sort
}

function getKnownServicePackages {
  cat ./scripts/.service-packages
}

function findServicePackagesToRun {
  servicesToRun=()

  available=$(getAvailableServicePackages)
  known=$(getKnownServicePackages)

  # find a list of services which won't be run automatically
  for service in $available
  do
    # shellcheck disable=SC2076
    if [[ ! "${known[*]}" =~ "${service}" ]]; then
      servicesToRun+=($service)
    fi
  done

  echo "${servicesToRun[@]}"
}

function runTestsForServicePackage {
  local service=$1

  echo "Running Tests and Linters for Service Package $service.."
  ./scripts/service-package.sh "$service"
  if (($? > 1)); then
      echo "Build Failed."
      exit 1
  fi
}

function main {
  packagesToRun=$(findServicePackagesToRun)
  for package in $packagesToRun
  do
    runTestsForServicePackage "$package"
    if (($? > 1)); then
        echo "Build Failed."
        exit 1
    fi
  done
}

main
