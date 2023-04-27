package schema_rules

import "github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"

type BreakingChangeRule interface {
	Check(base providerjson.SchemaJSON, current providerjson.SchemaJSON, propertyName string) *string
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
