// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

var _ BreakingChangeRule = optionalRemoveComputed{}

type optionalRemoveComputed struct {
}

// Check - Checks that Computed is not removed from Optional properties as user configs may not supply the value, but the state will contain one, causing a diff./
func (optionalRemoveComputed) Check(base providerjson.SchemaJSON, current providerjson.SchemaJSON, propertyName string) *string {
	if (base.Optional && base.Computed) && (current.Optional && !current.Computed) {
		return pointer.To(fmt.Sprintf("cannot remove Computed from the Optional property %q as the user config may not supply the value, thus causing a diff", propertyName))
	}

	return nil
}
