// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WorkspaceAADAdminId struct {
	SubscriptionId    string
	ResourceGroup     string
	WorkspaceName     string
	AdministratorName string
}

func NewWorkspaceAADAdminID(subscriptionId, resourceGroup, workspaceName, administratorName string) WorkspaceAADAdminId {
	return WorkspaceAADAdminId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		WorkspaceName:     workspaceName,
		AdministratorName: administratorName,
	}
}

func (id WorkspaceAADAdminId) String() string {
	segments := []string{
		fmt.Sprintf("Administrator Name %q", id.AdministratorName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Workspace A A D Admin", segmentsStr)
}

func (id WorkspaceAADAdminId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/administrators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.AdministratorName)
}

// WorkspaceAADAdminID parses a WorkspaceAADAdmin ID into an WorkspaceAADAdminId struct
func WorkspaceAADAdminID(input string) (*WorkspaceAADAdminId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an WorkspaceAADAdmin ID: %+v", input, err)
	}

	resourceId := WorkspaceAADAdminId{
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
	if resourceId.AdministratorName, err = id.PopSegment("administrators"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
