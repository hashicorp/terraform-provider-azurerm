// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func attributePlanModificationTypableError(schemaPath path.Path, value attr.Value) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		schemaPath,
		"Unexpected Attribute Plan Modifier Type Conversion Error",
		"An unexpected issue occurred while trying to get the correct type during attribute plan modification. "+
			"Expected the Valuable implementation Type() method to return a Typable. "+
			"This is likely an implementation error in terraform-plugin-framework and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("Value Type: %T\n", value)+
			fmt.Sprintf("Path: %s", schemaPath),
	)
}

func schemaDataValueError(ctx context.Context, value attr.Value, description fwschemadata.DataDescription, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		description.Title()+" Value Error",
		"An unexpected error occurred while fetching a "+value.Type(ctx).String()+" element value in the "+description.String()+". "+
			"This is an issue with the provider and should be reported to the provider developers.\n\n"+
			"Original Error: "+err.Error(),
	)
}

func schemaDataWalkError(schemaPath path.Path, value attr.Value) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		schemaPath,
		"Schema Data Walk Error",
		"An unexpected error occurred while walking the schema for data modification. "+
			"This is an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("unknown attribute value type (%T) at path: %s", value, schemaPath),
	)
}
