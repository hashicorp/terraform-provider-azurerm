// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package typed_upgrade

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

// derivePaths returns the generated file path and the renamed (legacy) file
// path for a resource. A resource at "foo_resource.go" becomes:
//   - generated:  "foo_resource.go"  (same path — the original is renamed first)
//   - renamed:    "foo_resource_gen.go" → kept temporarily so callers can diff
//
// Actually: the convention from the migration skill is to rename the original
// to *_legacy.go, write the new file at the original path.
func derivePaths(origPath string) (generatedPath, renamedPath string) {
	dir := filepath.Dir(origPath)
	base := filepath.Base(origPath)
	stem := strings.TrimSuffix(base, ".go")
	return origPath, filepath.Join(dir, stem+"_legacy.go")
}

// deriveRegistrationPath finds the registration.go in the same directory.
func deriveRegistrationPath(resourcePath string) string {
	dir := filepath.Dir(resourcePath)
	candidate := filepath.Join(dir, "registration.go")
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

// writeFiles renames the original resource to its legacy path and writes the
// generated source to the original path (so git tracks it as a modification
// rather than a delete + add).
func writeFiles(origPath, generatedPath, renamedPath string, generatedSrc []byte, overwrite bool) error {
	// Check whether the generated path is the same as origPath (it always is
	// by our convention). If so, we must rename first.
	if generatedPath == origPath {
		if _, err := os.Stat(renamedPath); err == nil && !overwrite {
			return fmt.Errorf("legacy file %q already exists; pass -overwrite to replace it", renamedPath)
		}
		if err := os.Rename(origPath, renamedPath); err != nil {
			return fmt.Errorf("renaming %q to %q: %w", origPath, renamedPath, err)
		}
	}

	return writeFile(generatedPath, generatedSrc, overwrite)
}

// writeFile writes src to path. It formats the source with gofmt first.
func writeFile(path string, src []byte, overwrite bool) error {
	if _, err := os.Stat(path); err == nil && !overwrite {
		return fmt.Errorf("file %q already exists; pass -overwrite to replace it", path)
	}

	// Best-effort gofmt — write even if formatting fails.
	if formatted, err := format.Source(src); err == nil {
		src = formatted
	}

	if err := os.WriteFile(path, src, 0o644); err != nil {
		return fmt.Errorf("writing %q: %w", path, err)
	}
	return nil
}
