// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package scaffold

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

// idCriticalPaths are properties whose Go field feeds New<Resource>ID(); they
// must not be removed and their Go field must stay stable, so only their schema
// key may be renamed.
var idCriticalPaths = map[string]bool{
	"name":                true,
	"resource_group_name": true,
}

// Apply transforms the resolved IR in place according to the mapping's
// per-attribute overrides, matched by source path. It returns non-fatal
// warnings (e.g. an override that matches nothing, or the removal of a required
// attribute), a map of property renames (oldTFName -> newTFName), and a fatal
// error when the result would not compile (e.g. a rename that collides with
// another attribute).
func Apply(res *ir.ResourceIR, m *Mapping) (warnings []string, renames map[string]string, err error) {
	overrides := make(map[string]*AttributeOverride, len(m.Attributes))
	for i := range m.Attributes {
		overrides[m.Attributes[i].SourcePath] = &m.Attributes[i]
	}
	seen := map[string]bool{}
	renames = make(map[string]string)

	// Collect path-qualified moves while applying flat overrides in place.
	var pendingMoves []pendingMove

	res.TopLevel = applyToPropsCollect(res.TopLevel, overrides, seen, &warnings, &pendingMoves, true, renames)
	for _, b := range res.Blocks {
		b.Properties = applyToPropsCollect(b.Properties, overrides, seen, &warnings, &pendingMoves, false, renames)
	}

	// Execute moves: insert each moved property into its resolved target block
	// (or TopLevel if blockPath is empty).
	for _, mv := range pendingMoves {
		renameProperty(mv.prop, mv.leafName, idCriticalPaths[mv.prop.SourcePath])
		if len(mv.blockPath) == 0 {
			// Move to TopLevel (promotion from nested).
			res.TopLevel = append(res.TopLevel, mv.prop)
		} else {
			// Move to a nested block.
			targetBlock := resolveOrCreateBlock(res, mv.blockPath)
			targetBlock.Properties = append(targetBlock.Properties, mv.prop)
		}
	}

	for path := range overrides {
		if !seen[path] {
			warnings = append(warnings, fmt.Sprintf("override for %q matches no attribute (removed upstream?)", path))
		}
	}

	pruneUnreferencedBlocks(res)

	if conflicts := validateUnique(res); len(conflicts) > 0 {
		return warnings, renames, fmt.Errorf("schema overrides produced conflicts:\n  - %s", strings.Join(conflicts, "\n  - "))
	}
	return warnings, renames, nil
}

// pendingMove captures a property that must be relocated into a target block
// chain after all current-location removals have been processed.
type pendingMove struct {
	prop      *ir.Property
	blockPath []string // ordered block tf_name segments, e.g. ["node"] or ["a", "b"]
	leafName  string   // final leaf schema key, e.g. "virtual_machine_sku"
}

// applyToPropsCollect applies overrides to a property slice. For flat tf_name
// values on top-level properties it renames in place. For path-qualified tf_name
// values (or flat tf_name on nested properties being promoted to top-level) it
// removes the property from this slice and appends it to moves so it can be
// placed in the correct block later. Renames are tracked in the renames map.
//
// isTopLevel distinguishes: when false and a flat tf_name is applied to a
// property with nested SourcePath (e.g. "properties.encryption.key" → "key_version"),
// the property is queued to move to TopLevel instead of being renamed in place.
func applyToPropsCollect(props []*ir.Property, overrides map[string]*AttributeOverride, seen map[string]bool, warnings *[]string, moves *[]pendingMove, isTopLevel bool, renames map[string]string) []*ir.Property {
	out := make([]*ir.Property, 0, len(props))
	for _, p := range props {
		ov, ok := overrides[p.SourcePath]
		if !ok {
			out = append(out, p)
			continue
		}
		seen[p.SourcePath] = true

		if ov.Remove != nil && *ov.Remove {
			if idCriticalPaths[p.SourcePath] {
				*warnings = append(*warnings, fmt.Sprintf("cannot remove %q: it is required for resource ID construction; keeping it", p.SourcePath))
				out = append(out, p)
				continue
			}
			if p.Required && p.SDKField != "" {
				*warnings = append(*warnings, fmt.Sprintf("removed required attribute %q; the SDK field %q will be left unset", p.SourcePath, p.SDKField))
			}
			continue // drop the property
		}

		applyFlags(p, ov)

		if ov.TFName != nil && *ov.TFName != "" {
			blockPath, leafName := splitTFNamePath(*ov.TFName)
			if len(blockPath) > 0 {
				// Path-qualified rename: remove from current location and queue move.
				oldTFName := p.TFName
				*moves = append(*moves, pendingMove{prop: p, blockPath: blockPath, leafName: leafName})
				renames[oldTFName] = leafName
				continue // do NOT append to out — property is being relocated
			}
			// Flat rename (no dot).
			// When processing a nested property (not TopLevel) with a flat tf_name,
			// it should be promoted to TopLevel unless the SourcePath is also flat.
			if !isTopLevel && strings.Contains(p.SourcePath, ".") {
				// Promote nested property to TopLevel: queue it as a move with empty blockPath.
				oldTFName := p.TFName
				*moves = append(*moves, pendingMove{prop: p, blockPath: []string{}, leafName: leafName})
				renames[oldTFName] = leafName
				continue // do NOT append to out — property is being relocated
			}
			// Flat rename at TopLevel or on a top-level source property: rename in place.
			if leafName != p.TFName {
				oldTFName := p.TFName
				renameProperty(p, leafName, idCriticalPaths[p.SourcePath])
				renames[oldTFName] = leafName
			}
		}
		out = append(out, p)
	}
	return out
}

// splitTFNamePath splits a tf_name value on "." into block-path segments and
// the leaf name. A plain name with no "." returns blockPath == nil.
func splitTFNamePath(tfName string) (blockPath []string, leafName string) {
	parts := strings.Split(tfName, ".")
	if len(parts) == 1 {
		return nil, tfName
	}
	return parts[:len(parts)-1], parts[len(parts)-1]
}

// resolveOrCreateBlock walks (creating synthetic blocks as needed) the block
// chain described by blockPath and returns the innermost BlockModel. Each
// segment is the tf_name (schema key) of the container property at its level.
func resolveOrCreateBlock(res *ir.ResourceIR, blockPath []string) *ir.BlockModel {
	// We walk level-by-level. At the top level we search res.TopLevel; once
	// inside a block we search that block's Properties. Create as needed.

	// currentProps points to the property slice of the current scope.
	// We use a pointer-to-slice so we can append to it.
	currentProps := &res.TopLevel

	// Index blocks by Name for fast lookup; updated lazily as we may add blocks.
	var current *ir.BlockModel

	for i, segment := range blockPath {
		_ = i
		// Find an existing container property for this segment in currentProps.
		var containerProp *ir.Property
		for _, p := range *currentProps {
			if p.TFName == segment && p.IsBlock {
				containerProp = p
				break
			}
		}

		if containerProp == nil {
			// Create a new synthetic block and its container property.
			goName := ir.GoFieldFromTFName(segment)
			newBlock := &ir.BlockModel{
				Name:       goName,
				SDKModel:   "",
				Properties: []*ir.Property{},
			}
			res.Blocks = append(res.Blocks, newBlock)

			newProp := &ir.Property{
				TFName:     segment,
				GoField:    goName,
				TFType:     "TypeList",
				MaxItems:   1,
				IsBlock:    true,
				BlockName:  goName,
				GoType:     "[]" + goName,
				Optional:   true,
				SourcePath: "", // synthetic; no API source path
			}
			*currentProps = append(*currentProps, newProp)
			current = newBlock
		} else {
			// Find the existing block by name.
			current = findBlockByName(res, containerProp.BlockName)
		}

		// Next iteration operates inside the current block's properties.
		currentProps = &current.Properties
	}
	return current
}

// findBlockByName returns the BlockModel with the given Name, or nil.
func findBlockByName(res *ir.ResourceIR, name string) *ir.BlockModel {
	for _, b := range res.Blocks {
		if b.Name == name {
			return b
		}
	}
	return nil
}

// applyFlags applies the metadata overrides that are set on the override.
func applyFlags(p *ir.Property, ov *AttributeOverride) {
	if ov.Required != nil {
		p.Required = *ov.Required
		if p.Required {
			p.Optional = false
			p.Computed = false
		}
	}
	if ov.Optional != nil {
		p.Optional = *ov.Optional
		if p.Optional {
			p.Required = false
		}
	}
	if ov.Computed != nil {
		p.Computed = *ov.Computed
	}
	if ov.Sensitive != nil {
		p.Sensitive = *ov.Sensitive
	}
	if ov.ForceNew != nil {
		p.ForceNew = *ov.ForceNew
	}
	if ov.Description != nil {
		p.Description = *ov.Description
	}
}

// renameProperty changes a property's schema key, keeping the Go struct field in
// sync (via the same rule the resolver uses) unless the property is ID-critical,
// in which case only the schema key changes so New<Resource>ID() still resolves.
func renameProperty(p *ir.Property, tfName string, idCritical bool) {
	p.TFName = tfName
	if !idCritical {
		p.GoField = ir.GoFieldFromTFName(tfName)
	}
}

// pruneUnreferencedBlocks drops nested block models that are no longer reachable
// from the top-level schema after removals, preserving order.
func pruneUnreferencedBlocks(res *ir.ResourceIR) {
	byName := make(map[string]*ir.BlockModel, len(res.Blocks))
	for _, b := range res.Blocks {
		byName[b.Name] = b
	}
	referenced := map[string]bool{}
	var mark func(props []*ir.Property)
	mark = func(props []*ir.Property) {
		for _, p := range props {
			if p.IsBlock && p.BlockName != "" && !referenced[p.BlockName] {
				referenced[p.BlockName] = true
				if b, ok := byName[p.BlockName]; ok {
					mark(b.Properties)
				}
			}
		}
	}
	mark(res.TopLevel)

	out := make([]*ir.BlockModel, 0, len(res.Blocks))
	for _, b := range res.Blocks {
		if referenced[b.Name] {
			out = append(out, b)
		}
	}
	res.Blocks = out
}

// validateUnique reports duplicate schema keys or Go fields within a struct,
// which would occur if a rename collided with a sibling attribute.
func validateUnique(res *ir.ResourceIR) []string {
	var conflicts []string
	check := func(scope string, props []*ir.Property) {
		tf := map[string]bool{}
		gf := map[string]bool{}
		for _, p := range props {
			if tf[p.TFName] {
				conflicts = append(conflicts, fmt.Sprintf("%s: duplicate schema key %q", scope, p.TFName))
			}
			tf[p.TFName] = true
			if gf[p.GoField] {
				conflicts = append(conflicts, fmt.Sprintf("%s: duplicate Go field %q", scope, p.GoField))
			}
			gf[p.GoField] = true
		}
	}
	check("top level", res.TopLevel)
	for _, b := range res.Blocks {
		check("block "+b.Name, b.Properties)
	}
	return conflicts
}
