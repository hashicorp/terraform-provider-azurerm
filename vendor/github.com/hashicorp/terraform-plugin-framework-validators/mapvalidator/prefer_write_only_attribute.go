// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
)

// PreferWriteOnlyAttribute returns a warning if the Terraform client supports
// write-only attributes, and the attribute that the validator is applied to has a value.
// It takes in a path.Expression that represents the write-only attribute schema location,
// and the warning message will indicate that the write-only attribute should be preferred.
//
// This validator should only be used for resource attributes as other schema types do not
// support write-only attributes.
//
// This implements the validation logic declaratively within the schema.
// Refer to [resourcevalidator.PreferWriteOnlyAttribute]
// for declaring this type of validation outside the schema definition.
func PreferWriteOnlyAttribute(writeOnlyAttribute path.Expression) validator.Map {
	return schemavalidator.PreferWriteOnlyAttribute{
		WriteOnlyAttribute: writeOnlyAttribute,
	}
}
