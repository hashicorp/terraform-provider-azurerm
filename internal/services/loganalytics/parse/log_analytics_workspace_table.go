package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LogAnalyticsWorkspaceTableId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkspaceName  string
	TableName      string
}

func NewLogAnalyticsWorkspaceTableID(subscriptionId, resourceGroup, workspaceName, tableName string) LogAnalyticsWorkspaceTableId {
	return LogAnalyticsWorkspaceTableId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkspaceName:  workspaceName,
		TableName:      tableName,
	}
}

func (id LogAnalyticsWorkspaceTableId) String() string {
	segments := []string{
		fmt.Sprintf("Table Name %q", id.TableName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Log Analytics Workspace Table", segmentsStr)
}

func (id LogAnalyticsWorkspaceTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/tables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.TableName)
}

// LogAnalyticsWorkspaceTableID parses a LogAnalyticsWorkspace ID into an LogAnalyticsWorkspaceId struct
func LogAnalyticsWorkspaceTableID(input string) (*LogAnalyticsWorkspaceTableId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsWorkspaceTableId{
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

	if resourceId.TableName, err = id.PopSegment("tables"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
