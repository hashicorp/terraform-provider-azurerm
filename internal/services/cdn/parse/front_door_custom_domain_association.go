// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorCustomDomainAssociationId struct {
	SubscriptionId  string
	ResourceGroup   string
	ProfileName     string
	AssociationName string
}

func NewFrontDoorCustomDomainAssociationID(subscriptionId, resourceGroup, profileName, associationName string) FrontDoorCustomDomainAssociationId {
	return FrontDoorCustomDomainAssociationId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ProfileName:     profileName,
		AssociationName: associationName,
	}
}

func (id FrontDoorCustomDomainAssociationId) String() string {
	segments := []string{
		fmt.Sprintf("Association Name %q", id.AssociationName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Custom Domain Association", segmentsStr)
}

func (id FrontDoorCustomDomainAssociationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/associations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AssociationName)
}

// FrontDoorCustomDomainAssociationID parses a FrontDoorCustomDomainAssociation ID into an FrontDoorCustomDomainAssociationId struct
func FrontDoorCustomDomainAssociationID(input string) (*FrontDoorCustomDomainAssociationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FrontDoorCustomDomainAssociation ID: %+v", input, err)
	}

	resourceId := FrontDoorCustomDomainAssociationId{
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
	if resourceId.AssociationName, err = id.PopSegment("associations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
