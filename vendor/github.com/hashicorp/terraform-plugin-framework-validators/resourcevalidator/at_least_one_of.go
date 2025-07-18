// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// AtLeastOneOf checks that a set of path.Expression has at least one non-null
// or unknown value.
func AtLeastOneOf(expressions ...path.Expression) resource.ConfigValidator {
	return &configvalidator.AtLeastOneOfValidator{
		PathExpressions: expressions,
	}
}
