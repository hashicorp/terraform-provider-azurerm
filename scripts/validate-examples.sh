#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


echo "==> Checking examples validate with 'terraform validate'..."

exampleDirs=$(find ./examples -mindepth 2 -maxdepth 2 -type d '!' -exec test -e "{}/*.tf" ';' -print | sort | egrep -v "tfc-checks")
examplesWithErrors=()
hasError=false

# Setup a local Terraform config file for setting up the dev_overrides for this provider.
terraformrc=$(mktemp)
cat << EOF > $terraformrc
provider_installation {
  dev_overrides {
    "hashicorp/azurerm" = "$(go env GOPATH)/bin"
  }
  direct {}
}
EOF

# first check each example validates via `terraform validate`..
for d in $exampleDirs; do
  echo "Validating $d.."
  exampleHasErrors=false
  # Though we are using the local built azurerm provider to validate example,
  # we still need to call `terraform init` here as examples might contain other providers.
  TF_CLI_CONFIG_FILE=$terraformrc terraform -chdir=$d init > /dev/null || exampleHasErrors=true
  if ! ${exampleHasErrors}; then
    # Always use the local built azurerm provider to validate examples, to avoid examples using
    # unreleased features leading to error during CI.
    TF_CLI_CONFIG_FILE=$terraformrc terraform -chdir=$d validate > /dev/null || exampleHasErrors=true
  fi
  if ${exampleHasErrors}; then
    examplesWithErrors[${#examplesWithErrors[@]}]=$d
    hasError=true
  fi
done

if ${hasError}; then
  echo "------------------------------------------------"
  echo ""
  echo "The directories listed below failed to validate using 'terraform validate'"
  echo "Please fix the validation errors for these, these can be found by running"
  echo "'terraform init' and then 'terraform validate':"
  for exampleDir in "${examplesWithErrors[@]}"
  do
       echo "- $exampleDir"
  done
  exit 1
fi

# Finally check everything is formatted
for d in $exampleDirs; do
  echo "Checking formatting for $d.."
  exampleHasErrors=false
  terraform fmt -check $d > /dev/null || exampleHasErrors=true
  if ${exampleHasErrors}; then
    examplesWithErrors[${#examplesWithErrors[@]}]=$d
    hasError=true
  fi
done

if ${hasError}; then
  echo "------------------------------------------------"
  echo ""
  echo "The directories listed below aren't formatted using 'terraform fmt'."
  echo "Please fix the validation errors for these, these can be found by running"
  echo "'terraform fmt':"
  for exampleDir in "${examplesWithErrors[@]}"
  do
       echo "- $exampleDir"
  done
  exit 1
fi

echo "All Examples Validate and Format"
exit 0
