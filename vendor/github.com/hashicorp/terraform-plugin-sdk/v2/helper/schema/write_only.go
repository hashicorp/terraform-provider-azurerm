// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/internal/configs/configschema"
)

// setWriteOnlyNullValues takes a cty.Value, and compares it to the schema setting any non-null
// values that are writeOnly to null.
func setWriteOnlyNullValues(val cty.Value, schema *configschema.Block) cty.Value {
	if !val.IsKnown() || val.IsNull() {
		return val
	}

	valMap := val.AsValueMap()
	newVals := make(map[string]cty.Value)

	for name, attr := range schema.Attributes {
		v := valMap[name]

		if attr.WriteOnly && !v.IsNull() {
			newVals[name] = cty.NullVal(attr.Type)
			continue
		}

		newVals[name] = v
	}

	for name, blockS := range schema.BlockTypes {
		blockVal := valMap[name]
		if blockVal.IsNull() || !blockVal.IsKnown() {
			newVals[name] = blockVal
			continue
		}

		blockValType := blockVal.Type()
		blockElementType := blockS.Block.ImpliedType()

		// This switches on the value type here, so we can correctly switch
		// between Tuples/Lists and Maps/Objects.
		switch {
		case blockS.Nesting == configschema.NestingSingle || blockS.Nesting == configschema.NestingGroup:
			// NestingSingle is the only exception here, where we treat the
			// block directly as an object
			newVals[name] = setWriteOnlyNullValues(blockVal, &blockS.Block)

		case blockValType.IsSetType(), blockValType.IsListType(), blockValType.IsTupleType():
			listVals := blockVal.AsValueSlice()
			newListVals := make([]cty.Value, 0, len(listVals))

			for _, v := range listVals {
				newListVals = append(newListVals, setWriteOnlyNullValues(v, &blockS.Block))
			}

			switch {
			case blockValType.IsSetType():
				switch len(newListVals) {
				case 0:
					newVals[name] = cty.SetValEmpty(blockElementType)
				default:
					newVals[name] = cty.SetVal(newListVals)
				}
			case blockValType.IsListType():
				switch len(newListVals) {
				case 0:
					newVals[name] = cty.ListValEmpty(blockElementType)
				default:
					newVals[name] = cty.ListVal(newListVals)
				}
			case blockValType.IsTupleType():
				newVals[name] = cty.TupleVal(newListVals)
			}

		case blockValType.IsMapType(), blockValType.IsObjectType():
			mapVals := blockVal.AsValueMap()
			newMapVals := make(map[string]cty.Value)

			for k, v := range mapVals {
				newMapVals[k] = setWriteOnlyNullValues(v, &blockS.Block)
			}

			switch {
			case blockValType.IsMapType():
				switch len(newMapVals) {
				case 0:
					newVals[name] = cty.MapValEmpty(blockElementType)
				default:
					newVals[name] = cty.MapVal(newMapVals)
				}
			case blockValType.IsObjectType():
				if len(newMapVals) == 0 {
					// We need to populate empty values to make a valid object.
					for attr, ty := range blockElementType.AttributeTypes() {
						newMapVals[attr] = cty.NullVal(ty)
					}
				}
				newVals[name] = cty.ObjectVal(newMapVals)
			}

		default:
			panic(fmt.Sprintf("failed to set null values for nested block %q:%#v", name, blockValType))
		}
	}

	return cty.ObjectVal(newVals)
}

// validateWriteOnlyNullValues validates that write-only attribute values
// are null to ensure that write-only values are not sent to unsupported
// Terraform client versions.
//
// it takes a cty.Value, and compares it to the schema and throws an
// error diagnostic for each non-null writeOnly attribute value.
func validateWriteOnlyNullValues(val cty.Value, schema *configschema.Block, path cty.Path) diag.Diagnostics {
	if !val.IsKnown() || val.IsNull() {
		return diag.Diagnostics{}
	}

	valMap := val.AsValueMap()
	diags := make([]diag.Diagnostic, 0)

	var attrNames []string
	for k := range schema.Attributes {
		attrNames = append(attrNames, k)
	}

	// Sort the attribute names to produce diags in a consistent order.
	sort.Strings(attrNames)

	for _, name := range attrNames {
		attr := schema.Attributes[name]
		v := valMap[name]

		if attr.WriteOnly && !v.IsNull() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Write-only Attribute Not Allowed",
				Detail: fmt.Sprintf("The resource contains a non-null value for write-only attribute %q ", name) +
					"Write-only attributes are only supported in Terraform 1.11 and later.",
				AttributePath: append(path, cty.GetAttrStep{Name: name}),
			})
		}
	}

	var blockNames []string
	for k := range schema.BlockTypes {
		blockNames = append(blockNames, k)
	}

	// Sort the block names to produce diags in a consistent order.
	sort.Strings(blockNames)

	for _, name := range blockNames {
		blockS := schema.BlockTypes[name]
		blockVal := valMap[name]
		if blockVal.IsNull() || !blockVal.IsKnown() {
			continue
		}

		blockValType := blockVal.Type()
		blockPath := append(path, cty.GetAttrStep{Name: name})

		// This switches on the value type here, so we can correctly switch
		// between Tuples/Lists and Maps/Objects.
		switch {
		case blockS.Nesting == configschema.NestingSingle || blockS.Nesting == configschema.NestingGroup:
			// NestingSingle is the only exception here, where we treat the
			// block directly as an object
			diags = append(diags, validateWriteOnlyNullValues(blockVal, &blockS.Block, blockPath)...)
		case blockValType.IsSetType():
			setVals := blockVal.AsValueSlice()

			for _, v := range setVals {
				setBlockPath := append(blockPath, cty.IndexStep{
					Key: v,
				})
				diags = append(diags, validateWriteOnlyNullValues(v, &blockS.Block, setBlockPath)...)
			}

		case blockValType.IsListType(), blockValType.IsTupleType():
			listVals := blockVal.AsValueSlice()

			for i, v := range listVals {
				listBlockPath := append(blockPath, cty.IndexStep{
					Key: cty.NumberIntVal(int64(i)),
				})
				diags = append(diags, validateWriteOnlyNullValues(v, &blockS.Block, listBlockPath)...)
			}

		case blockValType.IsMapType(), blockValType.IsObjectType():
			mapVals := blockVal.AsValueMap()

			for k, v := range mapVals {
				mapBlockPath := append(blockPath, cty.IndexStep{
					Key: cty.StringVal(k),
				})
				diags = append(diags, validateWriteOnlyNullValues(v, &blockS.Block, mapBlockPath)...)
			}

		default:
			panic(fmt.Sprintf("failed to validate WriteOnly values for nested block %q:%#v", name, blockValType))
		}
	}

	return diags
}
