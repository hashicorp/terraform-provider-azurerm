// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwtype

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ContainsCollectionWithDynamic will return true if an attr.Type is a complex type that either is or contains any
// collection types with dynamic types, which are not supported by the framework type system. Primitives, invalid
// types (missingType), or nil will return false.
//
// Unsupported collection types include:
//   - Lists that contain a dynamic type
//   - Maps that contain a dynamic type
//   - Sets that contain a dynamic type
func ContainsCollectionWithDynamic(typ attr.Type) bool {
	switch attrType := typ.(type) {
	// We haven't run into a collection type yet, so it's valid for this to be a dynamic type
	case basetypes.DynamicTypable:
		return false
	// Lists, maps, sets
	case attr.TypeWithElementType:
		// We found a collection, need to ensure there are no dynamics from this point on.
		return containsDynamic(attrType.ElementType())
	// Tuples
	case attr.TypeWithElementTypes:
		for _, elemType := range attrType.ElementTypes() {
			hasDynamic := ContainsCollectionWithDynamic(elemType)
			if hasDynamic {
				return true
			}
		}
		return false
	// Objects
	case attr.TypeWithAttributeTypes:
		for _, objAttrType := range attrType.AttributeTypes() {
			hasDynamic := ContainsCollectionWithDynamic(objAttrType)
			if hasDynamic {
				return true
			}
		}
		return false
	// Primitives, missing types, etc.
	default:
		return false
	}
}

// containsDynamic will return true if `typ` is a dynamic type or has any nested types that contain a dynamic type.
func containsDynamic(typ attr.Type) bool {
	switch attrType := typ.(type) {
	// Found a dynamic!
	case basetypes.DynamicTypable:
		return true
	// Lists, maps, sets
	case attr.TypeWithElementType:
		return containsDynamic(attrType.ElementType())
	// Tuples
	case attr.TypeWithElementTypes:
		for _, elemType := range attrType.ElementTypes() {
			hasDynamic := containsDynamic(elemType)
			if hasDynamic {
				return true
			}
		}
		return false
	// Objects
	case attr.TypeWithAttributeTypes:
		for _, objAttrType := range attrType.AttributeTypes() {
			hasDynamic := containsDynamic(objAttrType)
			if hasDynamic {
				return true
			}
		}
		return false
	// Primitives, missing types, etc.
	default:
		return false
	}
}

func AttributeCollectionWithDynamicTypeDiag(attributePath path.Path) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Invalid Schema Implementation",
		"When validating the schema, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("%q is an attribute that contains a collection type with a nested dynamic type.\n\n", attributePath)+
			"Dynamic types inside of collections are not currently supported in terraform-plugin-framework. "+
			fmt.Sprintf("If underlying dynamic values are required, replace the %q attribute definition with DynamicAttribute instead.", attributePath),
	)
}

func BlockCollectionWithDynamicTypeDiag(attributePath path.Path) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Invalid Schema Implementation",
		"When validating the schema, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("%q is a block that contains a collection type with a nested dynamic type.\n\n", attributePath)+
			"Dynamic types inside of collections are not currently supported in terraform-plugin-framework. "+
			fmt.Sprintf("If underlying dynamic values are required, replace the %q block definition with a DynamicAttribute.", attributePath),
	)
}

func ParameterCollectionWithDynamicTypeDiag(argument int64, name string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Invalid Function Definition",
		"When validating the function definition, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("Parameter %q at position %d contains a collection type with a nested dynamic type.\n\n", name, argument)+
			"Dynamic types inside of collections are not currently supported in terraform-plugin-framework. "+
			fmt.Sprintf("If underlying dynamic values are required, replace the %q parameter definition with DynamicParameter instead.", name),
	)
}

func VariadicParameterCollectionWithDynamicTypeDiag(name string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Invalid Function Definition",
		"When validating the function definition, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("Variadic parameter %q contains a collection type with a nested dynamic type.\n\n", name)+
			"Dynamic types inside of collections are not currently supported in terraform-plugin-framework. "+
			"If underlying dynamic values are required, replace the variadic parameter definition with DynamicParameter instead.",
	)
}

func ReturnCollectionWithDynamicTypeDiag() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Invalid Function Definition",
		"When validating the function definition, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			"Return contains a collection type with a nested dynamic type.\n\n"+
			"Dynamic types inside of collections are not currently supported in terraform-plugin-framework. "+
			"If underlying dynamic values are required, replace the return definition with DynamicReturn instead.",
	)
}
