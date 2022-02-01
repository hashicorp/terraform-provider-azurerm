package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorProfileSecurityPolicyId struct {
	SubscriptionId     string
	ResourceGroup      string
	ProfileName        string
	SecurityPolicyName string
}

func NewFrontdoorProfileSecurityPolicyID(subscriptionId, resourceGroup, profileName, securityPolicyName string) FrontdoorProfileSecurityPolicyId {
	return FrontdoorProfileSecurityPolicyId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		ProfileName:        profileName,
		SecurityPolicyName: securityPolicyName,
	}
}

func (id FrontdoorProfileSecurityPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Security Policy Name %q", id.SecurityPolicyName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Profile Security Policy", segmentsStr)
}

func (id FrontdoorProfileSecurityPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/securityPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.SecurityPolicyName)
}

// FrontdoorProfileSecurityPolicyID parses a FrontdoorProfileSecurityPolicy ID into an FrontdoorProfileSecurityPolicyId struct
func FrontdoorProfileSecurityPolicyID(input string) (*FrontdoorProfileSecurityPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorProfileSecurityPolicyId{
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
