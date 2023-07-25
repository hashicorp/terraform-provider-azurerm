// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type UrlPathMapId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	Name                   string
}

func NewUrlPathMapID(subscriptionId, resourceGroup, applicationGatewayName, name string) UrlPathMapId {
	return UrlPathMapId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ApplicationGatewayName: applicationGatewayName,
		Name:                   name,
	}
}

func (id UrlPathMapId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Url Path Map", segmentsStr)
}

func (id UrlPathMapId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/urlPathMaps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.Name)
}

// UrlPathMapID parses a UrlPathMap ID into an UrlPathMapId struct
func UrlPathMapID(input string) (*UrlPathMapId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an UrlPathMap ID: %+v", input, err)
	}

	resourceId := UrlPathMapId{
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
	if resourceId.Name, err = id.PopSegment("urlPathMaps"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// UrlPathMapIDInsensitively parses an UrlPathMap ID into an UrlPathMapId struct, insensitively
// This should only be used to parse an ID for rewriting, the UrlPathMapID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func UrlPathMapIDInsensitively(input string) (*UrlPathMapId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := UrlPathMapId{
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

	// find the correct casing for the 'urlPathMaps' segment
	urlPathMapsKey := "urlPathMaps"
	for key := range id.Path {
		if strings.EqualFold(key, urlPathMapsKey) {
			urlPathMapsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(urlPathMapsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
