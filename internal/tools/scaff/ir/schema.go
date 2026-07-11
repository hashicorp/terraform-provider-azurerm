// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ir

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/pandora"
)

// envelope fields that are metadata and never surfaced as configurable schema.
var skipEnvelopeFields = map[string]bool{
	"Id":                true,
	"Type":              true,
	"SystemData":        true,
	"Etag":              true,
	"ProvisioningState": true,
}

// goIdent normalises an SDK identifier to the provider's Go field/struct casing
// (e.g. VMSize -> VmSize, APIServerProfile -> ApiServerProfile, DiskSizeGB ->
// DiskSizeGb) by round-tripping through snake_case.
func goIdent(s string) string {
	return camel(snake(s))
}

// walkSchema traverses the create model's graph, flattens the ARM envelope
// (hoisting name/location/tags, descending into Properties), and populates the
// IR's top-level properties and nested block models.
func (res *ResourceIR) walkSchema(schema *pandora.SchemaResponse) error {
	createModel, ok := schema.Models[res.CreateModel]
	if !ok {
		return fmt.Errorf("create model %q not present in schema", res.CreateModel)
	}

	w := &schemaWalker{
		schema:   schema,
		blocks:   map[string]*BlockModel{},
		visiting: map[string]bool{},
	}

	// The resource group is sourced from the ID, not the model body.
	if res.HasResourceGroup {
		res.TopLevel = append(res.TopLevel, &Property{
			TFName:   "resource_group_name",
			TFType:   "TypeString",
			Required: true,
			ForceNew: true,
			GoField:  "ResourceGroup",
			GoType:   "string",
		})
	}

	for _, key := range sortedFieldKeys(createModel.Fields) {
		f := createModel.Fields[key]
		switch key {
		case "Name":
			res.TopLevel = append(res.TopLevel, &Property{
				TFName: "name", TFType: "TypeString", Required: true, ForceNew: true,
				GoField: "Name", GoType: "string", SDKField: "Name", JSONName: f.JSONName,
				Description: f.Description,
			})
		case "Location":
			res.TopLevel = append(res.TopLevel, &Property{
				TFName: "location", TFType: "TypeString", Required: true, ForceNew: true,
				GoField: "Location", GoType: "string", SDKField: "Location", JSONName: f.JSONName,
			})
		case "Tags":
			res.TopLevel = append(res.TopLevel, &Property{
				TFName: "tags", TFType: "TypeMap", Optional: true,
				GoField: "Tags", GoType: "map[string]string", SDKField: "Tags", JSONName: f.JSONName,
			})
		case "Identity":
			// Identity (UserAssignedIdentityMap) needs dedicated commonschema
			// support; skip for now rather than emit an unmappable block.
			continue
		case "Properties":
			propsModel, err := w.referencedModel(f.ObjectDefinition)
			if err != nil {
				return fmt.Errorf("resolving Properties: %w", err)
			}
			if f.ObjectDefinition.ReferenceName != nil {
				res.PropertiesModel = *f.ObjectDefinition.ReferenceName
			}
			for _, pk := range sortedFieldKeys(propsModel.Fields) {
				if skipEnvelopeFields[pk] {
					continue
				}
				if prop := w.resolveField(pk, propsModel.Fields[pk]); prop != nil {
					prop.UnderProperties = true
					res.TopLevel = append(res.TopLevel, prop)
				}
			}
		default:
			if skipEnvelopeFields[key] {
				continue
			}
			if prop := w.resolveField(key, f); prop != nil {
				res.TopLevel = append(res.TopLevel, prop)
			}
		}
	}

	// Emit blocks in a stable, alphabetical order.
	names := make([]string, 0, len(w.blocks))
	for name := range w.blocks {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		res.Blocks = append(res.Blocks, w.blocks[name])
	}

	return nil
}

// resolveUpdateModel inspects the update request model to determine what can be
// updated in place: whether it carries Tags, the name of its Properties
// sub-model, and the set of SDK field names present under Properties. This keeps
// Update() generation faithful to the API's PATCH surface.
func (res *ResourceIR) resolveUpdateModel(schema *pandora.SchemaResponse) {
	res.UpdatablePropSDKFields = map[string]bool{}
	if res.UpdateModel == "" {
		return
	}
	model, ok := schema.Models[res.UpdateModel]
	if !ok {
		return
	}
	for key, f := range model.Fields {
		switch key {
		case "Tags":
			res.UpdateHasTags = true
		case "Properties":
			if f.ObjectDefinition.Type == "Reference" && f.ObjectDefinition.ReferenceName != nil {
				res.UpdatePropertiesModel = *f.ObjectDefinition.ReferenceName
				if props, ok := schema.Models[res.UpdatePropertiesModel]; ok {
					for pk := range props.Fields {
						res.UpdatablePropSDKFields[pk] = true
					}
				}
			}
		}
	}
}

type schemaWalker struct {
	schema   *pandora.SchemaResponse
	blocks   map[string]*BlockModel
	visiting map[string]bool
}

// resolveField converts a single Pandora field into an IR Property, recursing
// into referenced models to build nested block structs.
func (w *schemaWalker) resolveField(key string, f pandora.Field) *Property {
	p := &Property{
		TFName:      snake(key),
		GoField:     goIdent(key),
		SDKField:    key,
		JSONName:    f.JSONName,
		Description: f.Description,
		Required:    f.Required,
		Optional:    f.Optional && !f.Required,
		Computed:    f.ReadOnly,
		Sensitive:   f.Sensitive,
	}
	if !p.Required && !p.Optional && !p.Computed {
		p.Optional = true
	}

	if !w.applyType(p, f.ObjectDefinition) {
		return nil
	}
	return p
}

// applyType maps a Pandora ObjectDefinition onto the Terraform + Go types of a
// Property, recursing into referenced models and list item types. It returns
// false when the type cannot yet be faithfully represented, in which case the
// field is skipped rather than emitting code that will not compile.
func (w *schemaWalker) applyType(p *Property, def pandora.ObjectDefinition) bool {
	switch def.Type {
	case "String", "DateTime", "RawFile", "Base64":
		p.TFType, p.GoType = "TypeString", "string"
	case "Integer":
		p.TFType, p.GoType = "TypeInt", "int64"
	case "Boolean":
		p.TFType, p.GoType = "TypeBool", "bool"
	case "Float":
		p.TFType, p.GoType = "TypeFloat", "float64"
	case "Location":
		p.TFType, p.GoType = "TypeString", "string"
	case "Tags":
		p.TFType, p.GoType = "TypeMap", "map[string]string"
	case "Dictionary":
		return w.applyDictionary(p, def)
	case "Reference":
		return w.applyReference(p, def)
	case "List":
		return w.applyList(p, def)
	default:
		// Unsupported/unknown types (e.g. UserAssignedIdentityMap) are skipped.
		return false
	}
	return true
}

// applyDictionary handles Type == "Dictionary". Maps of primitives become
// TypeMap; maps of objects are not yet representable and are skipped.
func (w *schemaWalker) applyDictionary(p *Property, def pandora.ObjectDefinition) bool {
	if def.NestedItem != nil && def.NestedItem.Type == "Reference" && def.NestedItem.ReferenceName != nil {
		if _, isConst := w.schema.Constants[*def.NestedItem.ReferenceName]; !isConst {
			return false
		}
	}
	p.TFType, p.GoType = "TypeMap", "map[string]string"
	return true
}

// applyReference handles Type == "Reference", distinguishing enum constants from
// nested object models.
func (w *schemaWalker) applyReference(p *Property, def pandora.ObjectDefinition) bool {
	if def.ReferenceName == nil {
		p.TFType, p.GoType = "TypeString", "string"
		return true
	}
	ref := *def.ReferenceName

	if c, ok := w.schema.Constants[ref]; ok {
		p.TFType, p.GoType = "TypeString", "string"
		p.IsEnum = true
		p.EnumType = ref
		p.EnumValues = sortedConstantValues(c)
		return true
	}

	block := w.resolveModelBlock(ref)
	p.TFType = "TypeList"
	p.MaxItems = 1
	p.IsBlock = true
	p.BlockName = block.Name
	p.GoType = "[]" + block.Name
	return true
}

// applyList handles Type == "List", mapping the nested item type.
func (w *schemaWalker) applyList(p *Property, def pandora.ObjectDefinition) bool {
	p.TFType = "TypeList"
	if def.NestedItem == nil {
		p.GoType = "[]string"
		return true
	}
	item := *def.NestedItem
	if item.Type == "Reference" && item.ReferenceName != nil {
		if _, ok := w.schema.Constants[*item.ReferenceName]; ok {
			// list of enum values -> []string
			p.GoType = "[]string"
			p.IsEnum = true
			p.EnumType = *item.ReferenceName
			p.EnumValues = sortedConstantValues(w.schema.Constants[*item.ReferenceName])
			return true
		}
		block := w.resolveModelBlock(*item.ReferenceName)
		p.IsBlock = true
		p.BlockName = block.Name
		p.GoType = "[]" + block.Name
		return true
	}
	switch item.Type {
	case "Integer":
		p.GoType = "[]int64"
	case "Boolean":
		p.GoType = "[]bool"
	default:
		p.GoType = "[]string"
	}
	return true
}

// resolveModelBlock ensures a nested block model exists (recursing into its
// fields), deduplicating by Go struct name and guarding against cycles.
func (w *schemaWalker) resolveModelBlock(sdkModel string) *BlockModel {
	name := goIdent(sdkModel)
	if existing, ok := w.blocks[name]; ok {
		return existing
	}
	if w.visiting[name] {
		// Cycle - return a stub; fields already being resolved higher up.
		return &BlockModel{Name: name, SDKModel: sdkModel}
	}
	w.visiting[name] = true

	block := &BlockModel{Name: name, SDKModel: sdkModel}
	if model, ok := w.schema.Models[sdkModel]; ok {
		for _, key := range sortedFieldKeys(model.Fields) {
			// Only SystemData is stripped from nested models; unlike the
			// envelope, a nested "Id" is usually a real, meaningful field.
			if key == "SystemData" {
				continue
			}
			if prop := w.resolveField(key, model.Fields[key]); prop != nil {
				block.Properties = append(block.Properties, prop)
			}
		}
	}

	w.blocks[name] = block
	delete(w.visiting, name)
	return block
}

// referencedModel returns the model pointed at by a Reference ObjectDefinition.
func (w *schemaWalker) referencedModel(def pandora.ObjectDefinition) (pandora.Model, error) {
	if def.Type != "Reference" || def.ReferenceName == nil {
		return pandora.Model{}, fmt.Errorf("expected a model reference, got type %q", def.Type)
	}
	model, ok := w.schema.Models[*def.ReferenceName]
	if !ok {
		return pandora.Model{}, fmt.Errorf("referenced model %q not present in schema", *def.ReferenceName)
	}
	return model, nil
}

// sortedFieldKeys returns model field keys sorted by their snake_case Terraform
// name for deterministic output.
func sortedFieldKeys(fields map[string]pandora.Field) []string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return snake(keys[i]) < snake(keys[j])
	})
	return keys
}

// sortedConstantValues returns the allowed values of a constant, sorted.
func sortedConstantValues(c pandora.Constant) []string {
	vals := make([]string, 0, len(c.Values))
	for _, v := range c.Values {
		vals = append(vals, v)
	}
	sort.Strings(vals)
	return vals
}
