// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudServiceRegistryId struct {
	SubscriptionId      string
	ResourceGroup       string
	SpringName          string
	ServiceRegistryName string
}

func NewSpringCloudServiceRegistryID(subscriptionId, resourceGroup, springName, serviceRegistryName string) SpringCloudServiceRegistryId {
	return SpringCloudServiceRegistryId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		SpringName:          springName,
		ServiceRegistryName: serviceRegistryName,
	}
}

func (id SpringCloudServiceRegistryId) String() string {
	segments := []string{
		fmt.Sprintf("Service Registry Name %q", id.ServiceRegistryName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Service Registry", segmentsStr)
}

func (id SpringCloudServiceRegistryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/serviceRegistries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ServiceRegistryName)
}

// SpringCloudServiceRegistryID parses a SpringCloudServiceRegistry ID into an SpringCloudServiceRegistryId struct
func SpringCloudServiceRegistryID(input string) (*SpringCloudServiceRegistryId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudServiceRegistry ID: %+v", input, err)
	}

	resourceId := SpringCloudServiceRegistryId{
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
	if resourceId.ServiceRegistryName, err = id.PopSegment("serviceRegistries"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudServiceRegistryIDInsensitively parses an SpringCloudServiceRegistry ID into an SpringCloudServiceRegistryId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudServiceRegistryID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudServiceRegistryIDInsensitively(input string) (*SpringCloudServiceRegistryId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudServiceRegistryId{
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

	// find the correct casing for the 'serviceRegistries' segment
	serviceRegistriesKey := "serviceRegistries"
	for key := range id.Path {
		if strings.EqualFold(key, serviceRegistriesKey) {
			serviceRegistriesKey = key
			break
		}
	}
	if resourceId.ServiceRegistryName, err = id.PopSegment(serviceRegistriesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
