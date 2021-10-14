#!/usr/bin/env bash

echo "==> Checking examples validate with 'terraform validate'..."

exampleDirs=$(find ./examples -mindepth 2 -maxdepth 2 -type d '!' -exec test -e "{}/*.tf" ';' -print | sort)
examplesWithErrors=()
hasError=false

for d in $exampleDirs; do
  echo "Validating $d.."
  exampleHasErrors=false
  terraform -chdir=$d init > /dev/null || exampleHasErrors=true
  if ! ${exampleHasErrors}; then
    terraform -chdir=$d validate > /dev/null || exampleHasErrors=true
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

exit 0
