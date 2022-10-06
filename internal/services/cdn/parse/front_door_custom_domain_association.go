package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorCustomDomainAssociationId struct {
	SubscriptionId   string
	ResourceGroup    string
	ProfileName      string
	CustomDomainName string
	AssociationName  string
}

func NewFrontDoorCustomDomainAssociationID(subscriptionId, resourceGroup, profileName, customDomainName, associationName string) FrontDoorCustomDomainAssociationId {
	return FrontDoorCustomDomainAssociationId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		ProfileName:      profileName,
		CustomDomainName: customDomainName,
		AssociationName:  associationName,
	}
}

func (id FrontDoorCustomDomainAssociationId) String() string {
	segments := []string{
		fmt.Sprintf("Association Name %q", id.AssociationName),
		fmt.Sprintf("Custom Domain Name %q", id.CustomDomainName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Custom Domain Association", segmentsStr)
}

func (id FrontDoorCustomDomainAssociationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/customDomains/%s/associations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.CustomDomainName, id.AssociationName)
}

// FrontDoorCustomDomainAssociationID parses a FrontDoorCustomDomainAssociation ID into an FrontDoorCustomDomainAssociationId struct
func FrontDoorCustomDomainAssociationID(input string) (*FrontDoorCustomDomainAssociationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
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
	if resourceId.CustomDomainName, err = id.PopSegment("customDomains"); err != nil {
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

// FrontDoorCustomDomainAssociationIDInsensitively parses an FrontDoorCustomDomainAssociation ID into an FrontDoorCustomDomainAssociationId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontDoorCustomDomainAssociationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontDoorCustomDomainAssociationIDInsensitively(input string) (*FrontDoorCustomDomainAssociationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
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

	// find the correct casing for the 'customDomains' segment
	customDomainsKey := "customDomains"
	for key := range id.Path {
		if strings.EqualFold(key, customDomainsKey) {
			customDomainsKey = key
			break
		}
	}
	if resourceId.CustomDomainName, err = id.PopSegment(customDomainsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'associations' segment
	associationsKey := "associations"
	for key := range id.Path {
		if strings.EqualFold(key, associationsKey) {
			associationsKey = key
			break
		}
	}
	if resourceId.AssociationName, err = id.PopSegment(associationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
