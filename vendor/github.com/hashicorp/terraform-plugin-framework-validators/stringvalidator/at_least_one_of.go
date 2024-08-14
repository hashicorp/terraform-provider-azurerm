// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AtLeastOneOf checks that of a set of path.Expression,
// including the attribute this validator is applied to,
// at least one has a non-null value.
//
// This implements the validation logic declaratively within the tfsdk.Schema.
// Refer to [datasourcevalidator.AtLeastOneOf],
// [providervalidator.AtLeastOneOf], or [resourcevalidator.AtLeastOneOf]
// for declaring this type of validation outside the schema definition.
//
// Any relative path.Expression will be resolved using the attribute being
// validated.
func AtLeastOneOf(expressions ...path.Expression) validator.String {
	return schemavalidator.AtLeastOneOfValidator{
		PathExpressions: expressions,
	}
}
