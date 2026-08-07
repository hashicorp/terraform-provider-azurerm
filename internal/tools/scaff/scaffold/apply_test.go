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
	warnings, _, err := Apply(res, m)
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
	if _, _, err := Apply(res, m); err != nil {
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
	if _, _, err := Apply(res, m); err != nil {
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
	if _, _, err := Apply(res, m); err != nil {
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
	_, _, err := Apply(res, m)
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
	warnings, _, err := Apply(res, m)
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
	warnings, _, err := Apply(res, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(warnings) == 0 || !strings.Contains(strings.Join(warnings, " "), "matches no attribute") {
		t.Errorf("expected an unmatched-path warning, got %v", warnings)
	}
}

// ---------------------------------------------------------------------------
// Path-qualified tf_name tests
// ---------------------------------------------------------------------------

// TestApply_PathRename_MovesToExistingBlock verifies that a dotted tf_name
// whose first segment matches an existing block's schema key moves the property
// into that block and renames it to the leaf.
//
//	attribute "properties.profile.sku" { tf_name = "profile.stock_keeping_unit" }
//
// Expected: "sku" is removed from Profile.Properties and re-inserted as
// "stock_keeping_unit"; no new block is created.
func TestApply_PathRename_MovesToExistingBlock(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.profile.sku", TFName: strp("profile.stock_keeping_unit")},
	}}
	_, _, err := Apply(res, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// The property must still be findable by source path.
	p := findProp(res, "properties.profile.sku")
	if p == nil {
		t.Fatal("sku property missing after path rename")
	}
	if p.TFName != "stock_keeping_unit" {
		t.Errorf("TFName = %q, want stock_keeping_unit", p.TFName)
	}
	if p.GoField != "StockKeepingUnit" {
		t.Errorf("GoField = %q, want StockKeepingUnit", p.GoField)
	}

	// The property must live inside the Profile block.
	profileBlock := findBlockByName(res, "Profile")
	if profileBlock == nil {
		t.Fatal("Profile block missing")
	}
	found := false
	for _, bp := range profileBlock.Properties {
		if bp.SourcePath == "properties.profile.sku" {
			found = true
		}
	}
	if !found {
		t.Error("renamed property not found inside Profile block")
	}

	// No new blocks should have been created.
	if len(res.Blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(res.Blocks))
	}
}

// TestApply_PathRename_CreatesNewBlock verifies that a dotted tf_name whose
// first segment does not match any existing block creates a synthetic block.
//
//	attribute "properties.size" { tf_name = "node.virtual_machine_sku" }
//
// Expected:
//   - "size" is removed from TopLevel.
//   - A new synthetic "Node" block exists in res.Blocks.
//   - A container property {TFName:"node", IsBlock:true, BlockName:"Node", SourcePath:""} exists in TopLevel.
//   - "virtual_machine_sku" lives inside the Node block.
func TestApply_PathRename_CreatesNewBlock(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.size", TFName: strp("node.virtual_machine_sku")},
	}}
	_, _, err := Apply(res, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// Original top-level "size" must be gone.
	for _, p := range res.TopLevel {
		if p.SourcePath == "properties.size" && !p.IsBlock {
			t.Error("scalar 'size' property must have been removed from TopLevel")
		}
	}

	// Synthetic container property at top level.
	var containerProp *ir.Property
	for _, p := range res.TopLevel {
		if p.TFName == "node" {
			containerProp = p
			break
		}
	}
	if containerProp == nil {
		t.Fatal("synthetic container property 'node' not found in TopLevel")
	}
	if !containerProp.IsBlock {
		t.Error("container property must have IsBlock == true")
	}
	if containerProp.BlockName != "Node" {
		t.Errorf("BlockName = %q, want Node", containerProp.BlockName)
	}
	if containerProp.TFType != "TypeList" {
		t.Errorf("TFType = %q, want TypeList", containerProp.TFType)
	}
	if containerProp.MaxItems != 1 {
		t.Errorf("MaxItems = %d, want 1", containerProp.MaxItems)
	}
	if containerProp.GoType != "[]Node" {
		t.Errorf("GoType = %q, want []Node", containerProp.GoType)
	}
	if !containerProp.Optional {
		t.Error("container property must be Optional")
	}
	if containerProp.SourcePath != "" {
		t.Errorf("container property SourcePath must be empty, got %q", containerProp.SourcePath)
	}

	// New Node block must exist.
	nodeBlock := findBlockByName(res, "Node")
	if nodeBlock == nil {
		t.Fatal("synthetic 'Node' block not found in res.Blocks")
	}
	if nodeBlock.SDKModel != "" {
		t.Errorf("SDKModel = %q, want empty for synthetic block", nodeBlock.SDKModel)
	}

	// Leaf property must live inside Node.
	p := findProp(res, "properties.size")
	if p == nil {
		t.Fatal("moved property not findable by source path")
	}
	if p.TFName != "virtual_machine_sku" {
		t.Errorf("TFName = %q, want virtual_machine_sku", p.TFName)
	}
	if p.GoField != "VirtualMachineSku" {
		t.Errorf("GoField = %q, want VirtualMachineSku", p.GoField)
	}
	foundInBlock := false
	for _, bp := range nodeBlock.Properties {
		if bp.SourcePath == "properties.size" {
			foundInBlock = true
		}
	}
	if !foundInBlock {
		t.Error("leaf property not found inside Node block")
	}
}

// TestApply_PathRename_DeepPath verifies that a three-segment tf_name creates
// nested synthetic blocks A → B and places the leaf inside B.
//
//	attribute "properties.size" { tf_name = "a.b.leaf" }
func TestApply_PathRename_DeepPath(t *testing.T) {
	res := fixtureIR()
	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.size", TFName: strp("a.b.leaf")},
	}}
	_, _, err := Apply(res, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// Block A must exist at top level.
	var containerA *ir.Property
	for _, p := range res.TopLevel {
		if p.TFName == "a" && p.IsBlock {
			containerA = p
			break
		}
	}
	if containerA == nil {
		t.Fatal("synthetic container property 'a' not found in TopLevel")
	}
	blockA := findBlockByName(res, containerA.BlockName)
	if blockA == nil {
		t.Fatalf("block %q not found", containerA.BlockName)
	}

	// Block B must exist inside A.
	var containerB *ir.Property
	for _, p := range blockA.Properties {
		if p.TFName == "b" && p.IsBlock {
			containerB = p
			break
		}
	}
	if containerB == nil {
		t.Fatal("synthetic container property 'b' not found inside block A")
	}
	blockB := findBlockByName(res, containerB.BlockName)
	if blockB == nil {
		t.Fatalf("block %q not found", containerB.BlockName)
	}

	// Leaf must live inside B.
	p := findProp(res, "properties.size")
	if p == nil {
		t.Fatal("moved property not findable by source path")
	}
	if p.TFName != "leaf" {
		t.Errorf("TFName = %q, want leaf", p.TFName)
	}
	foundInB := false
	for _, bp := range blockB.Properties {
		if bp.SourcePath == "properties.size" {
			foundInB = true
		}
	}
	if !foundInB {
		t.Error("leaf property not found inside block B")
	}
}

// TestApply_PathRename_FlatUnchanged verifies that a plain tf_name with no dot
// still renames in place without any block movement — the existing behaviour is
// completely preserved.
func TestApply_PathRename_FlatUnchanged(t *testing.T) {
	res := fixtureIR()
	origTopLevelCount := len(res.TopLevel)
	origBlockCount := len(res.Blocks)

	m := &Mapping{Attributes: []AttributeOverride{
		{SourcePath: "properties.size", TFName: strp("vm_sku")},
	}}
	_, _, err := Apply(res, m)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// Property must still be at top level, renamed.
	p := findProp(res, "properties.size")
	if p == nil {
		t.Fatal("property missing after flat rename")
	}
	if p.TFName != "vm_sku" {
		t.Errorf("TFName = %q, want vm_sku", p.TFName)
	}
	if p.GoField != "VmSku" {
		t.Errorf("GoField = %q, want VmSku", p.GoField)
	}

	// No structural changes: same counts.
	if len(res.TopLevel) != origTopLevelCount {
		t.Errorf("TopLevel count changed: got %d, want %d", len(res.TopLevel), origTopLevelCount)
	}
	if len(res.Blocks) != origBlockCount {
		t.Errorf("Blocks count changed: got %d, want %d", len(res.Blocks), origBlockCount)
	}
}
