// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IotSecuritySolutionId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewIotSecuritySolutionID(subscriptionId, resourceGroup, name string) IotSecuritySolutionId {
	return IotSecuritySolutionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id IotSecuritySolutionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Iot Security Solution", segmentsStr)
}

func (id IotSecuritySolutionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Security/iotSecuritySolutions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// IotSecuritySolutionID parses a IotSecuritySolution ID into an IotSecuritySolutionId struct
func IotSecuritySolutionID(input string) (*IotSecuritySolutionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an IotSecuritySolution ID: %+v", input, err)
	}

	resourceId := IotSecuritySolutionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("iotSecuritySolutions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// IotSecuritySolutionIDInsensitively parses an IotSecuritySolution ID into an IotSecuritySolutionId struct, insensitively
// This should only be used to parse an ID for rewriting, the IotSecuritySolutionID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func IotSecuritySolutionIDInsensitively(input string) (*IotSecuritySolutionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IotSecuritySolutionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'iotSecuritySolutions' segment
	iotSecuritySolutionsKey := "iotSecuritySolutions"
	for key := range id.Path {
		if strings.EqualFold(key, iotSecuritySolutionsKey) {
			iotSecuritySolutionsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(iotSecuritySolutionsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
