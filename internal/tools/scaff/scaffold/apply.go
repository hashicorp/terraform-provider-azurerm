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
// attribute) and a fatal error when the result would not compile (e.g. a rename
// that collides with another attribute).
func Apply(res *ir.ResourceIR, m *Mapping) (warnings []string, err error) {
	overrides := make(map[string]*AttributeOverride, len(m.Attributes))
	for i := range m.Attributes {
		overrides[m.Attributes[i].SourcePath] = &m.Attributes[i]
	}
	seen := map[string]bool{}

	res.TopLevel = applyToProps(res.TopLevel, overrides, seen, &warnings)
	for _, b := range res.Blocks {
		b.Properties = applyToProps(b.Properties, overrides, seen, &warnings)
	}

	for path := range overrides {
		if !seen[path] {
			warnings = append(warnings, fmt.Sprintf("override for %q matches no attribute (removed upstream?)", path))
		}
	}

	pruneUnreferencedBlocks(res)

	if conflicts := validateUnique(res); len(conflicts) > 0 {
		return warnings, fmt.Errorf("schema overrides produced conflicts:\n  - %s", strings.Join(conflicts, "\n  - "))
	}
	return warnings, nil
}

// applyToProps applies overrides to a property slice, returning it with removed
// properties filtered out.
func applyToProps(props []*ir.Property, overrides map[string]*AttributeOverride, seen map[string]bool, warnings *[]string) []*ir.Property {
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
		if ov.TFName != nil && *ov.TFName != "" && *ov.TFName != p.TFName {
			renameProperty(p, *ov.TFName, idCriticalPaths[p.SourcePath])
		}
		out = append(out, p)
	}
	return out
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
