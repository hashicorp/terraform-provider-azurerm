// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/rules"
)

// lintSrc writes src to a temporary *_resource.go file and lints it. Because the
// linter is purely syntactic, the source need not type-check or resolve imports.
func lintSrc(t *testing.T, src string) []Finding {
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
	return findings
}

// rulesAtPath returns the set of rule IDs reported for the given property path.
func rulesAtPath(findings []Finding, path string) map[string]bool {
	out := map[string]bool{}
	for _, f := range findings {
		if f.Path == path {
			out[f.RuleID] = true
		}
	}
	return out
}

func wrap(body string) string {
	return `package x
func r() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
` + body + `
		},
	}
}
`
}

// TestRunResolvesContext exercises the end-to-end pipeline: a finding carries
// the resolved resource type, kind and property path, and an opaque helper-call
// property is skipped by schema-dependent rules.
func TestRunResolvesContext(t *testing.T) {
	f := lintSrc(t, `package x

type FooResource struct{}

func (FooResource) ResourceType() string { return "azurerm_foo" }

func (FooResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name":     {Type: pluginsdk.TypeString, Optional: true, MaxItems: 1},
		"location": commonschema.Location(),
	}
}
`)

	if len(rulesAtPath(f, "location")) != 0 {
		t.Error("helper-call property should be skipped by schema-dependent rules")
	}
	if !rulesAtPath(f, "name")["SL003"] {
		t.Error("expected SL003 on name (MaxItems on a non-collection)")
	}
	for _, x := range f {
		if x.Path == "name" && (x.Resource != "azurerm_foo" || x.Kind != "resource") {
			t.Errorf("expected azurerm_foo/resource context, got %q/%q", x.Resource, x.Kind)
		}
	}
}

// TestRegistryComplete checks every rule is registered with distinct, valid
// metadata.
func TestRegistryComplete(t *testing.T) {
	if len(rules.Registry) != 13 {
		t.Fatalf("expected 13 rules, got %d", len(rules.Registry))
	}
	seen := map[string]bool{}
	for _, r := range rules.Registry {
		if r.ID == "" || r.Name == "" || r.Check == nil {
			t.Errorf("rule %+v has empty metadata", r)
		}
		if seen[r.ID] {
			t.Errorf("duplicate rule ID %q", r.ID)
		}
		seen[r.ID] = true
		if r.Severity != rules.Error && r.Severity != rules.Warning {
			t.Errorf("rule %s has invalid severity %q", r.ID, r.Severity)
		}
	}
}

func TestRuleFilters(t *testing.T) {
	src := wrap(`"s": {Type: pluginsdk.TypeString, Optional: true},`)
	dir := t.TempDir()
	p := filepath.Join(dir, "x_resource.go")
	if err := os.WriteFile(p, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}

	only, err := Run([]string{p}, Options{Only: map[string]bool{"SL001": true}})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range only {
		if f.RuleID != "SL001" {
			t.Errorf("with -rules SL001, unexpected %s", f.RuleID)
		}
	}

	disabled, err := Run([]string{p}, Options{Disable: map[string]bool{"SL001": true}})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range disabled {
		if f.RuleID == "SL001" {
			t.Error("SL001 should have been disabled")
		}
	}
}

// TestDiffFiltersBySpan checks that diff mode keeps a finding when a property's
// body changed, even though the key line it is anchored on did not.
func TestDiffFiltersBySpan(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "x_resource.go")
	// Line 5 is the "sku_name" key; lines 6-7 are its body.
	src := `package x
func r() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
		},
	}
}
`
	if err := os.WriteFile(p, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}

	// A change to the body (line 6) but not the key line (line 5) is kept.
	body := &Changes{files: map[string]map[int]bool{p: {6: true}}}
	got, err := Run([]string{p}, Options{Changes: body})
	if err != nil {
		t.Fatal(err)
	}
	if !rulesAtPath(got, "sku_name")["SL005"] {
		t.Error("expected SL005 to be kept when the property body changed")
	}

	// A change outside every property's span is filtered out.
	outside := &Changes{files: map[string]map[int]bool{p: {1: true}}}
	got, err = Run([]string{p}, Options{Changes: outside})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("expected no findings for a change outside any property, got %d", len(got))
	}
}
