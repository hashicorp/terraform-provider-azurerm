// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

// Checks that the Schema Property Type has not changed

var _ BreakingChangeRule = propertyType{}

type propertyType struct{}

// Check - Checks for invalid type changes. At the time of writing the only allowed change is a Set to a List
func (propertyType) Check(base providerjson.SchemaJSON, current providerjson.SchemaJSON, propertyName string) *string {
	if (base.Type != "" && current.Type != "" && base.Type != providerjson.SchemaTypeSet) && base.Type != current.Type {
		return pointer.To(fmt.Sprintf("schema type has changed for %q (%+v to %+v)", propertyName, base.Type, current.Type))
	}

	return nil
}
