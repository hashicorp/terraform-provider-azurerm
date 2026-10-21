// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// applySrc writes src to a temporary *_resource.go file, lints it, applies the
// auto-fixable findings and returns the rewritten file content plus the renames
// performed.
func applySrc(t *testing.T, src string) (string, []AppliedRename) {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "x_resource.go")
	if err := os.WriteFile(p, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	findings, err := Run([]string{p}, Options{})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	applied, err := Apply(findings)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	out, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(out), applied
}

// TestApplyRenamesEveryOccurrence confirms a rename fix replaces every quoted
// occurrence of the property name — not just the schema key — while leaving
// unrelated text (bare identifiers, superstrings) alone.
func TestApplyRenamesEveryOccurrence(t *testing.T) {
	src := `package x

func r() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"vm_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"vm_count", "location"},
			},
		},
	}
}

type model struct {
	VMCount int ` + "`tfschema:\"vm_count\"`" + `
}

func read(d *pluginsdk.ResourceData) {
	_ = d.Get("vm_count")
	_ = d.Get("vm_count_total")
}
`
	out, applied := applySrc(t, src)

	if len(applied) != 1 || applied[0].From != "vm_count" || applied[0].To != "virtual_machine_count" {
		t.Fatalf("expected one vm_count->virtual_machine_count rename, got %+v", applied)
	}
	for _, want := range []string{
		`"virtual_machine_count": {`,
		`[]string{"virtual_machine_count", "location"}`,
		"`tfschema:\"virtual_machine_count\"`",
		`d.Get("virtual_machine_count")`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected rewritten file to contain %q\n---\n%s", want, out)
		}
	}
	// A superstring token and the Go identifier must be untouched.
	if !strings.Contains(out, `"vm_count_total"`) {
		t.Errorf("superstring \"vm_count_total\" should be untouched\n---\n%s", out)
	}
	if !strings.Contains(out, "VMCount int") {
		t.Errorf("Go identifier VMCount should be untouched\n---\n%s", out)
	}
	if strings.Contains(out, `"vm_count"`) {
		t.Errorf("no bare \"vm_count\" token should remain\n---\n%s", out)
	}
}

// TestApplyRenamesDottedPathReferences confirms a rename also updates dotted
// cross-field references where the property is a path segment (e.g.
// AtLeastOneOf: []string{"network_profile.0.security_enabled"}), not just the
// standalone "network_profile" schema key.
func TestApplyRenamesDottedPathReferences(t *testing.T) {
	src := `package x

func r() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"network_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"security_enabled": {
							Type:         pluginsdk.TypeBool,
							Optional:     true,
							AtLeastOneOf: []string{"network_profile.0.observability_enabled", "network_profile.0.security_enabled"},
						},
					},
				},
			},
		},
	}
}

func read(d *pluginsdk.ResourceData) {
	_ = d.Get("network_profile.0.security_enabled")
}
`
	out, applied := applySrc(t, src)

	if len(applied) != 1 || applied[0].From != "network_profile" || applied[0].To != "network" {
		t.Fatalf("expected one network_profile->network rename, got %+v", applied)
	}
	for _, want := range []string{
		`"network": {`,
		`"network.0.observability_enabled"`,
		`"network.0.security_enabled"`,
		`d.Get("network.0.security_enabled")`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected rewritten file to contain %q\n---\n%s", want, out)
		}
	}
	if strings.Contains(out, "network_profile") {
		t.Errorf("no network_profile reference should remain\n---\n%s", out)
	}
}

// TestApplyLeavesNonRenamableFindings confirms Apply only touches rename fixes
// and reports the rest as unfixed.
func TestApplyLeavesNonRenamableFindings(t *testing.T) {
	src := wrap(`"scalar": {Type: pluginsdk.TypeString, Optional: true, MaxItems: 1},`)
	dir := t.TempDir()
	p := filepath.Join(dir, "x_resource.go")
	if err := os.WriteFile(p, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	findings, err := Run([]string{p}, Options{})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	applied, err := Apply(findings)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(applied) != 0 {
		t.Fatalf("expected no renames applied, got %+v", applied)
	}

	after, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != src {
		t.Errorf("file should be unchanged when there is nothing to auto-fix")
	}

	unfixed := Unfixed(findings)
	if len(unfixed) != len(findings) {
		t.Errorf("expected all %d findings to remain unfixed, got %d", len(findings), len(unfixed))
	}
	found := false
	for _, f := range unfixed {
		if f.RuleID == "SL003" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SL003 (limits-on-non-collection) among unfixed findings")
	}
}
