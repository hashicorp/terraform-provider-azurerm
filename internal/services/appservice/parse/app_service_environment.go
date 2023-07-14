// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AppServiceEnvironmentId struct {
	SubscriptionId         string
	ResourceGroup          string
	HostingEnvironmentName string
}

func NewAppServiceEnvironmentID(subscriptionId, resourceGroup, hostingEnvironmentName string) AppServiceEnvironmentId {
	return AppServiceEnvironmentId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		HostingEnvironmentName: hostingEnvironmentName,
	}
}

func (id AppServiceEnvironmentId) String() string {
	segments := []string{
		fmt.Sprintf("Hosting Environment Name %q", id.HostingEnvironmentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "App Service Environment", segmentsStr)
}

func (id AppServiceEnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/hostingEnvironments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.HostingEnvironmentName)
}

// AppServiceEnvironmentID parses a AppServiceEnvironment ID into an AppServiceEnvironmentId struct
func AppServiceEnvironmentID(input string) (*AppServiceEnvironmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AppServiceEnvironment ID: %+v", input, err)
	}

	resourceId := AppServiceEnvironmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.HostingEnvironmentName, err = id.PopSegment("hostingEnvironments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// AppServiceEnvironmentIDInsensitively parses an AppServiceEnvironment ID into an AppServiceEnvironmentId struct, insensitively
// This should only be used to parse an ID for rewriting, the AppServiceEnvironmentID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func AppServiceEnvironmentIDInsensitively(input string) (*AppServiceEnvironmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AppServiceEnvironmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'hostingEnvironments' segment
	hostingEnvironmentsKey := "hostingEnvironments"
	for key := range id.Path {
		if strings.EqualFold(key, hostingEnvironmentsKey) {
			hostingEnvironmentsKey = key
			break
		}
	}
	if resourceId.HostingEnvironmentName, err = id.PopSegment(hostingEnvironmentsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
