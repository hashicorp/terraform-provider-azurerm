package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ConsumptionBudgetSubscriptionId struct {
	SubscriptionId string
	BudgetName     string
}

func NewConsumptionBudgetSubscriptionID(subscriptionId, budgetName string) ConsumptionBudgetSubscriptionId {
	return ConsumptionBudgetSubscriptionId{
		SubscriptionId: subscriptionId,
		BudgetName:     budgetName,
	}
}

func (id ConsumptionBudgetSubscriptionId) String() string {
	segments := []string{
		fmt.Sprintf("Budget Name %q", id.BudgetName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Consumption Budget Subscription", segmentsStr)
}

func (id ConsumptionBudgetSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Consumption/budgets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.BudgetName)
}

// ConsumptionBudgetSubscriptionID parses a ConsumptionBudgetSubscription ID into an ConsumptionBudgetSubscriptionId struct
func ConsumptionBudgetSubscriptionID(input string) (*ConsumptionBudgetSubscriptionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConsumptionBudgetSubscriptionId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.BudgetName, err = id.PopSegment("budgets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
