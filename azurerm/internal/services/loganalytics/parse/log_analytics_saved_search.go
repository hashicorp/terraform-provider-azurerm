package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsSavedSearchId struct {
	SubscriptionId   string
	ResourceGroup    string
	WorkspaceName    string
	SavedSearcheName string
}

func NewLogAnalyticsSavedSearchID(subscriptionId, resourceGroup, workspaceName, savedSearcheName string) LogAnalyticsSavedSearchId {
	return LogAnalyticsSavedSearchId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		WorkspaceName:    workspaceName,
		SavedSearcheName: savedSearcheName,
	}
}

func (id LogAnalyticsSavedSearchId) String() string {
	segments := []string{
		fmt.Sprintf("Saved Searche Name %q", id.SavedSearcheName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Log Analytics Saved Search", segmentsStr)
}

func (id LogAnalyticsSavedSearchId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/microsoft.operationalinsights/workspaces/%s/savedSearches/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SavedSearcheName)
}

// LogAnalyticsSavedSearchID parses a LogAnalyticsSavedSearch ID into an LogAnalyticsSavedSearchId struct
func LogAnalyticsSavedSearchID(input string) (*LogAnalyticsSavedSearchId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsSavedSearchId{
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
	if resourceId.SavedSearcheName, err = id.PopSegment("savedSearches"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
