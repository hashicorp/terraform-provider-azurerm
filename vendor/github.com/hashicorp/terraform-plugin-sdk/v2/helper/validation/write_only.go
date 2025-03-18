// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PreferWriteOnlyAttribute is a ValidateRawResourceConfigFunc that returns a warning
// if the Terraform client supports write-only attributes and the old attribute is
// not null.
// The last step in the path must be a cty.GetAttrStep{}.
// When creating a cty.IndexStep{} to into a nested attribute, use an unknown value
// of the index type to indicate any key value.
// For lists: cty.Index(cty.UnknownVal(cty.Number)),
// For maps: cty.Index(cty.UnknownVal(cty.String)),
// For sets: cty.Index(cty.UnknownVal(cty.Object(nil))),
func PreferWriteOnlyAttribute(oldAttribute cty.Path, writeOnlyAttribute cty.Path) schema.ValidateRawResourceConfigFunc {
	return func(ctx context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
		if !req.WriteOnlyAttributesAllowed {
			return
		}

		pathLen := len(writeOnlyAttribute)

		if pathLen == 0 {
			return
		}

		lastStep := writeOnlyAttribute[pathLen-1]

		// Only attribute steps have a Name field
		writeOnlyAttrStep, ok := lastStep.(cty.GetAttrStep)
		if !ok {
			resp.Diagnostics = diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Invalid writeOnlyAttribute path",
					Detail: "The Terraform Provider unexpectedly provided a path that does not match the current schema. " +
						"This can happen if the path does not correctly follow the schema in structure or types. " +
						"Please report this to the provider developers. \n\n" +
						"The writeOnlyAttribute path provided is invalid. The last step in the path must be a cty.GetAttrStep{}",
					AttributePath: writeOnlyAttribute,
				},
			}
			return
		}

		var oldAttrs []attribute

		err := cty.Walk(req.RawConfig, func(path cty.Path, value cty.Value) (bool, error) {
			if PathMatches(path, oldAttribute) {
				oldAttrs = append(oldAttrs, attribute{
					value: value,
					path:  path,
				})
			}

			return true, nil
		})
		if err != nil {
			return
		}

		for _, attr := range oldAttrs {
			attrPath := attr.path.Copy()

			pathLen = len(attrPath)

			if pathLen == 0 {
				return
			}

			lastStep = attrPath[pathLen-1]

			// Only attribute steps have a Name field
			attrStep, ok := lastStep.(cty.GetAttrStep)
			if !ok {
				resp.Diagnostics = diag.Diagnostics{
					{
						Severity: diag.Error,
						Summary:  "Invalid oldAttribute path",
						Detail: "The Terraform Provider unexpectedly provided a path that does not match the current schema. " +
							"This can happen if the path does not correctly follow the schema in structure or types. " +
							"Please report this to the provider developers. \n\n" +
							"The oldAttribute path provided is invalid. The last step in the path must be a cty.GetAttrStep{}",
						AttributePath: attrPath,
					},
				}
				return
			}

			if !attr.value.IsNull() {
				resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Available Write-only Attribute Alternative",
					Detail: fmt.Sprintf("The attribute %s has a write-only alternative %s available. "+
						"Use the write-only alternative of the attribute when possible.", attrStep.Name, writeOnlyAttrStep.Name),
					AttributePath: attr.path,
				})
			}
		}
	}
}

type attribute struct {
	value cty.Value
	path  cty.Path
}
