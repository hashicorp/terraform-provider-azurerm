// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceGroupPolicyRemediationId struct {
	SubscriptionId  string
	ResourceGroup   string
	RemediationName string
}

func NewResourceGroupPolicyRemediationID(subscriptionId, resourceGroup, remediationName string) ResourceGroupPolicyRemediationId {
	return ResourceGroupPolicyRemediationId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		RemediationName: remediationName,
	}
}

func (id ResourceGroupPolicyRemediationId) String() string {
	segments := []string{
		fmt.Sprintf("Remediation Name %q", id.RemediationName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Group Policy Remediation", segmentsStr)
}

func (id ResourceGroupPolicyRemediationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RemediationName)
}

// ResourceGroupPolicyRemediationID parses a ResourceGroupPolicyRemediation ID into an ResourceGroupPolicyRemediationId struct
func ResourceGroupPolicyRemediationID(input string) (*ResourceGroupPolicyRemediationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ResourceGroupPolicyRemediation ID: %+v", input, err)
	}

	resourceId := ResourceGroupPolicyRemediationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.RemediationName, err = id.PopSegment("remediations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
