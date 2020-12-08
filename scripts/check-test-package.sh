#!/usr/bin/env sh

files=$(find . | egrep "/azurerm/internal/services/[a-z]+/[a-z_]+(resource|data_source)[a-z_]+\.go$" | egrep "test.go")
error=false

echo "==> Checking that acceptance test packages are used..."

for f in $files; do
  line=$(head -n 1 $f)
  if [ "$line" = "${line%%_test}" ]; then
    echo $f
    error=true
  fi
done

if $error; then
  echo ""
  echo "------------------------------------------------"
  echo ""
  echo "The acceptance test files listed above are using the same package as the resource or data source code."
  echo "They must use a test package to prevent a circular dependency. To fix this change the first line:"
  echo ""
  echo "> package service"
  echo ""
  echo "to"
  echo ""
  echo "> package service_test"
  echo ""
  exit 1
fi

exit 0
