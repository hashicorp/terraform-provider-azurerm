// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorOriginId struct {
	SubscriptionId  string
	ResourceGroup   string
	ProfileName     string
	OriginGroupName string
	OriginName      string
}

func NewFrontDoorOriginID(subscriptionId, resourceGroup, profileName, originGroupName, originName string) FrontDoorOriginId {
	return FrontDoorOriginId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ProfileName:     profileName,
		OriginGroupName: originGroupName,
		OriginName:      originName,
	}
}

func (id FrontDoorOriginId) String() string {
	segments := []string{
		fmt.Sprintf("Origin Name %q", id.OriginName),
		fmt.Sprintf("Origin Group Name %q", id.OriginGroupName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Origin", segmentsStr)
}

func (id FrontDoorOriginId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/originGroups/%s/origins/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
}

// FrontDoorOriginID parses a FrontDoorOrigin ID into an FrontDoorOriginId struct
func FrontDoorOriginID(input string) (*FrontDoorOriginId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FrontDoorOrigin ID: %+v", input, err)
	}

	resourceId := FrontDoorOriginId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}
	if resourceId.OriginGroupName, err = id.PopSegment("originGroups"); err != nil {
		return nil, err
	}
	if resourceId.OriginName, err = id.PopSegment("origins"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontDoorOriginIDInsensitively parses an FrontDoorOrigin ID into an FrontDoorOriginId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontDoorOriginID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontDoorOriginIDInsensitively(input string) (*FrontDoorOriginId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorOriginId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'profiles' segment
	profilesKey := "profiles"
	for key := range id.Path {
		if strings.EqualFold(key, profilesKey) {
			profilesKey = key
			break
		}
	}
	if resourceId.ProfileName, err = id.PopSegment(profilesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'originGroups' segment
	originGroupsKey := "originGroups"
	for key := range id.Path {
		if strings.EqualFold(key, originGroupsKey) {
			originGroupsKey = key
			break
		}
	}
	if resourceId.OriginGroupName, err = id.PopSegment(originGroupsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'origins' segment
	originsKey := "origins"
	for key := range id.Path {
		if strings.EqualFold(key, originsKey) {
			originsKey = key
			break
		}
	}
	if resourceId.OriginName, err = id.PopSegment(originsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
