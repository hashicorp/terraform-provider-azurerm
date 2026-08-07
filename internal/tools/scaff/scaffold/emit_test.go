// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEmitRoundTrip(t *testing.T) {
	res := fixtureIR()
	out := Emit(res, EmitOptions{
		ResourceName: "widget",
		ARMType:      "Microsoft.Widgets/widgets",
		GenResource:  true,
		DataSource:   true,
	})

	// The emitted mapping must parse back into a Mapping.
	dir := t.TempDir()
	path := filepath.Join(dir, "widget.scaffold.hcl")
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		t.Fatalf("writing emitted mapping: %v", err)
	}
	m, err := Parse(path)
	if err != nil {
		t.Fatalf("parsing emitted mapping: %v", err)
	}

	if m.ResourceName != "widget" {
		t.Errorf("resource_name = %q, want widget", m.ResourceName)
	}
	if m.ARMType != "Microsoft.Widgets/widgets" {
		t.Errorf("arm_type = %q", m.ARMType)
	}
	if m.APIVersion != "2025-01-01" {
		t.Errorf("api_version = %q, want 2025-01-01", m.APIVersion)
	}
	if m.GoName != "Widget" {
		t.Errorf("go_name = %q, want Widget", m.GoName)
	}

	// Every resolved property must appear as an addressable attribute block.
	byPath := map[string]AttributeOverride{}
	for _, a := range m.Attributes {
		byPath[a.SourcePath] = a
	}
	for _, want := range []string{"name", "properties.size", "properties.profile", "properties.profile.sku", "properties.profile.tier"} {
		if _, ok := byPath[want]; !ok {
			t.Errorf("emitted mapping missing attribute %q", want)
		}
	}
	// Block properties are emitted with a path-qualified tf_name so the user can
	// see the nesting context and use the value directly as a path-qualified rename.
	if got := byPath["properties.profile.sku"]; got.TFName == nil || *got.TFName != "profile.sku" {
		t.Errorf("sku tf_name not round-tripped: %+v", got)
	}
}

// TestEmitApplyRoundTrip proves that emitting then re-applying the mapping is a
// no-op: the schema is unchanged when nothing is edited.
func TestEmitApplyRoundTrip(t *testing.T) {
	res := fixtureIR()
	out := Emit(res, EmitOptions{ResourceName: "widget", ARMType: "Microsoft.Widgets/widgets", GenResource: true})

	dir := t.TempDir()
	path := filepath.Join(dir, "widget.scaffold.hcl")
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		t.Fatalf("writing emitted mapping: %v", err)
	}
	m, err := Parse(path)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}

	fresh := fixtureIR()
	warnings, _, err := Apply(fresh, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings on no-op apply: %v", warnings)
	}
	// Names and fields must be identical to the untouched fixture.
	orig := fixtureIR()
	if len(fresh.TopLevel) != len(orig.TopLevel) || len(fresh.Blocks) != len(orig.Blocks) {
		t.Fatalf("no-op apply changed shape: top=%d/%d blocks=%d/%d", len(fresh.TopLevel), len(orig.TopLevel), len(fresh.Blocks), len(orig.Blocks))
	}
	for i, p := range fresh.TopLevel {
		if p.TFName != orig.TopLevel[i].TFName || p.GoField != orig.TopLevel[i].GoField {
			t.Errorf("no-op apply changed %q", p.SourcePath)
		}
	}
}
