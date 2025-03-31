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

type ManagedPrivateEndpointId struct {
	SubscriptionId            string
	ResourceGroup             string
	FactoryName               string
	ManagedVirtualNetworkName string
	Name                      string
}

func NewManagedPrivateEndpointID(subscriptionId, resourceGroup, factoryName, managedVirtualNetworkName, name string) ManagedPrivateEndpointId {
	return ManagedPrivateEndpointId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		FactoryName:               factoryName,
		ManagedVirtualNetworkName: managedVirtualNetworkName,
		Name:                      name,
	}
}

func (id ManagedPrivateEndpointId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Managed Virtual Network Name %q", id.ManagedVirtualNetworkName),
		fmt.Sprintf("Factory Name %q", id.FactoryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed Private Endpoint", segmentsStr)
}

func (id ManagedPrivateEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/managedVirtualNetworks/%s/managedPrivateEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FactoryName, id.ManagedVirtualNetworkName, id.Name)
}

// ManagedPrivateEndpointID parses a ManagedPrivateEndpoint ID into an ManagedPrivateEndpointId struct
func ManagedPrivateEndpointID(input string) (*ManagedPrivateEndpointId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedPrivateEndpoint ID: %+v", input, err)
	}

	resourceId := ManagedPrivateEndpointId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FactoryName, err = id.PopSegment("factories"); err != nil {
		return nil, err
	}
	if resourceId.ManagedVirtualNetworkName, err = id.PopSegment("managedVirtualNetworks"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("managedPrivateEndpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
