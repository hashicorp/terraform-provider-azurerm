// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gen

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

// blockUsage reports, for each block name, whether it is referenced as a single
// nested object (MaxItems 1) and/or as a list, so the correct expand/flatten
// helpers are generated.
func blockUsage(res *ir.ResourceIR) (single, list map[string]bool) {
	single = map[string]bool{}
	list = map[string]bool{}
	mark := func(p *ir.Property) {
		if !p.IsBlock {
			return
		}
		if p.MaxItems == 1 {
			single[p.BlockName] = true
		} else {
			list[p.BlockName] = true
		}
	}
	for _, p := range res.TopLevel {
		mark(p)
	}
	for _, b := range res.Blocks {
		for _, p := range b.Properties {
			mark(p)
		}
	}
	return single, list
}

// expandFieldExpr returns the RHS Go expression assigning a model field into an
// SDK model field (model -> SDK).
func expandFieldExpr(res *ir.ResourceIR, p *ir.Property, srcVar string) string {
	src := fmt.Sprintf("%s.%s", srcVar, p.GoField)
	switch {
	case p.IsBlock && p.MaxItems == 1:
		return fmt.Sprintf("expand%s(%s)", p.BlockName, src)
	case p.IsBlock:
		return fmt.Sprintf("expand%sList(%s)", p.BlockName, src)
	case p.IsEnum && p.TFType == "TypeString":
		return fmt.Sprintf("pointer.To(%s.%s(%s))", res.SDKPackage, p.EnumType, src)
	default:
		return fmt.Sprintf("pointer.To(%s)", src)
	}
}

// flattenFieldExpr returns the RHS Go expression assigning an SDK model field
// into a model field (SDK -> model).
func flattenFieldExpr(res *ir.ResourceIR, p *ir.Property, srcVar string) string {
	src := fmt.Sprintf("%s.%s", srcVar, p.SDKField)
	switch {
	case p.IsBlock && p.MaxItems == 1:
		return fmt.Sprintf("flatten%s(%s)", p.BlockName, src)
	case p.IsBlock:
		return fmt.Sprintf("flatten%sList(%s)", p.BlockName, src)
	case p.IsEnum && p.TFType == "TypeString":
		return fmt.Sprintf("string(pointer.From(%s))", src)
	default:
		return fmt.Sprintf("pointer.From(%s)", src)
	}
}

// RenderExpandFlatten renders the expand/flatten helper functions for every
// nested block model, generating single and/or list variants as used.
func RenderExpandFlatten(res *ir.ResourceIR) string {
	single, list := blockUsage(res)
	var sb strings.Builder
	for _, b := range res.Blocks {
		if single[b.Name] {
			renderExpandSingle(&sb, res, b)
			renderFlattenSingle(&sb, res, b)
		}
		if list[b.Name] {
			renderExpandList(&sb, res, b)
			renderFlattenList(&sb, res, b)
		}
	}
	return sb.String()
}

func renderExpandSingle(sb *strings.Builder, res *ir.ResourceIR, b *ir.BlockModel) {
	fmt.Fprintf(sb, "func expand%s(input []%s) *%s.%s {\n", b.Name, b.Name, res.SDKPackage, b.SDKModel)
	sb.WriteString("if len(input) == 0 {\nreturn nil\n}\n")
	sb.WriteString("v := input[0]\n")
	fmt.Fprintf(sb, "return &%s.%s{\n", res.SDKPackage, b.SDKModel)
	for _, p := range b.Properties {
		fmt.Fprintf(sb, "%s: %s,\n", p.SDKField, expandFieldExpr(res, p, "v"))
	}
	sb.WriteString("}\n}\n\n")
}

func renderFlattenSingle(sb *strings.Builder, res *ir.ResourceIR, b *ir.BlockModel) {
	fmt.Fprintf(sb, "func flatten%s(input *%s.%s) []%s {\n", b.Name, res.SDKPackage, b.SDKModel, b.Name)
	fmt.Fprintf(sb, "if input == nil {\nreturn []%s{}\n}\n", b.Name)
	fmt.Fprintf(sb, "return []%s{\n{\n", b.Name)
	for _, p := range b.Properties {
		fmt.Fprintf(sb, "%s: %s,\n", p.GoField, flattenFieldExpr(res, p, "input"))
	}
	sb.WriteString("},\n}\n}\n\n")
}

func renderExpandList(sb *strings.Builder, res *ir.ResourceIR, b *ir.BlockModel) {
	fmt.Fprintf(sb, "func expand%sList(input []%s) *[]%s.%s {\n", b.Name, b.Name, res.SDKPackage, b.SDKModel)
	sb.WriteString("if len(input) == 0 {\nreturn nil\n}\n")
	fmt.Fprintf(sb, "result := make([]%s.%s, 0, len(input))\n", res.SDKPackage, b.SDKModel)
	sb.WriteString("for _, v := range input {\n")
	fmt.Fprintf(sb, "result = append(result, %s.%s{\n", res.SDKPackage, b.SDKModel)
	for _, p := range b.Properties {
		fmt.Fprintf(sb, "%s: %s,\n", p.SDKField, expandFieldExpr(res, p, "v"))
	}
	sb.WriteString("})\n}\n")
	sb.WriteString("return &result\n}\n\n")
}

func renderFlattenList(sb *strings.Builder, res *ir.ResourceIR, b *ir.BlockModel) {
	fmt.Fprintf(sb, "func flatten%sList(input *[]%s.%s) []%s {\n", b.Name, res.SDKPackage, b.SDKModel, b.Name)
	fmt.Fprintf(sb, "if input == nil {\nreturn []%s{}\n}\n", b.Name)
	fmt.Fprintf(sb, "result := make([]%s, 0, len(*input))\n", b.Name)
	sb.WriteString("for _, v := range *input {\n")
	fmt.Fprintf(sb, "result = append(result, %s{\n", b.Name)
	for _, p := range b.Properties {
		fmt.Fprintf(sb, "%s: %s,\n", p.GoField, flattenFieldExpr(res, p, "v"))
	}
	sb.WriteString("})\n}\n")
	sb.WriteString("return result\n}\n\n")
}

// RenderCreatePayload renders the SDK request payload literal built from the
// plan model (variable "config"), including the ARM Properties envelope.
func RenderCreatePayload(res *ir.ResourceIR) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "payload := %s.%s{\n", res.SDKPackage, res.CreateModel)
	for _, p := range res.TopLevel {
		if p.UnderProperties {
			continue
		}
		switch p.TFName {
		case "name":
			sb.WriteString("Name: pointer.To(config.Name),\n")
		case "location":
			sb.WriteString("Location: location.Normalize(config.Location),\n")
		case "tags":
			sb.WriteString("Tags: pointer.To(config.Tags),\n")
		}
	}
	if res.PropertiesModel != "" {
		fmt.Fprintf(&sb, "Properties: &%s.%s{\n", res.SDKPackage, res.PropertiesModel)
		for _, p := range res.TopLevel {
			if !p.UnderProperties || p.Computed {
				continue
			}
			fmt.Fprintf(&sb, "%s: %s,\n", p.SDKField, expandFieldExpr(res, p, "config"))
		}
		sb.WriteString("},\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}

// RenderFlattenMethod renders the resource's flatten method, which maps a
// (non-nil) SDK model + parsed ID onto the model struct and encodes it. It is
// shared by Read and by the list resource.
func RenderFlattenMethod(res *ir.ResourceIR) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "func (r %sResource) flatten(metadata sdk.ResourceMetaData, id *%s.%s, model *%s.%s) error {\n",
		res.Name, res.SDKPackage, res.IDTypeName, res.SDKPackage, res.ReadModel)
	fmt.Fprintf(&sb, "state := %s{\n", res.ModelStructName)
	if res.IDNameSegment != "" {
		fmt.Fprintf(&sb, "Name: id.%s,\n", res.IDNameSegment)
	}
	if res.HasResourceGroup {
		sb.WriteString("ResourceGroup: id.ResourceGroupName,\n")
	}
	sb.WriteString("}\n\n")

	for _, p := range res.TopLevel {
		if p.UnderProperties {
			continue
		}
		switch p.TFName {
		case "location":
			sb.WriteString("state.Location = location.Normalize(model.Location)\n")
		case "tags":
			sb.WriteString("state.Tags = pointer.From(model.Tags)\n")
		}
	}

	if res.PropertiesModel != "" {
		sb.WriteString("if props := model.Properties; props != nil {\n")
		for _, p := range res.TopLevel {
			if !p.UnderProperties {
				continue
			}
			fmt.Fprintf(&sb, "state.%s = %s\n", p.GoField, flattenFieldExpr(res, p, "props"))
		}
		sb.WriteString("}\n")
	}

	sb.WriteString("\nreturn metadata.Encode(&state)\n}\n")
	return sb.String()
}

// RenderUpdateBody renders the full body of the Update func, building the PATCH
// payload from properties that have changed. Only properties present in the
// update model's surface are emitted, keeping the output compile-clean.
func RenderUpdateBody(res *ir.ResourceIR) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "client := metadata.Client.%s.%s\n\n", res.ServiceName, res.ClientField)
	fmt.Fprintf(&sb, "id, err := %s.%s(metadata.ResourceData.Id())\n", res.SDKPackage, res.IDParseFunc)
	sb.WriteString("if err != nil {\nreturn err\n}\n\n")
	fmt.Fprintf(&sb, "var state %s\n", res.ModelStructName)
	sb.WriteString("if err := metadata.Decode(&state); err != nil {\nreturn fmt.Errorf(\"decoding: %+v\", err)\n}\n\n")
	fmt.Fprintf(&sb, "payload := %s.%s{}\n\n", res.SDKPackage, res.UpdateModel)

	stateUsed := false
	if res.UpdateHasTags {
		sb.WriteString("if metadata.ResourceData.HasChange(\"tags\") {\npayload.Tags = pointer.To(state.Tags)\n}\n\n")
		stateUsed = true
	}

	for _, p := range res.TopLevel {
		if !p.UnderProperties || p.Computed || p.ForceNew {
			continue
		}
		if !res.UpdatablePropSDKFields[p.SDKField] {
			continue
		}
		fmt.Fprintf(&sb, "if metadata.ResourceData.HasChange(%q) {\n", p.TFName)
		fmt.Fprintf(&sb, "payload.Properties.%s = %s\n", p.SDKField, expandFieldExpr(res, p, "state"))
		sb.WriteString("}\n\n")
		stateUsed = true
	}

	if !stateUsed {
		sb.WriteString("_ = state\n\n")
	}

	if res.UpdateLRO {
		fmt.Fprintf(&sb, "if err := client.%sThenPoll(ctx, *id, payload); err != nil {\nreturn fmt.Errorf(\"updating %%s: %%+v\", id, err)\n}\n\n", res.UpdateOp)
	} else {
		fmt.Fprintf(&sb, "if _, err := client.%s(ctx, *id, payload); err != nil {\nreturn fmt.Errorf(\"updating %%s: %%+v\", id, err)\n}\n\n", res.UpdateOp)
	}
	sb.WriteString("return nil\n")
	return sb.String()
}
