// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SqlPoolWorkloadGroupId struct {
	SubscriptionId    string
	ResourceGroup     string
	WorkspaceName     string
	SqlPoolName       string
	WorkloadGroupName string
}

func NewSqlPoolWorkloadGroupID(subscriptionId, resourceGroup, workspaceName, sqlPoolName, workloadGroupName string) SqlPoolWorkloadGroupId {
	return SqlPoolWorkloadGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		WorkspaceName:     workspaceName,
		SqlPoolName:       sqlPoolName,
		WorkloadGroupName: workloadGroupName,
	}
}

func (id SqlPoolWorkloadGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Workload Group Name %q", id.WorkloadGroupName),
		fmt.Sprintf("Sql Pool Name %q", id.SqlPoolName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sql Pool Workload Group", segmentsStr)
}

func (id SqlPoolWorkloadGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/sqlPools/%s/workloadGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName)
}

// SqlPoolWorkloadGroupID parses a SqlPoolWorkloadGroup ID into an SqlPoolWorkloadGroupId struct
func SqlPoolWorkloadGroupID(input string) (*SqlPoolWorkloadGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SqlPoolWorkloadGroup ID: %+v", input, err)
	}

	resourceId := SqlPoolWorkloadGroupId{
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
	if resourceId.SqlPoolName, err = id.PopSegment("sqlPools"); err != nil {
		return nil, err
	}
	if resourceId.WorkloadGroupName, err = id.PopSegment("workloadGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
