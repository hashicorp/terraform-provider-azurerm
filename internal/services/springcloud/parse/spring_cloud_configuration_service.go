// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudConfigurationServiceId struct {
	SubscriptionId           string
	ResourceGroup            string
	SpringName               string
	ConfigurationServiceName string
}

func NewSpringCloudConfigurationServiceID(subscriptionId, resourceGroup, springName, configurationServiceName string) SpringCloudConfigurationServiceId {
	return SpringCloudConfigurationServiceId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		SpringName:               springName,
		ConfigurationServiceName: configurationServiceName,
	}
}

func (id SpringCloudConfigurationServiceId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Service Name %q", id.ConfigurationServiceName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Configuration Service", segmentsStr)
}

func (id SpringCloudConfigurationServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/configurationServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ConfigurationServiceName)
}

// SpringCloudConfigurationServiceID parses a SpringCloudConfigurationService ID into an SpringCloudConfigurationServiceId struct
func SpringCloudConfigurationServiceID(input string) (*SpringCloudConfigurationServiceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudConfigurationService ID: %+v", input, err)
	}

	resourceId := SpringCloudConfigurationServiceId{
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
	if resourceId.ConfigurationServiceName, err = id.PopSegment("configurationServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudConfigurationServiceIDInsensitively parses an SpringCloudConfigurationService ID into an SpringCloudConfigurationServiceId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudConfigurationServiceID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudConfigurationServiceIDInsensitively(input string) (*SpringCloudConfigurationServiceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudConfigurationServiceId{
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

	// find the correct casing for the 'configurationServices' segment
	configurationServicesKey := "configurationServices"
	for key := range id.Path {
		if strings.EqualFold(key, configurationServicesKey) {
			configurationServicesKey = key
			break
		}
	}
	if resourceId.ConfigurationServiceName, err = id.PopSegment(configurationServicesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
