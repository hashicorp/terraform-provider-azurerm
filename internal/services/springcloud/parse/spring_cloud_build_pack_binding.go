// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudBuildPackBindingId struct {
	SubscriptionId       string
	ResourceGroup        string
	SpringName           string
	BuildServiceName     string
	BuilderName          string
	BuildPackBindingName string
}

func NewSpringCloudBuildPackBindingID(subscriptionId, resourceGroup, springName, buildServiceName, builderName, buildPackBindingName string) SpringCloudBuildPackBindingId {
	return SpringCloudBuildPackBindingId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		SpringName:           springName,
		BuildServiceName:     buildServiceName,
		BuilderName:          builderName,
		BuildPackBindingName: buildPackBindingName,
	}
}

func (id SpringCloudBuildPackBindingId) String() string {
	segments := []string{
		fmt.Sprintf("Build Pack Binding Name %q", id.BuildPackBindingName),
		fmt.Sprintf("Builder Name %q", id.BuilderName),
		fmt.Sprintf("Build Service Name %q", id.BuildServiceName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Build Pack Binding", segmentsStr)
}

func (id SpringCloudBuildPackBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/buildServices/%s/builders/%s/buildPackBindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName)
}

// SpringCloudBuildPackBindingID parses a SpringCloudBuildPackBinding ID into an SpringCloudBuildPackBindingId struct
func SpringCloudBuildPackBindingID(input string) (*SpringCloudBuildPackBindingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudBuildPackBinding ID: %+v", input, err)
	}

	resourceId := SpringCloudBuildPackBindingId{
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
	if resourceId.BuildPackBindingName, err = id.PopSegment("buildPackBindings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudBuildPackBindingIDInsensitively parses an SpringCloudBuildPackBinding ID into an SpringCloudBuildPackBindingId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudBuildPackBindingID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudBuildPackBindingIDInsensitively(input string) (*SpringCloudBuildPackBindingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudBuildPackBindingId{
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

	// find the correct casing for the 'buildPackBindings' segment
	buildPackBindingsKey := "buildPackBindings"
	for key := range id.Path {
		if strings.EqualFold(key, buildPackBindingsKey) {
			buildPackBindingsKey = key
			break
		}
	}
	if resourceId.BuildPackBindingName, err = id.PopSegment(buildPackBindingsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
