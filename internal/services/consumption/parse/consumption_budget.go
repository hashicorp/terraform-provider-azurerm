package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type ConsumptionBudgetId struct {
	Name  string
	Scope string
}

func (id ConsumptionBudgetId) String() string {
	segments := []string{
		fmt.Sprintf("Consumption Budget Name %q", id.Name),
		fmt.Sprintf("Scope %q", id.Scope),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Consumption Budget ID", segmentsStr)
}

func (id ConsumptionBudgetId) ID() string {
	fmtString := "%s/providers/Microsoft.Consumption/budgets/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.Name)
}

func NewConsumptionBudgetId(scope, name string) ConsumptionBudgetId {
	return ConsumptionBudgetId{
		Name:  name,
		Scope: scope,
	}
}

func ConsumptionBudgetID(input string) (*ConsumptionBudgetId, error) {
	// in general, the id of a budget should be:
	// {scope}/providers/Microsoft.Consumption/budgets/{name}
	regex := regexp.MustCompile(`/providers/Microsoft\.Consumption/budgets/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Consumption Budget ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Consumption Budget ID %q: Expected 2 segments after split", input)
	}

	scope := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Consumption Budget ID %q: budget name is empty", input)
	}

	return &ConsumptionBudgetId{
		Name:  name,
		Scope: scope,
	}, nil
}
