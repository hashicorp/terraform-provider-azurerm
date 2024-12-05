// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwtype

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ContainsMissingUnderlyingType will return true if an attr.Type is
// a complex type that either is or contains any collection types with missing
// element or attribute types. Primitives will return false. Nil will return
// true.
func ContainsMissingUnderlyingType(typ attr.Type) bool {
	// The below logic must use AttrTypes/ElemType/ElemTypes directly, or the
	// types package will return the unexported missingType type, which cannot
	// be caught here.
	switch attrType := typ.(type) {
	case nil:
		return true
	case basetypes.ListType:
		return ContainsMissingUnderlyingType(attrType.ElemType)
	case basetypes.MapType:
		return ContainsMissingUnderlyingType(attrType.ElemType)
	case basetypes.ObjectType:
		for _, objAttrType := range attrType.AttrTypes {
			if ContainsMissingUnderlyingType(objAttrType) {
				return true
			}
		}

		return false
	case basetypes.SetType:
		return ContainsMissingUnderlyingType(attrType.ElemType)
	case basetypes.TupleType:
		for _, elemType := range attrType.ElemTypes {
			if ContainsMissingUnderlyingType(elemType) {
				return true
			}
		}

		return false
	// Everything else (primitives, custom types, etc.)
	default:
		return false
	}
}

func ParameterMissingUnderlyingTypeDiag(name string, position *int64) diag.Diagnostic {
	if position == nil {
		return diag.NewErrorDiagnostic(
			"Invalid Function Definition",
			"When validating the function definition, an implementation issue was found. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				"Variadic parameter is missing underlying type.\n\n"+
				"Collection element and object attribute types are always required in Terraform.",
		)
	}

	return diag.NewErrorDiagnostic(
		"Invalid Function Definition",
		"When validating the function definition, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("Parameter %q at position %d is missing underlying type.\n\n", name, *position)+
			"Collection element and object attribute types are always required in Terraform.",
	)
}

func ReturnMissingUnderlyingTypeDiag() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Invalid Function Definition",
		"When validating the function definition, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			"Return is missing underlying type.\n\n"+
			"Collection element and object attribute types are always required in Terraform.",
	)
}
