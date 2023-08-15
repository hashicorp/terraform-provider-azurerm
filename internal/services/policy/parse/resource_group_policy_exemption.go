// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceGroupPolicyExemptionId struct {
	SubscriptionId      string
	ResourceGroup       string
	PolicyExemptionName string
}

func NewResourceGroupPolicyExemptionID(subscriptionId, resourceGroup, policyExemptionName string) ResourceGroupPolicyExemptionId {
	return ResourceGroupPolicyExemptionId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		PolicyExemptionName: policyExemptionName,
	}
}

func (id ResourceGroupPolicyExemptionId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Exemption Name %q", id.PolicyExemptionName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Group Policy Exemption", segmentsStr)
}

func (id ResourceGroupPolicyExemptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Authorization/policyExemptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PolicyExemptionName)
}

// ResourceGroupPolicyExemptionID parses a ResourceGroupPolicyExemption ID into an ResourceGroupPolicyExemptionId struct
func ResourceGroupPolicyExemptionID(input string) (*ResourceGroupPolicyExemptionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ResourceGroupPolicyExemption ID: %+v", input, err)
	}

	resourceId := ResourceGroupPolicyExemptionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PolicyExemptionName, err = id.PopSegment("policyExemptions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
