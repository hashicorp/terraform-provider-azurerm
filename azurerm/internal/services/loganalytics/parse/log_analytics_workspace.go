package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsWorkspaceId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkspaceName  string
}

func NewLogAnalyticsWorkspaceID(subscriptionId, resourceGroup, workspaceName string) LogAnalyticsWorkspaceId {
	return LogAnalyticsWorkspaceId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkspaceName:  workspaceName,
	}
}

func (id LogAnalyticsWorkspaceId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
	}
	return strings.Join(segments, " / ")
}

func (id LogAnalyticsWorkspaceId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
}

// LogAnalyticsWorkspaceID parses a LogAnalyticsWorkspace ID into an LogAnalyticsWorkspaceId struct
func LogAnalyticsWorkspaceID(input string) (*LogAnalyticsWorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsWorkspaceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
