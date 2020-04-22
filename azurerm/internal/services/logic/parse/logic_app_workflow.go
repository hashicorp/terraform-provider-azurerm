package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogicAppWorkflowId struct {
	Subscription  string
	ResourceGroup string
	Name          string
}

func (id LogicAppWorkflowId) String() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s",
		id.Subscription, id.ResourceGroup, id.Name)
}

func LogicAppWorkflowID(input string) (*LogicAppWorkflowId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Logic App Workflow ID %q: %+v", input, err)
	}

	workflow := LogicAppWorkflowId{
		Subscription:  id.SubscriptionID,
		ResourceGroup: id.ResourceGroup,
	}

	if workflow.Name, err = id.PopSegment("workflows"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &workflow, nil
}
