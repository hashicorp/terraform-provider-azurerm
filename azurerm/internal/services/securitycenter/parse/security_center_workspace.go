package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterWorkspaceId struct {
	SubscriptionId       string
	WorkspaceSettingName string
}

func NewSecurityCenterWorkspaceID(subscriptionId, workspaceSettingName string) SecurityCenterWorkspaceId {
	return SecurityCenterWorkspaceId{
		SubscriptionId:       subscriptionId,
		WorkspaceSettingName: workspaceSettingName,
	}
}

func (id SecurityCenterWorkspaceId) String() string {
	segments := []string{
		fmt.Sprintf("Workspace Setting Name %q", id.WorkspaceSettingName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Security Center Workspace", segmentsStr)
}

func (id SecurityCenterWorkspaceId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/workspaceSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.WorkspaceSettingName)
}

// SecurityCenterWorkspaceID parses a SecurityCenterWorkspace ID into an SecurityCenterWorkspaceId struct
func SecurityCenterWorkspaceID(input string) (*SecurityCenterWorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SecurityCenterWorkspaceId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.WorkspaceSettingName, err = id.PopSegment("workspaceSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
