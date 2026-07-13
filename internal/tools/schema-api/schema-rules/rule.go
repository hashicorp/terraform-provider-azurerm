// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import "github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"

type BreakingChangeRule interface {
	Check(base providerschema.SchemaJSON, current providerschema.SchemaJSON, propertyName string) *string
}

var BreakingChangeRules = []BreakingChangeRule{
	becomeComputedOnly{},
	newRequiredPropertyExistingResource{},
	optionalRemoveComputed{},
	optionalToRequired{},
	propertyType{},
}

var BreakingChangeRulesDataSource = []BreakingChangeRule{
	propertyType{},
}
