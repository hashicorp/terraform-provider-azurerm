// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Conflicting checks that a set of path.Expression, are not configured
// simultaneously.
func Conflicting(expressions ...path.Expression) resource.ConfigValidator {
	return &configvalidator.ConflictingValidator{
		PathExpressions: expressions,
	}
}
