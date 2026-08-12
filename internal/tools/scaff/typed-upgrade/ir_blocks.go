// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package typed_upgrade

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/gen"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/pandora"
)

// BlockResolveOptions carries the optional Pandora connection + target details
// needed to resolve the full resource IR for nested block generation.
type BlockResolveOptions struct {
	PandoraURL string
	ARMType    string
	Service    string
	Resource   string
	APIVersion string
}

// ResolveBlocks optionally resolves the Pandora IR and stores it on info,
// populating info.ExpandFlattenSrc with the generated expand/flatten helpers.
// It is a no-op (returns nil) when no ARM target is specified or the server is
// unreachable. Non-fatal errors are returned so callers can warn rather than fail.
func ResolveBlocks(info *Info, opts BlockResolveOptions) error {
	if opts.ARMType == "" && (opts.Service == "" || opts.Resource == "") {
		return nil
	}
	url := opts.PandoraURL
	if url == "" {
		url = "http://localhost:8080"
	}
	client := pandora.NewClient(url)
	name := toSnake(info.BaseName)
	resolved, err := ir.Resolve(client, ir.Options{
		ARMType:        opts.ARMType,
		Service:        opts.Service,
		Resource:       opts.Resource,
		APIVersion:     opts.APIVersion,
		Name:           name,
		GoName:         info.GoName,
		ServicePackage: info.Package,
	})
	if err != nil {
		return fmt.Errorf("resolving Pandora IR for %q: %w", info.TerraformType, err)
	}
	info.PandoraIR = resolved
	info.ExpandFlattenSrc = renderBlockHelpers(info)
	return nil
}

// renderBlockHelpers generates typed expand/flatten functions for every nested
// block model in the Pandora IR that also appears in the AST-extracted schema
// fields.
func renderBlockHelpers(info *Info) string {
	if info.PandoraIR == nil || len(info.PandoraIR.Blocks) == 0 {
		return ""
	}

	// Cross-reference IR blocks with AST schema fields by TF name.
	blocksByTF := buildIRBlockByTFName(info.PandoraIR)
	matched := matchedBlocks(info.Fields, blocksByTF)
	if len(matched) == 0 {
		return ""
	}

	// Build a minimal IR for the renderer — just SDKPackage + the matched blocks.
	// Populate TopLevel so blockUsage() in gen/crud.go can determine MaxItems.
	minIR := &ir.ResourceIR{
		SDKPackage: info.SDKPackage,
		Blocks:     matched,
	}
	for _, f := range info.Fields {
		if f.Kind != FieldListObj {
			continue
		}
		blockName := irBlockNameForField(f, blocksByTF)
		if blockName == "" {
			continue
		}
		maxItems := f.MaxItems
		minIR.TopLevel = append(minIR.TopLevel, &ir.Property{
			TFName:    f.TFName,
			IsBlock:   true,
			BlockName: blockName,
			MaxItems:  maxItems,
		})
	}
	return gen.RenderExpandFlatten(minIR)
}

// buildIRBlockByTFName builds a mapping from the schema TF name (derived from
// the IR property) to the IR BlockModel.
func buildIRBlockByTFName(res *ir.ResourceIR) map[string]*ir.BlockModel {
	// First index blocks by their Go name so we can look them up from TopLevel.
	byName := make(map[string]*ir.BlockModel, len(res.Blocks))
	for _, b := range res.Blocks {
		byName[b.Name] = b
	}
	// Build TFName → BlockModel using the TopLevel properties.
	out := make(map[string]*ir.BlockModel)
	for _, p := range res.TopLevel {
		if p.IsBlock && p.BlockName != "" {
			if b, ok := byName[p.BlockName]; ok {
				out[p.TFName] = b
			}
		}
	}
	return out
}

// matchedBlocks returns the IR BlockModels whose TF name appears as a
// FieldListObj in the AST schema fields.
func matchedBlocks(fields []*SchemaField, blocksByTF map[string]*ir.BlockModel) []*ir.BlockModel {
	seen := map[string]bool{}
	var out []*ir.BlockModel
	for _, f := range fields {
		if f.Kind != FieldListObj {
			continue
		}
		if b, ok := blocksByTF[f.TFName]; ok && !seen[b.Name] {
			seen[b.Name] = true
			out = append(out, b)
		}
	}
	return out
}

// irBlockNameForField returns the IR BlockModel.Name for the given field,
// or "" when not found.
func irBlockNameForField(f *SchemaField, blocksByTF map[string]*ir.BlockModel) string {
	if b, ok := blocksByTF[f.TFName]; ok {
		return b.Name
	}
	return ""
}

// --- CRUD body rewiring for blocks ----------------------------------------

// applyBlockTransforms replaces old expand/flatten calls for nested blocks with
// the typed versions generated from the IR. It is called after the scalar model
// transform pass has already replaced Get/Set calls.
//
//   - In Create/Update: expandOldFunc(model.Field) → expandIRName(model.Field)
//   - In Read: state.Field = flattenOldFunc(sdkExpr) → state.Field = flattenIRName(sdkExpr)
func applyBlockTransforms(body string, info *Info) string {
	if info.PandoraIR == nil {
		return body
	}
	blocksByTF := buildIRBlockByTFName(info.PandoraIR)
	for _, f := range info.Fields {
		if f.Kind != FieldListObj {
			continue
		}
		b, ok := blocksByTF[f.TFName]
		if !ok {
			continue
		}
		single := f.MaxItems == 1
		expandName := expandFuncName(b.Name, single)
		flattenName := flattenFuncName(b.Name, single)
		goField := f.GoField

		// Replace old expand: expandXxx(model.GoField) → expandIRName(model.GoField)
		reExpand := regexp.MustCompile(`\bexpand\w+\(model\.` + regexp.QuoteMeta(goField) + `\)`)
		body = reExpand.ReplaceAllString(body, expandName+"(model."+goField+")")

		// Replace old expand without model prefix (direct old-style call):
		// expandXxx(metadata.ResourceData.Get("tfName").([]interface{}))
		// was already handled by Get replacement but catch any residuals.
		reExpandDirect := regexp.MustCompile(`\bexpand\w+\(metadata\.ResourceData\.Get\("` +
			regexp.QuoteMeta(f.TFName) + `"\)\.\([^)]*\)\)`)
		body = reExpandDirect.ReplaceAllString(body, expandName+"(model."+goField+")")

		// Replace old flatten in Read Set: state.GoField = flattenXxx(expr)
		// → state.GoField = flattenIRName(expr)  (preserves the inner arg)
		reFlatten := regexp.MustCompile(`(state\.` + regexp.QuoteMeta(goField) + ` = )flattenOld\w+\(`)
		body = reFlatten.ReplaceAllString(body, "${1}"+flattenName+"(")
		// Catch any flatten call assigned to state.Field regardless of prefix:
		reFlattenAny := regexp.MustCompile(
			`(state\.` + regexp.QuoteMeta(goField) + ` = )flatten\w+\(`)
		body = reFlattenAny.ReplaceAllString(body, "${1}"+flattenName+"(")
	}
	return body
}

// expandFuncName returns the name of the expand function generated by
// gen.RenderExpandFlatten for the given block.
func expandFuncName(blockName string, single bool) string {
	if single {
		return "expand" + blockName
	}
	return "expand" + blockName + "List"
}

// flattenFuncName returns the name of the flatten function generated by
// gen.RenderExpandFlatten for the given block.
func flattenFuncName(blockName string, single bool) string {
	if single {
		return "flatten" + blockName
	}
	return "flatten" + blockName + "List"
}

// BlockModelGoType returns the Go type for a block field's model struct.
func BlockModelGoType(f *SchemaField) string {
	if f.MaxItems == 1 {
		return "[]" + f.BlockName // typed SDK uses []Struct even for MaxItems=1
	}
	return "[]" + f.BlockName
}

// irBlockNameToModelName converts an IR block name to the model struct name
// used in the generated typed resource. The IR uses plain names like "Admin";
// we append "Model" to avoid conflicts (e.g. AdminModel).
func irBlockNameToModelName(irName string) string {
	return irName + "Model"
}

// updateModelFieldsFromIR updates SchemaField.BlockName to match the IR block
// model name. This ensures the model struct name and the expand/flatten types
// agree.
func updateModelFieldsFromIR(fields []*SchemaField, blocksByTF map[string]*ir.BlockModel) {
	for _, f := range fields {
		if f.Kind != FieldListObj {
			continue
		}
		if b, ok := blocksByTF[f.TFName]; ok {
			f.BlockName = b.Name // Use IR name, not our AST-derived guess
		}
	}
}

// PlanBlockWarnings adds warnings to the Plan when blocks exist but no IR is available.
func PlanBlockWarnings(info *Info) []string {
	var w []string
	for _, f := range info.Fields {
		if f.Kind == FieldListObj {
			if info.PandoraIR == nil {
				w = append(w, fmt.Sprintf(
					"nested block %q detected — pass -arm-type to generate typed expand/flatten helpers",
					f.TFName))
				return w // one warning is enough
			}
			// Check if the block is matched in the IR.
			blocksByTF := buildIRBlockByTFName(info.PandoraIR)
			if _, ok := blocksByTF[f.TFName]; !ok {
				w = append(w, fmt.Sprintf(
					"nested block %q was not found in the Pandora IR; expand/flatten will need manual implementation",
					f.TFName))
			}
		}
	}
	return w
}

// IRBlockCount returns the number of nested block fields matched in the Pandora IR.
func IRBlockCount(info *Info) int {
	if info.PandoraIR == nil {
		return 0
	}
	blocksByTF := buildIRBlockByTFName(info.PandoraIR)
	n := 0
	for _, f := range info.Fields {
		if f.Kind == FieldListObj {
			if _, ok := blocksByTF[f.TFName]; ok {
				n++
			}
		}
	}
	return n
}

// nestedBlockFields returns all FieldListObj fields with a matching IR block.
func nestedBlockFields(info *Info) []*SchemaField {
	if info.PandoraIR == nil {
		return nil
	}
	blocksByTF := buildIRBlockByTFName(info.PandoraIR)
	var out []*SchemaField
	for _, f := range info.Fields {
		if f.Kind == FieldListObj {
			if _, ok := blocksByTF[f.TFName]; ok {
				out = append(out, f)
			}
		}
	}
	return out
}

// formatBlockNote generates the inline NOTE comment for a block field that
// still needs wiring in the generated body.
func formatBlockNote(f *SchemaField, single bool) string {
	return fmt.Sprintf("/* TODO: block %q: expand=%s, flatten=%s */",
		f.TFName, expandFuncName(f.BlockName, single), flattenFuncName(f.BlockName, single))
}

// renderBlockModelStructs generates the block model struct definitions for
// blocks resolved from the Pandora IR. These replace the placeholder stubs
// emitted when no IR is available.
func renderBlockModelStructs(info *Info) string {
	if info.PandoraIR == nil {
		return ""
	}
	blocksByTF := buildIRBlockByTFName(info.PandoraIR)
	seen := map[string]bool{}
	var sb strings.Builder
	for _, f := range info.Fields {
		if f.Kind != FieldListObj {
			continue
		}
		b, ok := blocksByTF[f.TFName]
		if !ok || seen[b.Name] {
			continue
		}
		seen[b.Name] = true
		fmt.Fprintf(&sb, "type %s struct {\n", b.Name)
		for _, p := range b.Properties {
			fmt.Fprintf(&sb, "\t%s %s `tfschema:%q`\n", p.GoField, p.GoType, p.TFName)
		}
		sb.WriteString("}\n\n")
	}
	return sb.String()
}
