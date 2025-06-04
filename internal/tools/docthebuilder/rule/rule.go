package rule

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/data"
)

type Rule interface {
	// ID returns the ID of the rule, e.g. `G001`.
	ID() string
	// Name returns a friendly name of the rule, e.g. `Validate API Section`.
	Name() string
	// Description returns an explanation of what this rule is validating.
	Description() string
	// Run implements the validation logic, the boolean parameter is true when cmd is "fix"
	Run(*data.ResourceData, bool) []error
}

func GetRule(name string) Rule {
	if rule, ok := Registration[name]; ok {
		return rule
	}

	return nil
}

// Registration for rules
// G<number> -- global/section agnostic rules, e.g. note formatting
// S<number> -- section specific rules, e.g. API section validation
var Registration = map[string]Rule{
	// global
	G001{}.ID(): G001{},
	G002{}.ID(): G002{},

	// section
	S001{}.ID(): S001{},
	S002{}.ID(): S002{},
}
