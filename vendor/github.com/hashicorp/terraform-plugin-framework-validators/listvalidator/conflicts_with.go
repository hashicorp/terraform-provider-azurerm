// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ConflictsWith checks that a set of path.Expression,
// including the attribute or block the validator is applied to,
// do not have a value simultaneously.
//
// This implements the validation logic declaratively within the schema.
// Refer to [datasourcevalidator.Conflicting],
// [providervalidator.Conflicting], or [resourcevalidator.Conflicting]
// for declaring this type of validation outside the schema definition.
//
// Relative path.Expression will be resolved using the attribute or block
// being validated.
func ConflictsWith(expressions ...path.Expression) validator.List {
	return schemavalidator.ConflictsWithValidator{
		PathExpressions: expressions,
	}
}
