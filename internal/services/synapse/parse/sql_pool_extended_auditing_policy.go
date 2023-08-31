// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SqlPoolExtendedAuditingPolicyId struct {
	SubscriptionId              string
	ResourceGroup               string
	WorkspaceName               string
	SqlPoolName                 string
	ExtendedAuditingSettingName string
}

func NewSqlPoolExtendedAuditingPolicyID(subscriptionId, resourceGroup, workspaceName, sqlPoolName, extendedAuditingSettingName string) SqlPoolExtendedAuditingPolicyId {
	return SqlPoolExtendedAuditingPolicyId{
		SubscriptionId:              subscriptionId,
		ResourceGroup:               resourceGroup,
		WorkspaceName:               workspaceName,
		SqlPoolName:                 sqlPoolName,
		ExtendedAuditingSettingName: extendedAuditingSettingName,
	}
}

func (id SqlPoolExtendedAuditingPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Extended Auditing Setting Name %q", id.ExtendedAuditingSettingName),
		fmt.Sprintf("Sql Pool Name %q", id.SqlPoolName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sql Pool Extended Auditing Policy", segmentsStr)
}

func (id SqlPoolExtendedAuditingPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/sqlPools/%s/extendedAuditingSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.ExtendedAuditingSettingName)
}

// SqlPoolExtendedAuditingPolicyID parses a SqlPoolExtendedAuditingPolicy ID into an SqlPoolExtendedAuditingPolicyId struct
func SqlPoolExtendedAuditingPolicyID(input string) (*SqlPoolExtendedAuditingPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SqlPoolExtendedAuditingPolicy ID: %+v", input, err)
	}

	resourceId := SqlPoolExtendedAuditingPolicyId{
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
	if resourceId.ExtendedAuditingSettingName, err = id.PopSegment("extendedAuditingSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
