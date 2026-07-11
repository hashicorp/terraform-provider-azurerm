// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gen

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

// blockByName indexes a resource's nested block models by their Go struct name.
func blockByName(res *ir.ResourceIR) map[string]*ir.BlockModel {
	m := make(map[string]*ir.BlockModel, len(res.Blocks))
	for _, b := range res.Blocks {
		m[b.Name] = b
	}
	return m
}

// isAttribute reports whether a property is a read-only (computed) attribute
// rather than a configurable argument.
func isAttribute(p *ir.Property) bool {
	return p.Computed && !p.Required && !p.Optional
}

// RenderModelStructs renders the top-level model struct and every nested block
// struct as Go source.
func RenderModelStructs(res *ir.ResourceIR) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "type %s struct {\n", res.ModelStructName)
	for _, p := range res.TopLevel {
		fmt.Fprintf(&sb, "\t%s %s `tfschema:%q`\n", p.GoField, p.GoType, p.TFName)
	}
	sb.WriteString("}\n\n")

	for _, b := range res.Blocks {
		fmt.Fprintf(&sb, "type %s struct {\n", b.Name)
		for _, p := range b.Properties {
			fmt.Fprintf(&sb, "\t%s %s `tfschema:%q`\n", p.GoField, p.GoType, p.TFName)
		}
		sb.WriteString("}\n\n")
	}

	return sb.String()
}

// RenderArguments renders the body of the Arguments() schema map (configurable
// properties).
func RenderArguments(res *ir.ResourceIR) string {
	blocks := blockByName(res)
	var sb strings.Builder
	for _, p := range res.TopLevel {
		if isAttribute(p) {
			continue
		}
		renderSchemaProp(&sb, p, blocks, false)
	}
	return sb.String()
}

// RenderAttributes renders the body of the Attributes() schema map (computed
// properties).
func RenderAttributes(res *ir.ResourceIR) string {
	blocks := blockByName(res)
	var sb strings.Builder
	for _, p := range res.TopLevel {
		if !isAttribute(p) {
			continue
		}
		renderSchemaProp(&sb, p, blocks, false)
	}
	return sb.String()
}

// renderSchemaProp renders a single schema entry (recursing into nested blocks).
// When computed is true every attribute is emitted as Computed (data source
// mode). Output is intentionally minimally-indented; gofmt normalises it.
func renderSchemaProp(sb *strings.Builder, p *ir.Property, blocks map[string]*ir.BlockModel, computed bool) {
	// Common schema helpers for the well-known top-level properties.
	switch p.TFName {
	case "location":
		if computed {
			sb.WriteString("\"location\": commonschema.LocationComputed(),\n")
		} else {
			sb.WriteString("\"location\": commonschema.Location(),\n")
		}
		return
	case "resource_group_name":
		sb.WriteString("\"resource_group_name\": commonschema.ResourceGroupName(),\n")
		return
	case "tags":
		if computed {
			sb.WriteString("\"tags\": commonschema.TagsDataSource(),\n")
		} else {
			sb.WriteString("\"tags\": commonschema.Tags(),\n")
		}
		return
	}

	fmt.Fprintf(sb, "%q: {\n", p.TFName)
	fmt.Fprintf(sb, "Type: pluginsdk.%s,\n", p.TFType)
	if p.Description != "" {
		fmt.Fprintf(sb, "Description: %q,\n", p.Description)
	}

	if computed {
		sb.WriteString("Computed: true,\n")
	} else {
		switch {
		case p.Required:
			sb.WriteString("Required: true,\n")
		case p.Optional:
			sb.WriteString("Optional: true,\n")
		}
		if p.Computed {
			sb.WriteString("Computed: true,\n")
		}
		if p.ForceNew {
			sb.WriteString("ForceNew: true,\n")
		}
		if p.Sensitive {
			sb.WriteString("Sensitive: true,\n")
		}
	}
	if p.MaxItems == 1 {
		sb.WriteString("MaxItems: 1,\n")
	}

	switch {
	case p.IsBlock:
		block := blocks[p.BlockName]
		sb.WriteString("Elem: &pluginsdk.Resource{\n")
		sb.WriteString("Schema: map[string]*pluginsdk.Schema{\n")
		if block != nil {
			for _, child := range block.Properties {
				renderSchemaProp(sb, child, blocks, computed)
			}
		}
		sb.WriteString("},\n")
		sb.WriteString("},\n")
	case p.TFType == "TypeList":
		// list of primitives (or enum values); any validation applies to the
		// element schema rather than the list itself.
		itemType := listItemTFType(p.GoType)
		sb.WriteString("Elem: &pluginsdk.Schema{\n")
		fmt.Fprintf(sb, "Type: pluginsdk.%s,\n", itemType)
		if isConfigurable(p, computed) {
			if vf := scalarValidateFunc(itemType, p.IsEnum, p.EnumValues); vf != "" {
				fmt.Fprintf(sb, "ValidateFunc: %s,\n", vf)
			}
		}
		sb.WriteString("},\n")
	default:
		// Scalar primitives: always attach a default ValidateFunc for
		// configurable, non-boolean arguments (e.g. StringIsNotEmpty,
		// IntAtLeast(0)). Booleans and maps have no meaningful default.
		if isConfigurable(p, computed) {
			if vf := scalarValidateFunc(p.TFType, p.IsEnum, p.EnumValues); vf != "" {
				fmt.Fprintf(sb, "ValidateFunc: %s,\n", vf)
			}
		}
	}

	sb.WriteString("},\n")
}

// isConfigurable reports whether a rendered property accepts user input and so
// should carry a ValidateFunc. In data-source (computed) mode nothing is
// configurable; otherwise a property is configurable when it is Required or
// Optional (computed-only attributes are excluded).
func isConfigurable(p *ir.Property, computed bool) bool {
	return !computed && (p.Required || p.Optional)
}

// scalarValidateFunc returns the validation.* expression (without the
// "ValidateFunc:" key or trailing comma) appropriate for a scalar schema of the
// given Terraform type. Enums validate against their allowed values; strings,
// ints and floats get a permissive non-boolean default. Booleans, maps and
// anything else return "" (no validation).
func scalarValidateFunc(tfType string, isEnum bool, enumValues []string) string {
	if isEnum && tfType == "TypeString" {
		return fmt.Sprintf("validation.StringInSlice([]string{\n%s}, false)", renderEnumValues(enumValues))
	}
	switch tfType {
	case "TypeString":
		return "validation.StringIsNotEmpty"
	case "TypeInt":
		return "validation.IntAtLeast(0)"
	case "TypeFloat":
		return "validation.FloatAtLeast(0)"
	default:
		return ""
	}
}

// renderEnumValues renders a quoted, comma-separated list of enum values.
func renderEnumValues(values []string) string {
	var sb strings.Builder
	for _, v := range values {
		fmt.Fprintf(&sb, "%q,\n", v)
	}
	return sb.String()
}

// listItemTFType maps a Go slice type to the pluginsdk element type.
func listItemTFType(goType string) string {
	switch goType {
	case "[]int64":
		return "TypeInt"
	case "[]bool":
		return "TypeBool"
	case "[]float64":
		return "TypeFloat"
	default:
		return "TypeString"
	}
}
