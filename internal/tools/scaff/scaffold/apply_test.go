// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package scaffold
// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package scaffold

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

func strp(s string) *string { return &s }
func boolp(b bool) *bool    { return &b }

// fixtureIR returns a small resolved-IR-like structure for exercising Apply.
func fixtureIR() *ir.ResourceIR {
	return &ir.ResourceIR{
		Name:           "Widget",
		TerraformType:  "azurerm_widget",
		ServicePackage: "widgets",
		APIVersion:     "2025-01-01",
		SDKPackage:     "widgets",
		TopLevel: []*ir.Property{
			{TFName: "name", GoField: "Name", SDKField: "Name", SourcePath: "name", Required: true, ForceNew: true, TFType: "TypeString", GoType: "string"},
			{TFName: "size", GoField: "Size", SDKField: "Size", SourcePath: "properties.size", Optional: true, TFType: "TypeString", GoType: "string"},
			{TFName: "profile", GoField: "Profile", SourcePath: "properties.profile", IsBlock: true, BlockName: "Profile", MaxItems: 1, TFType: "TypeList", GoType: "[]Profile", Optional: true},
		},
		Blocks: []*ir.BlockModel{
			{Name: "Profile", SDKModel: "Profile", Properties: []*ir.Property{
				{TFName: "sku", GoField: "Sku", SDKField: "Sku", SourcePath: "properties.profile.sku", Optional: true, TFType: "TypeString", GoType: "string"},
				{TFName: "tier", GoField: "Tier", SDKField: "Tier", SourcePath: "properties.profile.tier", Optional: true, TFType: "TypeString", GoType: "string"},
			}},
		},
	}
}

// findProp locates a property by source path across top-level and block scopes.
func findProp(res *ir.ResourceIR, sourcePath string) *ir.Property {
	for _, p := range res.TopLevel {
		if p.SourcePath == sourcePath {
			return p
		}
	}
	for _, b := range res.Blocks {
		for _, p := range b.Properties {
			if p.SourcePath == sourcePath {
				return p
			}
		}
	}
	return nil
}

func TestApply_Rename(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.profile.sku", TFName: strp("stock_keeping_unit")},
	}}
	warnings, err := Apply(res, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	p := findProp(res, "properties.profile.sku")
	if p == nil {
		t.Fatal("sku property missing after rename")
	}
	if p.TFName != "stock_keeping_unit" {
		t.Errorf("TFName = %q, want stock_keeping_unit", p.TFName)
	}
	if p.GoField != "StockKeepingUnit" {
		t.Errorf("GoField = %q, want StockKeepingUnit", p.GoField)
	}
	if p.SDKField != "Sku" {
		t.Errorf("SDKField changed to %q, want Sku (must stay for expand/flatten)", p.SDKField)
	}
}

func TestApply_RemoveScalar(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.size", Remove: boolp(true)},
	}}
	if _, err := Apply(res, m); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if findProp(res, "properties.size") != nil {
		t.Error("expected properties.size to be removed")
	}
}

func TestApply_RemoveBlockPrunes(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.profile", Remove: boolp(true)},
	}}
	if _, err := Apply(res, m); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if findProp(res, "properties.profile") != nil {
		t.Error("expected block reference to be removed")
	}
	for _, b := range res.Blocks {
		if b.Name == "Profile" {
			t.Error("expected the now-unreferenced Profile block to be pruned")
		}
	}
}

func TestApply_Flags(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.size", Required: boolp(true), Sensitive: boolp(true), ForceNew: boolp(true)},
	}}
	if _, err := Apply(res, m); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	p := findProp(res, "properties.size")
	if !p.Required || !p.Sensitive || !p.ForceNew {
		t.Errorf("flags not applied: required=%v sensitive=%v force_new=%v", p.Required, p.Sensitive, p.ForceNew)
	}
	if p.Optional {
		t.Error("expected Optional to be cleared when Required is set")
	}
}

func TestApply_RenameCollisionErrors(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.profile.tier", TFName: strp("sku")}, // collides with existing sku
	}}
	_, err := Apply(res, m)
	if err == nil {
		t.Fatal("expected a conflict error for the colliding rename")
	}
	if !strings.Contains(err.Error(), "duplicate schema key") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestApply_GuardIDCritical(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "name", Remove: boolp(true)},
	}}
	warnings, err := Apply(res, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if findProp(res, "name") == nil {
		t.Error("name must not be removed (ID-critical)")
	}
	if len(warnings) == 0 || !strings.Contains(strings.Join(warnings, " "), "cannot remove") {
		t.Errorf("expected a guard warning, got %v", warnings)
	}
}

func TestApply_UnknownPathWarns(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.doesnotexist", TFName: strp("nope")},
	}}
	warnings, err := Apply(res, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(warnings) == 0 || !strings.Contains(strings.Join(warnings, " "), "matches no attribute") {
		t.Errorf("expected an unmatched-path warning, got %v", warnings)
	}
}
