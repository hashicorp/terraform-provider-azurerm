// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

type defaultValueChange struct {
}

var _ BreakingChangeRule = defaultValueChange{}

// Check - Checks that an Optional or Required property is not updated to become Computed only
func (o defaultValueChange) Check(base providerjson.SchemaJSON, current providerjson.SchemaJSON, propertyName string) *string {
	if base.Default != current.Default {
		return pointer.To(fmt.Sprintf("Cannot change property %q to Computed only", propertyName))
	}

	return nil
}
