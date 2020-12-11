package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsLinkedServiceId struct {
	SubscriptionId    string
	ResourceGroup     string
	WorkspaceName     string
	LinkedServiceName string
}

func NewLogAnalyticsLinkedServiceID(subscriptionId, resourceGroup, workspaceName, linkedServiceName string) LogAnalyticsLinkedServiceId {
	return LogAnalyticsLinkedServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		WorkspaceName:     workspaceName,
		LinkedServiceName: linkedServiceName,
	}
}

func (id LogAnalyticsLinkedServiceId) String() string {
	segments := []string{
		fmt.Sprintf("Linked Service Name %q", id.LinkedServiceName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Log Analytics Linked Service", segmentsStr)
}

func (id LogAnalyticsLinkedServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/linkedServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.LinkedServiceName)
}

// LogAnalyticsLinkedServiceID parses a LogAnalyticsLinkedService ID into an LogAnalyticsLinkedServiceId struct
func LogAnalyticsLinkedServiceID(input string) (*LogAnalyticsLinkedServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsLinkedServiceId{
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
	if resourceId.LinkedServiceName, err = id.PopSegment("linkedServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
