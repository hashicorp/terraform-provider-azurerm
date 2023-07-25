// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BackendHttpSettingsCollectionId struct {
	SubscriptionId                    string
	ResourceGroup                     string
	ApplicationGatewayName            string
	BackendHttpSettingsCollectionName string
}

func NewBackendHttpSettingsCollectionID(subscriptionId, resourceGroup, applicationGatewayName, backendHttpSettingsCollectionName string) BackendHttpSettingsCollectionId {
	return BackendHttpSettingsCollectionId{
		SubscriptionId:                    subscriptionId,
		ResourceGroup:                     resourceGroup,
		ApplicationGatewayName:            applicationGatewayName,
		BackendHttpSettingsCollectionName: backendHttpSettingsCollectionName,
	}
}

func (id BackendHttpSettingsCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Backend Http Settings Collection Name %q", id.BackendHttpSettingsCollectionName),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Backend Http Settings Collection", segmentsStr)
}

func (id BackendHttpSettingsCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/backendHttpSettingsCollection/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.BackendHttpSettingsCollectionName)
}

// BackendHttpSettingsCollectionID parses a BackendHttpSettingsCollection ID into an BackendHttpSettingsCollectionId struct
func BackendHttpSettingsCollectionID(input string) (*BackendHttpSettingsCollectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an BackendHttpSettingsCollection ID: %+v", input, err)
	}

	resourceId := BackendHttpSettingsCollectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ApplicationGatewayName, err = id.PopSegment("applicationGateways"); err != nil {
		return nil, err
	}
	if resourceId.BackendHttpSettingsCollectionName, err = id.PopSegment("backendHttpSettingsCollection"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// BackendHttpSettingsCollectionIDInsensitively parses an BackendHttpSettingsCollection ID into an BackendHttpSettingsCollectionId struct, insensitively
// This should only be used to parse an ID for rewriting, the BackendHttpSettingsCollectionID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func BackendHttpSettingsCollectionIDInsensitively(input string) (*BackendHttpSettingsCollectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := BackendHttpSettingsCollectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'applicationGateways' segment
	applicationGatewaysKey := "applicationGateways"
	for key := range id.Path {
		if strings.EqualFold(key, applicationGatewaysKey) {
			applicationGatewaysKey = key
			break
		}
	}
	if resourceId.ApplicationGatewayName, err = id.PopSegment(applicationGatewaysKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'backendHttpSettingsCollection' segment
	backendHttpSettingsCollectionKey := "backendHttpSettingsCollection"
	for key := range id.Path {
		if strings.EqualFold(key, backendHttpSettingsCollectionKey) {
			backendHttpSettingsCollectionKey = key
			break
		}
	}
	if resourceId.BackendHttpSettingsCollectionName, err = id.PopSegment(backendHttpSettingsCollectionKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
