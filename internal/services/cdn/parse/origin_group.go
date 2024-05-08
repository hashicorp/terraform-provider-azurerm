// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type OriginGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	ProfileName    string
	EndpointName   string
	Name           string
}

func NewOriginGroupID(subscriptionId, resourceGroup, profileName, endpointName, name string) OriginGroupId {
	return OriginGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ProfileName:    profileName,
		EndpointName:   endpointName,
		Name:           name,
	}
}

func (id OriginGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Endpoint Name %q", id.EndpointName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Origin Group", segmentsStr)
}

func (id OriginGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/endpoints/%s/originGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
}

// OriginGroupID parses a OriginGroup ID into an OriginGroupId struct
func OriginGroupID(input string) (*OriginGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an OriginGroup ID: %+v", input, err)
	}

	resourceId := OriginGroupId{
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
	if resourceId.EndpointName, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("originGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
