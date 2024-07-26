// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// BlockWithListValidators is an optional interface on Block which
// enables List validation support.
type BlockWithListValidators interface {
	fwschema.Block

	// ListValidators should return a list of List validators.
	ListValidators() []validator.List
}

// BlockWithObjectValidators is an optional interface on Block which
// enables Object validation support.
type BlockWithObjectValidators interface {
	fwschema.Block

	// ObjectValidators should return a list of Object validators.
	ObjectValidators() []validator.Object
}

// BlockWithSetValidators is an optional interface on Block which
// enables Set validation support.
type BlockWithSetValidators interface {
	fwschema.Block

	// SetValidators should return a list of Set validators.
	SetValidators() []validator.Set
}
