// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SharedAccessPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	IotHubName     string
	IotHubKeyName  string
}

func NewSharedAccessPolicyID(subscriptionId, resourceGroup, iotHubName, iotHubKeyName string) SharedAccessPolicyId {
	return SharedAccessPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotHubName:     iotHubName,
		IotHubKeyName:  iotHubKeyName,
	}
}

func (id SharedAccessPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Iot Hub Key Name %q", id.IotHubKeyName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Shared Access Policy", segmentsStr)
}

func (id SharedAccessPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/iotHubs/%s/iotHubKeys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.IotHubKeyName)
}

// SharedAccessPolicyID parses a SharedAccessPolicy ID into an SharedAccessPolicyId struct
func SharedAccessPolicyID(input string) (*SharedAccessPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SharedAccessPolicy ID: %+v", input, err)
	}

	resourceId := SharedAccessPolicyId{
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
	if resourceId.IotHubKeyName, err = id.PopSegment("iotHubKeys"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SharedAccessPolicyIDInsensitively parses an SharedAccessPolicy ID into an SharedAccessPolicyId struct, insensitively
// This should only be used to parse an ID for rewriting, the SharedAccessPolicyID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SharedAccessPolicyIDInsensitively(input string) (*SharedAccessPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SharedAccessPolicyId{
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

	// find the correct casing for the 'iotHubKeys' segment
	iotHubKeysKey := "iotHubKeys"
	for key := range id.Path {
		if strings.EqualFold(key, iotHubKeysKey) {
			iotHubKeysKey = key
			break
		}
	}
	if resourceId.IotHubKeyName, err = id.PopSegment(iotHubKeysKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
