// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ServicePlanId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerfarmName string
}

func NewServicePlanID(subscriptionId, resourceGroup, serverfarmName string) ServicePlanId {
	return ServicePlanId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerfarmName: serverfarmName,
	}
}

func (id ServicePlanId) String() string {
	segments := []string{
		fmt.Sprintf("Serverfarm Name %q", id.ServerfarmName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Service Plan", segmentsStr)
}

func (id ServicePlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/serverfarms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerfarmName)
}

// ServicePlanID parses a ServicePlan ID into an ServicePlanId struct
func ServicePlanID(input string) (*ServicePlanId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ServicePlan ID: %+v", input, err)
	}

	resourceId := ServicePlanId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerfarmName, err = id.PopSegment("serverfarms"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ServicePlanIDInsensitively parses an ServicePlan ID into an ServicePlanId struct, insensitively
// This should only be used to parse an ID for rewriting, the ServicePlanID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ServicePlanIDInsensitively(input string) (*ServicePlanId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ServicePlanId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'serverfarms' segment
	serverfarmsKey := "serverfarms"
	for key := range id.Path {
		if strings.EqualFold(key, serverfarmsKey) {
			serverfarmsKey = key
			break
		}
	}
	if resourceId.ServerfarmName, err = id.PopSegment(serverfarmsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
