// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// NestedBlockObjectWithPlanModifiers is an optional interface on
// NestedBlockObject which enables Object plan modification support.
type NestedBlockObjectWithPlanModifiers interface {
	fwschema.NestedBlockObject

	// ObjectPlanModifiers should return a list of Object plan modifiers.
	ObjectPlanModifiers() []planmodifier.Object
}
