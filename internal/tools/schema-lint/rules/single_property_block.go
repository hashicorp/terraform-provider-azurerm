// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

var (
	_ PropertyRule = singlePropertyBlock{}
	_ Fixer        = singlePropertyBlock{}
)

// singlePropertyBlock flags a MaxItems:1 block (a TypeList/TypeSet with a
// resource Elem) that contains exactly one nested property. Such blocks can
// usually be flattened into a single top-level property, per the schema design
// considerations ("Flattening nested properties" and the "Enabled" toggle
// pattern).
type singlePropertyBlock struct{}

func (singlePropertyBlock) ID() string   { return "SL006" }
func (singlePropertyBlock) Name() string { return "single-property-block" }

func (singlePropertyBlock) Description() string {
	return "A MaxItems:1 block with a single nested property should usually be flattened into one top-level property."
}

func (singlePropertyBlock) DefaultSeverity() Severity { return SeverityWarning }

func (singlePropertyBlock) FixHint() string { return "flatten the block" }

func (r singlePropertyBlock) CheckProperty(ctx PropertyContext) []Finding {
	if ctx.Schema.MaxItems != 1 {
		return nil
	}

	block, ok := BlockElem(ctx.Schema)
	if !ok || len(block.Schema) != 1 {
		return nil
	}

	var childName string
	var child providerschema.SchemaJSON
	for k, v := range block.Schema {
		childName, child = k, v
	}

	message := fmt.Sprintf("block %q has a single nested property %q (MaxItems 1); consider flattening it", ctx.Path, childName)

	suggestion := fmt.Sprintf("flatten %q into a single top-level %q property", ctx.Path, ctx.Name+"_"+childName)
	if childName == "enabled" && child.Type == TypeBool {
		suggestion = fmt.Sprintf("replace the block %q with a single top-level boolean %q", ctx.Path, ctx.Name+"_enabled")
	}

	return []Finding{propertyFindingFix(r, ctx, message, suggestion)}
}
