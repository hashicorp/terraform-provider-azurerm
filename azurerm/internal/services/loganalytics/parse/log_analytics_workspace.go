package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsWorkspaceId struct {
	ResourceGroup string
	Name          string
}

func NewLogAnalyticsWorkspaceID(name, resourceGroup string) LogAnalyticsWorkspaceId {
	return LogAnalyticsWorkspaceId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id LogAnalyticsWorkspaceId) ID(subscriptionId string) string {
	// Log Analytics ID ignores casing
	// Issue tracked here - https://github.com/Azure/azure-sdk-for-go/issues/13268
	return fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/microsoft.operationalinsights/workspaces/%s", subscriptionId, strings.ToLower(id.ResourceGroup), strings.ToLower(id.Name))
}

func LogAnalyticsWorkspaceID(input string) (*LogAnalyticsWorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Log Analytics Workspace ID %q: %+v", input, err)
	}

	server := LogAnalyticsWorkspaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
