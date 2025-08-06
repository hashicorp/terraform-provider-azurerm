// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ExactlyOneOf checks that of a set of path.Expression,
// including the attribute or block the validator is applied to,
// one and only one attribute has a value.
// It will also cause a validation error if none are specified.
//
// This implements the validation logic declaratively within the schema.
// Refer to [datasourcevalidator.ExactlyOneOf],
// [providervalidator.ExactlyOneOf], or [resourcevalidator.ExactlyOneOf]
// for declaring this type of validation outside the schema definition.
//
// Relative path.Expression will be resolved using the attribute or block
// being validated.
func ExactlyOneOf(expressions ...path.Expression) validator.List {
	return schemavalidator.ExactlyOneOfValidator{
		PathExpressions: expressions,
	}
}
