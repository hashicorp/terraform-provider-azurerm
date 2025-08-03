#!/usr/bin/env bash

echo "==> Validating documentation..."

if ! go run ./internal/tools/document-fmt/main.go validate; then
  echo ""
  echo "------------------------------------------------"
  echo "Encountered errors validating the documentation."
  echo "To fix these errors, run \`make document-fix\`."
  echo "------------------------------------------------"
  echo ""

  exit 1
fi

exit 0
