// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type EndpointServiceBusTopicId struct {
	SubscriptionId string
	ResourceGroup  string
	IotHubName     string
	EndpointName   string
}

func NewEndpointServiceBusTopicID(subscriptionId, resourceGroup, iotHubName, endpointName string) EndpointServiceBusTopicId {
	return EndpointServiceBusTopicId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotHubName:     iotHubName,
		EndpointName:   endpointName,
	}
}

func (id EndpointServiceBusTopicId) String() string {
	segments := []string{
		fmt.Sprintf("Endpoint Name %q", id.EndpointName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Endpoint Service Bus Topic", segmentsStr)
}

func (id EndpointServiceBusTopicId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/iotHubs/%s/endpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.EndpointName)
}

// EndpointServiceBusTopicID parses a EndpointServiceBusTopic ID into an EndpointServiceBusTopicId struct
func EndpointServiceBusTopicID(input string) (*EndpointServiceBusTopicId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an EndpointServiceBusTopic ID: %+v", input, err)
	}

	resourceId := EndpointServiceBusTopicId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IotHubName, err = id.PopSegment("iotHubs"); err != nil {
		return nil, err
	}
	if resourceId.EndpointName, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// EndpointServiceBusTopicIDInsensitively parses an EndpointServiceBusTopic ID into an EndpointServiceBusTopicId struct, insensitively
// This should only be used to parse an ID for rewriting, the EndpointServiceBusTopicID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func EndpointServiceBusTopicIDInsensitively(input string) (*EndpointServiceBusTopicId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EndpointServiceBusTopicId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'iotHubs' segment
	iotHubsKey := "iotHubs"
	for key := range id.Path {
		if strings.EqualFold(key, iotHubsKey) {
			iotHubsKey = key
			break
		}
	}
	if resourceId.IotHubName, err = id.PopSegment(iotHubsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'endpoints' segment
	endpointsKey := "endpoints"
	for key := range id.Path {
		if strings.EqualFold(key, endpointsKey) {
			endpointsKey = key
			break
		}
	}
	if resourceId.EndpointName, err = id.PopSegment(endpointsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
