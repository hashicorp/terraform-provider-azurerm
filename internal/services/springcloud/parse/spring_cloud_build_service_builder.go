// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudBuildServiceBuilderId struct {
	SubscriptionId   string
	ResourceGroup    string
	SpringName       string
	BuildServiceName string
	BuilderName      string
}

func NewSpringCloudBuildServiceBuilderID(subscriptionId, resourceGroup, springName, buildServiceName, builderName string) SpringCloudBuildServiceBuilderId {
	return SpringCloudBuildServiceBuilderId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		SpringName:       springName,
		BuildServiceName: buildServiceName,
		BuilderName:      builderName,
	}
}

func (id SpringCloudBuildServiceBuilderId) String() string {
	segments := []string{
		fmt.Sprintf("Builder Name %q", id.BuilderName),
		fmt.Sprintf("Build Service Name %q", id.BuildServiceName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Build Service Builder", segmentsStr)
}

func (id SpringCloudBuildServiceBuilderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/buildServices/%s/builders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName)
}

// SpringCloudBuildServiceBuilderID parses a SpringCloudBuildServiceBuilder ID into an SpringCloudBuildServiceBuilderId struct
func SpringCloudBuildServiceBuilderID(input string) (*SpringCloudBuildServiceBuilderId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudBuildServiceBuilder ID: %+v", input, err)
	}

	resourceId := SpringCloudBuildServiceBuilderId{
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
	if resourceId.BuildServiceName, err = id.PopSegment("buildServices"); err != nil {
		return nil, err
	}
	if resourceId.BuilderName, err = id.PopSegment("builders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudBuildServiceBuilderIDInsensitively parses an SpringCloudBuildServiceBuilder ID into an SpringCloudBuildServiceBuilderId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudBuildServiceBuilderID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudBuildServiceBuilderIDInsensitively(input string) (*SpringCloudBuildServiceBuilderId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudBuildServiceBuilderId{
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

	// find the correct casing for the 'buildServices' segment
	buildServicesKey := "buildServices"
	for key := range id.Path {
		if strings.EqualFold(key, buildServicesKey) {
			buildServicesKey = key
			break
		}
	}
	if resourceId.BuildServiceName, err = id.PopSegment(buildServicesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'builders' segment
	buildersKey := "builders"
	for key := range id.Path {
		if strings.EqualFold(key, buildersKey) {
			buildersKey = key
			break
		}
	}
	if resourceId.BuilderName, err = id.PopSegment(buildersKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
