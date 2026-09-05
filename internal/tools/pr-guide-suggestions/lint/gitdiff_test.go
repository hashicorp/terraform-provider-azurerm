// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import (
	"path/filepath"
	"testing"
)

func TestParseHunkNewStart(t *testing.T) {
	cases := map[string]int{
		"@@ -1,0 +23,4 @@":            23,
		"@@ -5 +7 @@":                 7,
		"@@ -1,2 +10,0 @@ func x() {": 10,
	}
	for header, want := range cases {
		got, ok := parseHunkNewStart(header)
		if !ok || got != want {
			t.Errorf("parseHunkNewStart(%q) = %d,%v; want %d", header, got, ok, want)
		}
	}
}

func TestParseDiffPath(t *testing.T) {
	root := "/repo"
	if got := parseDiffPath(root, "b/internal/x.go"); got != filepath.Join(root, "internal/x.go") {
		t.Errorf("unexpected path %q", got)
	}
	if got := parseDiffPath(root, "/dev/null"); got != "" {
		t.Errorf("expected empty path for /dev/null, got %q", got)
	}
}
