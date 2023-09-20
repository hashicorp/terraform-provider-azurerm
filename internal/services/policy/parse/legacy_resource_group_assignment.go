// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceGroupAssignmentId struct {
	SubscriptionId       string
	ResourceGroup        string
	PolicyAssignmentName string
}

func NewResourceGroupAssignmentID(subscriptionId, resourceGroup, policyAssignmentName string) ResourceGroupAssignmentId {
	return ResourceGroupAssignmentId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		PolicyAssignmentName: policyAssignmentName,
	}
}

func (id ResourceGroupAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Assignment Name %q", id.PolicyAssignmentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Group Assignment", segmentsStr)
}

func (id ResourceGroupAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Authorization/policyAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PolicyAssignmentName)
}

// ResourceGroupAssignmentID parses a ResourceGroupAssignment ID into an ResourceGroupAssignmentId struct
func ResourceGroupAssignmentID(input string) (*ResourceGroupAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ResourceGroupAssignment ID: %+v", input, err)
	}

	resourceId := ResourceGroupAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PolicyAssignmentName, err = id.PopSegment("policyAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
