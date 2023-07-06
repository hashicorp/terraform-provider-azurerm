// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

// Checks that the Schema Property Type has not changed

var _ BreakingChangeRule = newRequiredPropertyExistingResource{}

type newRequiredPropertyExistingResource struct{}

// Check - Checks that a newly introduced property is not marked as Required since this will not be in users configurations.
func (newRequiredPropertyExistingResource) Check(base providerjson.SchemaJSON, current providerjson.SchemaJSON, propertyName string) *string {
	if base.Type == "" && current.Required {
		return pointer.To(fmt.Sprintf("new property %q is Required", propertyName))
	}

	return nil
}
