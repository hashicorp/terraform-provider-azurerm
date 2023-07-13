// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SqlPoolSecurityAlertPolicyId struct {
	SubscriptionId          string
	ResourceGroup           string
	WorkspaceName           string
	SqlPoolName             string
	SecurityAlertPolicyName string
}

func NewSqlPoolSecurityAlertPolicyID(subscriptionId, resourceGroup, workspaceName, sqlPoolName, securityAlertPolicyName string) SqlPoolSecurityAlertPolicyId {
	return SqlPoolSecurityAlertPolicyId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		WorkspaceName:           workspaceName,
		SqlPoolName:             sqlPoolName,
		SecurityAlertPolicyName: securityAlertPolicyName,
	}
}

func (id SqlPoolSecurityAlertPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Security Alert Policy Name %q", id.SecurityAlertPolicyName),
		fmt.Sprintf("Sql Pool Name %q", id.SqlPoolName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sql Pool Security Alert Policy", segmentsStr)
}

func (id SqlPoolSecurityAlertPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/sqlPools/%s/securityAlertPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.SecurityAlertPolicyName)
}

// SqlPoolSecurityAlertPolicyID parses a SqlPoolSecurityAlertPolicy ID into an SqlPoolSecurityAlertPolicyId struct
func SqlPoolSecurityAlertPolicyID(input string) (*SqlPoolSecurityAlertPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SqlPoolSecurityAlertPolicy ID: %+v", input, err)
	}

	resourceId := SqlPoolSecurityAlertPolicyId{
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
	if resourceId.SecurityAlertPolicyName, err = id.PopSegment("securityAlertPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
