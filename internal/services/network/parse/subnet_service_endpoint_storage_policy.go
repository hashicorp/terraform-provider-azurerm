// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubnetServiceEndpointStoragePolicyId struct {
	SubscriptionId            string
	ResourceGroup             string
	ServiceEndpointPolicyName string
}

func NewSubnetServiceEndpointStoragePolicyID(subscriptionId, resourceGroup, serviceEndpointPolicyName string) SubnetServiceEndpointStoragePolicyId {
	return SubnetServiceEndpointStoragePolicyId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		ServiceEndpointPolicyName: serviceEndpointPolicyName,
	}
}

func (id SubnetServiceEndpointStoragePolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Service Endpoint Policy Name %q", id.ServiceEndpointPolicyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subnet Service Endpoint Storage Policy", segmentsStr)
}

func (id SubnetServiceEndpointStoragePolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/serviceEndpointPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceEndpointPolicyName)
}

// SubnetServiceEndpointStoragePolicyID parses a SubnetServiceEndpointStoragePolicy ID into an SubnetServiceEndpointStoragePolicyId struct
func SubnetServiceEndpointStoragePolicyID(input string) (*SubnetServiceEndpointStoragePolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SubnetServiceEndpointStoragePolicy ID: %+v", input, err)
	}

	resourceId := SubnetServiceEndpointStoragePolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceEndpointPolicyName, err = id.PopSegment("serviceEndpointPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
