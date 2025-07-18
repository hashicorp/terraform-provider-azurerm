// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// RequiredTogether checks that a set of path.Expression either has all known
// or all null values.
func RequiredTogether(expressions ...path.Expression) resource.ConfigValidator {
	return &configvalidator.RequiredTogetherValidator{
		PathExpressions: expressions,
	}
}
