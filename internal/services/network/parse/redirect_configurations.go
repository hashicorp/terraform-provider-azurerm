// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RedirectConfigurationsId struct {
	SubscriptionId            string
	ResourceGroup             string
	ApplicationGatewayName    string
	RedirectConfigurationName string
}

func NewRedirectConfigurationsID(subscriptionId, resourceGroup, applicationGatewayName, redirectConfigurationName string) RedirectConfigurationsId {
	return RedirectConfigurationsId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		ApplicationGatewayName:    applicationGatewayName,
		RedirectConfigurationName: redirectConfigurationName,
	}
}

func (id RedirectConfigurationsId) String() string {
	segments := []string{
		fmt.Sprintf("Redirect Configuration Name %q", id.RedirectConfigurationName),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Redirect Configurations", segmentsStr)
}

func (id RedirectConfigurationsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/redirectConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.RedirectConfigurationName)
}

// RedirectConfigurationsID parses a RedirectConfigurations ID into an RedirectConfigurationsId struct
func RedirectConfigurationsID(input string) (*RedirectConfigurationsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an RedirectConfigurations ID: %+v", input, err)
	}

	resourceId := RedirectConfigurationsId{
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
	if resourceId.RedirectConfigurationName, err = id.PopSegment("redirectConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// RedirectConfigurationsIDInsensitively parses an RedirectConfigurations ID into an RedirectConfigurationsId struct, insensitively
// This should only be used to parse an ID for rewriting, the RedirectConfigurationsID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func RedirectConfigurationsIDInsensitively(input string) (*RedirectConfigurationsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RedirectConfigurationsId{
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

	// find the correct casing for the 'redirectConfigurations' segment
	redirectConfigurationsKey := "redirectConfigurations"
	for key := range id.Path {
		if strings.EqualFold(key, redirectConfigurationsKey) {
			redirectConfigurationsKey = key
			break
		}
	}
	if resourceId.RedirectConfigurationName, err = id.PopSegment(redirectConfigurationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
