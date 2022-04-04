package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorSecurityPolicyId struct {
	SubscriptionId     string
	ResourceGroup      string
	ProfileName        string
	SecurityPolicyName string
}

func NewFrontdoorSecurityPolicyID(subscriptionId, resourceGroup, profileName, securityPolicyName string) FrontdoorSecurityPolicyId {
	return FrontdoorSecurityPolicyId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		ProfileName:        profileName,
		SecurityPolicyName: securityPolicyName,
	}
}

func (id FrontdoorSecurityPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Security Policy Name %q", id.SecurityPolicyName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Security Policy", segmentsStr)
}

func (id FrontdoorSecurityPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/securityPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
}

// FrontdoorSecurityPolicyID parses a FrontdoorSecurityPolicy ID into an FrontdoorSecurityPolicyId struct
func FrontdoorSecurityPolicyID(input string) (*FrontdoorSecurityPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorSecurityPolicyId{
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
	if resourceId.SecurityPolicyName, err = id.PopSegment("securityPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontdoorSecurityPolicyIDInsensitively parses an FrontdoorSecurityPolicy ID into an FrontdoorSecurityPolicyId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontdoorSecurityPolicyID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontdoorSecurityPolicyIDInsensitively(input string) (*FrontdoorSecurityPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorSecurityPolicyId{
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

	// find the correct casing for the 'securityPolicies' segment
	securityPoliciesKey := "securityPolicies"
	for key := range id.Path {
		if strings.EqualFold(key, securityPoliciesKey) {
			securityPoliciesKey = key
			break
		}
	}
	if resourceId.SecurityPolicyName, err = id.PopSegment(securityPoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
