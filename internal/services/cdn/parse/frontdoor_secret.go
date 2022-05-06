package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorSecretId struct {
	SubscriptionId string
	ResourceGroup  string
	ProfileName    string
	SecretName     string
}

func NewFrontdoorSecretID(subscriptionId, resourceGroup, profileName, secretName string) FrontdoorSecretId {
	return FrontdoorSecretId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ProfileName:    profileName,
		SecretName:     secretName,
	}
}

func (id FrontdoorSecretId) String() string {
	segments := []string{
		fmt.Sprintf("Secret Name %q", id.SecretName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Secret", segmentsStr)
}

func (id FrontdoorSecretId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/secrets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.SecretName)
}

// FrontdoorSecretID parses a FrontdoorSecret ID into an FrontdoorSecretId struct
func FrontdoorSecretID(input string) (*FrontdoorSecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorSecretId{
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
	if resourceId.SecretName, err = id.PopSegment("secrets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontdoorSecretIDInsensitively parses an FrontdoorSecret ID into an FrontdoorSecretId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontdoorSecretID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontdoorSecretIDInsensitively(input string) (*FrontdoorSecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorSecretId{
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
