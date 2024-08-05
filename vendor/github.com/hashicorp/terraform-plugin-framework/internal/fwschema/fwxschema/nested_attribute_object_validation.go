// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// NestedAttributeObjectWithValidators is an optional interface on
// NestedAttributeObject which enables Object validation support.
type NestedAttributeObjectWithValidators interface {
	fwschema.NestedAttributeObject

	// ObjectValidators should return a list of Object validators.
	ObjectValidators() []validator.Object
}
