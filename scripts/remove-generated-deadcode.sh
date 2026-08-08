#!/usr/bin/env bash
# Removes fully-dead generated resource-ID families after 'make deadcode-fix'.
#
# For each modified generated parse file (marker: "generated via 'go:generate'") whose
# NewXxxID constructor AND XxxID parser were both removed by the fix, the whole family is
# deleted - parse/validate files, their tests, and the '//go:generate ... -name=Xxx'
# directive in the service's resourceids.go - so 'make generate' cannot resurrect it.
# IDs with any surviving (live) function are skipped.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

git diff --name-only --diff-filter=M -- 'internal/services/*/parse/*.go' | grep -v '_test\.go$' | while read -r pf; do
  grep -q "generated via 'go:generate'" "$pf" || continue
  base=$(basename "$pf" .go)
  svc=$(dirname "$(dirname "$pf")")

  # the ID name comes from the pre-fix constructor, e.g. func NewApiDiagnosticID(
  name=$(git show "HEAD:$pf" | sed -n 's/^func New\([A-Za-z0-9]*\)ID(.*/\1/p' | head -1)
  if [ -z "$name" ]; then
    echo "skip (no constructor in HEAD): $pf"
    continue
  fi

  # fully dead only if the fix removed every generated function - the constructor, the
  # parser, and any XxxIDInsensitively variant (no paren anchor so variants match)
  if grep -qE "func (New)?${name}ID" "$pf"; then
    echo "skip (still has live funcs): $pf"
    continue
  fi

  for f in "$svc/parse/$base.go" "$svc/parse/${base}_test.go" "$svc/validate/${base}_id.go" "$svc/validate/${base}_id_test.go"; do
    if [ -f "$f" ]; then
      echo "  rm $f"
      rm "$f"
    fi
  done

  rid="$svc/resourceids.go"
  if [ -f "$rid" ] && grep -qE -- "-name=${name}( |$)" "$rid"; then
    grep -vE -- "-name=${name}( |$)" "$rid" > "$rid.tmp" && mv "$rid.tmp" "$rid"
    echo "  directive removed: $rid (-name=$name)"
    if ! grep -q 'go:generate' "$rid"; then
      echo "  rm $rid (no directives left)"
      rm "$rid"
    fi
  else
    echo "  WARNING: no -name=${name} directive found in $rid"
  fi
done

# remove now-empty parse/validate directories
find internal/services -type d \( -name parse -o -name validate \) -empty -print -delete
