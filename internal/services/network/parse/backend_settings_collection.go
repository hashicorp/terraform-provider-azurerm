// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BackendSettingsCollectionId struct {
	SubscriptionId                string
	ResourceGroup                 string
	ApplicationGatewayName        string
	BackendSettingsCollectionName string
}

func NewBackendSettingsCollectionID(subscriptionId, resourceGroup, applicationGatewayName, backendSettingsCollectionName string) BackendSettingsCollectionId {
	return BackendSettingsCollectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		ApplicationGatewayName:        applicationGatewayName,
		BackendSettingsCollectionName: backendSettingsCollectionName,
	}
}

func (id BackendSettingsCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Backend Settings Collection Name %q", id.BackendSettingsCollectionName),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Backend Settings Collection", segmentsStr)
}

func (id BackendSettingsCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/backendSettingsCollection/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.BackendSettingsCollectionName)
}

// BackendSettingsCollectionID parses a BackendSettingsCollection ID into an BackendSettingsCollectionId struct
func BackendSettingsCollectionID(input string) (*BackendSettingsCollectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an BackendSettingsCollection ID: %+v", input, err)
	}

	resourceId := BackendSettingsCollectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ApplicationGatewayName, err = id.PopSegment("applicationGateways"); err != nil {
		return nil, err
	}
	if resourceId.BackendSettingsCollectionName, err = id.PopSegment("backendSettingsCollection"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// BackendSettingsCollectionIDInsensitively parses an BackendSettingsCollection ID into an BackendSettingsCollectionId struct, insensitively
// This should only be used to parse an ID for rewriting, the BackendSettingsCollectionID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func BackendSettingsCollectionIDInsensitively(input string) (*BackendSettingsCollectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := BackendSettingsCollectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
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

	// find the correct casing for the 'backendSettingsCollection' segment
	backendSettingsCollectionKey := "backendSettingsCollection"
	for key := range id.Path {
		if strings.EqualFold(key, backendSettingsCollectionKey) {
			backendSettingsCollectionKey = key
			break
		}
	}
	if resourceId.BackendSettingsCollectionName, err = id.PopSegment(backendSettingsCollectionKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
