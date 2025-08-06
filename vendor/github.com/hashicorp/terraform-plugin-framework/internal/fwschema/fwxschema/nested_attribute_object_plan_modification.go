// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// NestedAttributeObjectWithPlanModifiers is an optional interface on
// NestedAttributeObject which enables Object plan modification support.
type NestedAttributeObjectWithPlanModifiers interface {
	fwschema.NestedAttributeObject

	// ObjectPlanModifiers should return a list of Object plan modifiers.
	ObjectPlanModifiers() []planmodifier.Object
}
