package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorCustomDomainTxtId struct {
	SubscriptionId   string
	ResourceGroup    string
	ProfileName      string
	CustomDomainName string
	DnsTxtName       string
}

func NewFrontdoorCustomDomainTxtID(subscriptionId, resourceGroup, profileName, customDomainName, dnsTxtName string) FrontdoorCustomDomainTxtId {
	return FrontdoorCustomDomainTxtId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		ProfileName:      profileName,
		CustomDomainName: customDomainName,
		DnsTxtName:       dnsTxtName,
	}
}

func (id FrontdoorCustomDomainTxtId) String() string {
	segments := []string{
		fmt.Sprintf("Dns Txt Name %q", id.DnsTxtName),
		fmt.Sprintf("Custom Domain Name %q", id.CustomDomainName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Custom Domain Txt", segmentsStr)
}

func (id FrontdoorCustomDomainTxtId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/customDomains/%s/dnsTxt/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.CustomDomainName, id.DnsTxtName)
}

// FrontdoorCustomDomainTxtID parses a FrontdoorCustomDomainTxt ID into an FrontdoorCustomDomainTxtId struct
func FrontdoorCustomDomainTxtID(input string) (*FrontdoorCustomDomainTxtId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorCustomDomainTxtId{
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
	if resourceId.DnsTxtName, err = id.PopSegment("dnsTxt"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontdoorCustomDomainTxtIDInsensitively parses an FrontdoorCustomDomainTxt ID into an FrontdoorCustomDomainTxtId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontdoorCustomDomainTxtID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontdoorCustomDomainTxtIDInsensitively(input string) (*FrontdoorCustomDomainTxtId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorCustomDomainTxtId{
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

	// find the correct casing for the 'dnsTxt' segment
	dnsTxtKey := "dnsTxt"
	for key := range id.Path {
		if strings.EqualFold(key, dnsTxtKey) {
			dnsTxtKey = key
			break
		}
	}
	if resourceId.DnsTxtName, err = id.PopSegment(dnsTxtKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
