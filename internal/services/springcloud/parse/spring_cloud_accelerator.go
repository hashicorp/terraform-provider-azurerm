// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudAcceleratorId struct {
	SubscriptionId             string
	ResourceGroup              string
	SpringName                 string
	ApplicationAcceleratorName string
}

func NewSpringCloudAcceleratorID(subscriptionId, resourceGroup, springName, applicationAcceleratorName string) SpringCloudAcceleratorId {
	return SpringCloudAcceleratorId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		SpringName:                 springName,
		ApplicationAcceleratorName: applicationAcceleratorName,
	}
}

func (id SpringCloudAcceleratorId) String() string {
	segments := []string{
		fmt.Sprintf("Application Accelerator Name %q", id.ApplicationAcceleratorName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Accelerator", segmentsStr)
}

func (id SpringCloudAcceleratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/applicationAccelerators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName)
}

// SpringCloudAcceleratorID parses a SpringCloudAccelerator ID into an SpringCloudAcceleratorId struct
func SpringCloudAcceleratorID(input string) (*SpringCloudAcceleratorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudAccelerator ID: %+v", input, err)
	}

	resourceId := SpringCloudAcceleratorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("spring"); err != nil {
		return nil, err
	}
	if resourceId.ApplicationAcceleratorName, err = id.PopSegment("applicationAccelerators"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudAcceleratorIDInsensitively parses an SpringCloudAccelerator ID into an SpringCloudAcceleratorId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudAcceleratorID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudAcceleratorIDInsensitively(input string) (*SpringCloudAcceleratorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudAcceleratorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'spring' segment
	springKey := "spring"
	for key := range id.Path {
		if strings.EqualFold(key, springKey) {
			springKey = key
			break
		}
	}
	if resourceId.SpringName, err = id.PopSegment(springKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'applicationAccelerators' segment
	applicationAcceleratorsKey := "applicationAccelerators"
	for key := range id.Path {
		if strings.EqualFold(key, applicationAcceleratorsKey) {
			applicationAcceleratorsKey = key
			break
		}
	}
	if resourceId.ApplicationAcceleratorName, err = id.PopSegment(applicationAcceleratorsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
