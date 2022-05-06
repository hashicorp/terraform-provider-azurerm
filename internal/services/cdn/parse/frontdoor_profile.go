package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorProfileId struct {
	SubscriptionId string
	ResourceGroup  string
	ProfileName    string
}

func NewFrontdoorProfileID(subscriptionId, resourceGroup, profileName string) FrontdoorProfileId {
	return FrontdoorProfileId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ProfileName:    profileName,
	}
}

func (id FrontdoorProfileId) String() string {
	segments := []string{
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Profile", segmentsStr)
}

func (id FrontdoorProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName)
}

// FrontdoorProfileID parses a FrontdoorProfile ID into an FrontdoorProfileId struct
func FrontdoorProfileID(input string) (*FrontdoorProfileId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorProfileId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontdoorProfileIDInsensitively parses an FrontdoorProfile ID into an FrontdoorProfileId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontdoorProfileID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontdoorProfileIDInsensitively(input string) (*FrontdoorProfileId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorProfileId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
