package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorPolicyId struct {
	SubscriptionId                            string
	ResourceGroup                             string
	FrontDoorWebApplicationFirewallPolicyName string
}

func NewFrontdoorPolicyID(subscriptionId, resourceGroup, frontDoorWebApplicationFirewallPolicyName string) FrontdoorPolicyId {
	return FrontdoorPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FrontDoorWebApplicationFirewallPolicyName: frontDoorWebApplicationFirewallPolicyName,
	}
}

func (id FrontdoorPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Front Door Web Application Firewall Policy Name %q", id.FrontDoorWebApplicationFirewallPolicyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Policy", segmentsStr)
}

func (id FrontdoorPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/FrontDoorWebApplicationFirewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
}

// FrontdoorPolicyID parses a FrontdoorPolicy ID into an FrontdoorPolicyId struct
func FrontdoorPolicyID(input string) (*FrontdoorPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FrontDoorWebApplicationFirewallPolicyName, err = id.PopSegment("FrontDoorWebApplicationFirewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontdoorPolicyIDInsensitively parses an FrontdoorPolicy ID into an FrontdoorPolicyId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontdoorPolicyID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontdoorPolicyIDInsensitively(input string) (*FrontdoorPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'FrontDoorWebApplicationFirewallPolicies' segment
	FrontDoorWebApplicationFirewallPoliciesKey := "FrontDoorWebApplicationFirewallPolicies"
	for key := range id.Path {
		if strings.EqualFold(key, FrontDoorWebApplicationFirewallPoliciesKey) {
			FrontDoorWebApplicationFirewallPoliciesKey = key
			break
		}
	}
	if resourceId.FrontDoorWebApplicationFirewallPolicyName, err = id.PopSegment(FrontDoorWebApplicationFirewallPoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
