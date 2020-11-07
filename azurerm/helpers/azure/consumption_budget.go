package azure

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
	"strings"
)

const (
	consumptionBudgetIDSeparator = "/providers/Microsoft.Consumption/budgets/"
)

type ConsumptionBudgetID struct {
	Name  string
	Scope string
}

// validation
func ValidateAzureConsumptionBudgetName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[-_a-zA-Z0-9]{1,63}$"),
		"The consumption budget name can contain only letters, numbers, underscores, and hyphens. The consumption budget name be between 6 and 63 characters long.",
	)
}

// /{scope}/providers/Microsoft.Consumption/budgets/{budgetName}
func ParseAzureConsumptionBudgetID(id string) (*ConsumptionBudgetID, error) {
	if !strings.Contains(id, consumptionBudgetIDSeparator) {
		return nil, fmt.Errorf("the provided ID %q does not contain the expected resource provider %s", id, consumptionBudgetIDSeparator)
	}

	components := strings.Split(id, consumptionBudgetIDSeparator)

	if len(components) != 2 {
		return nil, fmt.Errorf("the provided ID %q is not in the expected format", id)
	}

	name := components[1]
	scope := components[0]

	if name == "" || scope == "" {
		return nil, fmt.Errorf("name or scope cannot be empty strings. Name: '%s', Scope: '%s'", name, scope)
	}

	return &ConsumptionBudgetID{
		Name:  components[1],
		Scope: components[0],
	}, nil
}
