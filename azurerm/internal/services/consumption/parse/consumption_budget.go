package parse

// NOTE: this file is manually generated due to the structure of consumption budgets IDs

import (
	"fmt"
	"strings"
)

const (
	consumptionBudgetIDSeparator = "/providers/Microsoft.Consumption/budgets/"
)

type ConsumptionBudgetId struct {
	Name  string
	Scope string
}

func NewConsumptionBudgetID(name, scope string) ConsumptionBudgetId {
	return ConsumptionBudgetId{
		Name:  name,
		Scope: scope,
	}
}

func (id ConsumptionBudgetId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Scope %q", id.Scope),
	}

	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Consumption Budget", segmentsStr)
}

func ConsumptionBudgetID(input string) (*ConsumptionBudgetId, error) {
	if !strings.Contains(input, consumptionBudgetIDSeparator) {
		return nil, fmt.Errorf("the provided ID %q does not contain the expected resource provider %s", input, consumptionBudgetIDSeparator)
	}

	components := strings.Split(input, consumptionBudgetIDSeparator)

	if len(components) != 2 {
		return nil, fmt.Errorf("the provided ID %q is not in the expected format", input)
	}

	name := components[1]
	scope := components[0]

	if name == "" || scope == "" {
		return nil, fmt.Errorf("name or scope cannot be empty strings. Name: '%s', Scope: '%s'", name, scope)
	}

	return &ConsumptionBudgetId{
		Name:  components[1],
		Scope: components[0],
	}, nil
}
