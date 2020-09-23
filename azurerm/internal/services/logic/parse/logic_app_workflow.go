package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogicAppWorkflowId struct {
	ResourceGroup string
	Name          string
}

func NewLogicAppWorkflowID(resourceGroup, name string) LogicAppWorkflowId {
	return LogicAppWorkflowId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id LogicAppWorkflowId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func LogicAppWorkflowID(input string) (*LogicAppWorkflowId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Logic App Workflow ID %q: %+v", input, err)
	}

	workflow := LogicAppWorkflowId{
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
