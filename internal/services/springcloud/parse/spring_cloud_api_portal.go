// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudAPIPortalId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	ApiPortalName  string
}

func NewSpringCloudAPIPortalID(subscriptionId, resourceGroup, springName, apiPortalName string) SpringCloudAPIPortalId {
	return SpringCloudAPIPortalId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		ApiPortalName:  apiPortalName,
	}
}

func (id SpringCloudAPIPortalId) String() string {
	segments := []string{
		fmt.Sprintf("Api Portal Name %q", id.ApiPortalName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring CloudAPI Portal", segmentsStr)
}

func (id SpringCloudAPIPortalId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apiPortals/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ApiPortalName)
}

// SpringCloudAPIPortalID parses a SpringCloudAPIPortal ID into an SpringCloudAPIPortalId struct
func SpringCloudAPIPortalID(input string) (*SpringCloudAPIPortalId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudAPIPortal ID: %+v", input, err)
	}

	resourceId := SpringCloudAPIPortalId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("spring"); err != nil {
		return nil, err
	}
	if resourceId.ApiPortalName, err = id.PopSegment("apiPortals"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudAPIPortalIDInsensitively parses an SpringCloudAPIPortal ID into an SpringCloudAPIPortalId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudAPIPortalID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudAPIPortalIDInsensitively(input string) (*SpringCloudAPIPortalId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudAPIPortalId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
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

	// find the correct casing for the 'apiPortals' segment
	apiPortalsKey := "apiPortals"
	for key := range id.Path {
		if strings.EqualFold(key, apiPortalsKey) {
			apiPortalsKey = key
			break
		}
	}
	if resourceId.ApiPortalName, err = id.PopSegment(apiPortalsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
