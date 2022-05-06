package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorCustomDomainSecretId struct {
	SubscriptionId   string
	ResourceGroup    string
	ProfileName      string
	CustomDomainName string
	SecretName       string
}

func NewFrontdoorCustomDomainSecretID(subscriptionId, resourceGroup, profileName, customDomainName, secretName string) FrontdoorCustomDomainSecretId {
	return FrontdoorCustomDomainSecretId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		ProfileName:      profileName,
		CustomDomainName: customDomainName,
		SecretName:       secretName,
	}
}

func (id FrontdoorCustomDomainSecretId) String() string {
	segments := []string{
		fmt.Sprintf("Secret Name %q", id.SecretName),
		fmt.Sprintf("Custom Domain Name %q", id.CustomDomainName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Custom Domain Secret", segmentsStr)
}

func (id FrontdoorCustomDomainSecretId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/customDomains/%s/secrets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.CustomDomainName, id.SecretName)
}

// FrontdoorCustomDomainSecretID parses a FrontdoorCustomDomainSecret ID into an FrontdoorCustomDomainSecretId struct
func FrontdoorCustomDomainSecretID(input string) (*FrontdoorCustomDomainSecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorCustomDomainSecretId{
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
	if resourceId.SecretName, err = id.PopSegment("secrets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontdoorCustomDomainSecretIDInsensitively parses an FrontdoorCustomDomainSecret ID into an FrontdoorCustomDomainSecretId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontdoorCustomDomainSecretID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontdoorCustomDomainSecretIDInsensitively(input string) (*FrontdoorCustomDomainSecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorCustomDomainSecretId{
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

	// find the correct casing for the 'secrets' segment
	secretsKey := "secrets"
	for key := range id.Path {
		if strings.EqualFold(key, secretsKey) {
			secretsKey = key
			break
		}
	}
	if resourceId.SecretName, err = id.PopSegment(secretsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
