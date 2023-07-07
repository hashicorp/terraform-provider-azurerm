// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudContainerRegistryId struct {
	SubscriptionId        string
	ResourceGroup         string
	SpringName            string
	ContainerRegistryName string
}

func NewSpringCloudContainerRegistryID(subscriptionId, resourceGroup, springName, containerRegistryName string) SpringCloudContainerRegistryId {
	return SpringCloudContainerRegistryId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		SpringName:            springName,
		ContainerRegistryName: containerRegistryName,
	}
}

func (id SpringCloudContainerRegistryId) String() string {
	segments := []string{
		fmt.Sprintf("Container Registry Name %q", id.ContainerRegistryName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Container Registry", segmentsStr)
}

func (id SpringCloudContainerRegistryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/containerRegistries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ContainerRegistryName)
}

// SpringCloudContainerRegistryID parses a SpringCloudContainerRegistry ID into an SpringCloudContainerRegistryId struct
func SpringCloudContainerRegistryID(input string) (*SpringCloudContainerRegistryId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudContainerRegistry ID: %+v", input, err)
	}

	resourceId := SpringCloudContainerRegistryId{
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
	if resourceId.ContainerRegistryName, err = id.PopSegment("containerRegistries"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudContainerRegistryIDInsensitively parses an SpringCloudContainerRegistry ID into an SpringCloudContainerRegistryId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudContainerRegistryID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudContainerRegistryIDInsensitively(input string) (*SpringCloudContainerRegistryId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudContainerRegistryId{
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

	// find the correct casing for the 'containerRegistries' segment
	containerRegistriesKey := "containerRegistries"
	for key := range id.Path {
		if strings.EqualFold(key, containerRegistriesKey) {
			containerRegistriesKey = key
			break
		}
	}
	if resourceId.ContainerRegistryName, err = id.PopSegment(containerRegistriesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
