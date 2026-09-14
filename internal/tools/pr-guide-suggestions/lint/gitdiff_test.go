// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import (
	"path/filepath"
	"testing"
)

func TestParseDiffAddedLines(t *testing.T) {
	root := "/repo"
	diff := `diff --git a/internal/services/foo/foo_resource.go b/internal/services/foo/foo_resource.go
index bbb0728..ac91822 100644
--- a/internal/services/foo/foo_resource.go
+++ b/internal/services/foo/foo_resource.go
@@ -9,3 +9,6 @@ func r() *pluginsdk.Resource {
-				Removed: true,
+				Type:     pluginsdk.TypeString,
+				Required: true,
+			},
+			"new_config": {
+				Type:     pluginsdk.TypeString,
+				Optional: true,
`

	changes, err := parseDiff(root, diff)
	if err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(root, "internal/services/foo/foo_resource.go")
	// New-side lines 9..14 are added; the removal does not advance the counter.
	for _, line := range []int{9, 10, 11, 12, 13, 14} {
		if !changes.Added(file, line) {
			t.Errorf("expected line %d to be added", line)
		}
	}
	if changes.Added(file, 8) {
		t.Error("line 8 (unchanged context) should not be marked added")
	}
	if changes.Added(file, 15) {
		t.Error("line 15 should not be marked added")
	}
}

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
