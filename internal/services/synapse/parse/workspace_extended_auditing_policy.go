// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WorkspaceExtendedAuditingPolicyId struct {
	SubscriptionId              string
	ResourceGroup               string
	WorkspaceName               string
	ExtendedAuditingSettingName string
}

func NewWorkspaceExtendedAuditingPolicyID(subscriptionId, resourceGroup, workspaceName, extendedAuditingSettingName string) WorkspaceExtendedAuditingPolicyId {
	return WorkspaceExtendedAuditingPolicyId{
		SubscriptionId:              subscriptionId,
		ResourceGroup:               resourceGroup,
		WorkspaceName:               workspaceName,
		ExtendedAuditingSettingName: extendedAuditingSettingName,
	}
}

func (id WorkspaceExtendedAuditingPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Extended Auditing Setting Name %q", id.ExtendedAuditingSettingName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Workspace Extended Auditing Policy", segmentsStr)
}

func (id WorkspaceExtendedAuditingPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/extendedAuditingSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.ExtendedAuditingSettingName)
}

// WorkspaceExtendedAuditingPolicyID parses a WorkspaceExtendedAuditingPolicy ID into an WorkspaceExtendedAuditingPolicyId struct
func WorkspaceExtendedAuditingPolicyID(input string) (*WorkspaceExtendedAuditingPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an WorkspaceExtendedAuditingPolicy ID: %+v", input, err)
	}

	resourceId := WorkspaceExtendedAuditingPolicyId{
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
	if resourceId.ExtendedAuditingSettingName, err = id.PopSegment("extendedAuditingSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
