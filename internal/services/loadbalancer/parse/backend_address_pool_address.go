// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BackendAddressPoolAddressId struct {
	SubscriptionId         string
	ResourceGroup          string
	LoadBalancerName       string
	BackendAddressPoolName string
	AddressName            string
}

func NewBackendAddressPoolAddressID(subscriptionId, resourceGroup, loadBalancerName, backendAddressPoolName, addressName string) BackendAddressPoolAddressId {
	return BackendAddressPoolAddressId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		LoadBalancerName:       loadBalancerName,
		BackendAddressPoolName: backendAddressPoolName,
		AddressName:            addressName,
	}
}

func (id BackendAddressPoolAddressId) String() string {
	segments := []string{
		fmt.Sprintf("Address Name %q", id.AddressName),
		fmt.Sprintf("Backend Address Pool Name %q", id.BackendAddressPoolName),
		fmt.Sprintf("Load Balancer Name %q", id.LoadBalancerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Backend Address Pool Address", segmentsStr)
}

func (id BackendAddressPoolAddressId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/backendAddressPools/%s/addresses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, id.AddressName)
}

// BackendAddressPoolAddressID parses a BackendAddressPoolAddress ID into an BackendAddressPoolAddressId struct
func BackendAddressPoolAddressID(input string) (*BackendAddressPoolAddressId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an BackendAddressPoolAddress ID: %+v", input, err)
	}

	resourceId := BackendAddressPoolAddressId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}
	if resourceId.BackendAddressPoolName, err = id.PopSegment("backendAddressPools"); err != nil {
		return nil, err
	}
	if resourceId.AddressName, err = id.PopSegment("addresses"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
