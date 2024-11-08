// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WorkspaceSqlAADAdminId struct {
	SubscriptionId       string
	ResourceGroup        string
	WorkspaceName        string
	SqlAdministratorName string
}

func NewWorkspaceSqlAADAdminID(subscriptionId, resourceGroup, workspaceName, sqlAdministratorName string) WorkspaceSqlAADAdminId {
	return WorkspaceSqlAADAdminId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		WorkspaceName:        workspaceName,
		SqlAdministratorName: sqlAdministratorName,
	}
}

func (id WorkspaceSqlAADAdminId) String() string {
	segments := []string{
		fmt.Sprintf("Sql Administrator Name %q", id.SqlAdministratorName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Workspace Sql A A D Admin", segmentsStr)
}

func (id WorkspaceSqlAADAdminId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/sqlAdministrators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SqlAdministratorName)
}

// WorkspaceSqlAADAdminID parses a WorkspaceSqlAADAdmin ID into an WorkspaceSqlAADAdminId struct
func WorkspaceSqlAADAdminID(input string) (*WorkspaceSqlAADAdminId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an WorkspaceSqlAADAdmin ID: %+v", input, err)
	}

	resourceId := WorkspaceSqlAADAdminId{
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
	if resourceId.SqlAdministratorName, err = id.PopSegment("sqlAdministrators"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
