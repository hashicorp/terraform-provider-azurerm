package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsLinkedServiceId struct {
	ResourceGroup string
	WorkspaceName string
	Type          string
}

func NewLogAnalyticsLinkedServiceID(resourceGroup, serviceType, workspaceName string) LogAnalyticsLinkedServiceId {
	return LogAnalyticsLinkedServiceId{
		ResourceGroup: resourceGroup,
		WorkspaceName: workspaceName,
		Type:          serviceType,
	}
}

func (id LogAnalyticsLinkedServiceId) ID(subscriptionId string) string {
	// Log Analytics ID ignores casing
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/linkedServices/%s", subscriptionId, id.ResourceGroup, id.WorkspaceName, id.Type)
}

func LogAnalyticsLinkedServiceID(input string) (*LogAnalyticsLinkedServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Log Analytics Linked Service ID %q: %+v", input, err)
	}

	linkedService := LogAnalyticsLinkedServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if linkedService.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if linkedService.Type, err = id.PopSegment("linkedServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &linkedService, nil
}
