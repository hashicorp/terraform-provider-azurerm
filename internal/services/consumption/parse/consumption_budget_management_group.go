package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type ConsumptionBudgetManagementGroupId struct {
	ManagementGroupName string
	BudgetName          string
}

func NewConsumptionBudgetManagementGroupID(managementGroupName, budgetName string) ConsumptionBudgetManagementGroupId {
	return ConsumptionBudgetManagementGroupId{
		ManagementGroupName: managementGroupName,
		BudgetName:          budgetName,
	}
}

func (id ConsumptionBudgetManagementGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Budget Name %q", id.BudgetName),
		fmt.Sprintf("Management Group Name %q", id.ManagementGroupName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Consumption Budget Management Group", segmentsStr)
}

func (id ConsumptionBudgetManagementGroupId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Consumption/budgets/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupName, id.BudgetName)
}

// ConsumptionBudgetManagementGroupID parses a ConsumptionBudgetManagementGroup ID into an ConsumptionBudgetManagementGroupId struct
func ConsumptionBudgetManagementGroupID(input string) (*ConsumptionBudgetManagementGroupId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConsumptionBudgetManagementGroupId{}

	if resourceId.ManagementGroupName, err = id.PopSegment("managementGroups"); err != nil {
		return nil, err
	}

	if resourceId.BudgetName, err = id.PopSegment("budgets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
