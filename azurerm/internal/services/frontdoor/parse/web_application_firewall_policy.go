package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type WebApplicationFirewallPolicyId struct {
	SubscriptionId                            string
	ResourceGroup                             string
	FrontDoorWebApplicationFirewallPolicyName string
}

func NewWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, frontDoorWebApplicationFirewallPolicyName string) WebApplicationFirewallPolicyId {
	return WebApplicationFirewallPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FrontDoorWebApplicationFirewallPolicyName: frontDoorWebApplicationFirewallPolicyName,
	}
}

func (id WebApplicationFirewallPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Front Door Web Application Firewall Policy Name %q", id.FrontDoorWebApplicationFirewallPolicyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Web Application Firewall Policy", segmentsStr)
}

func (id WebApplicationFirewallPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
}

// WebApplicationFirewallPolicyID parses a WebApplicationFirewallPolicy ID into an WebApplicationFirewallPolicyId struct
func WebApplicationFirewallPolicyID(input string) (*WebApplicationFirewallPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebApplicationFirewallPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FrontDoorWebApplicationFirewallPolicyName, err = id.PopSegment("frontDoorWebApplicationFirewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// WebApplicationFirewallPolicyIDInsensitively parses an WebApplicationFirewallPolicy ID into an WebApplicationFirewallPolicyId struct, insensitively
// This should only be used to parse an ID for rewriting, the WebApplicationFirewallPolicyID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func WebApplicationFirewallPolicyIDInsensitively(input string) (*WebApplicationFirewallPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebApplicationFirewallPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'frontDoorWebApplicationFirewallPolicies' segment
	frontDoorWebApplicationFirewallPoliciesKey := "frontDoorWebApplicationFirewallPolicies"
	for key := range id.Path {
		if strings.EqualFold(key, frontDoorWebApplicationFirewallPoliciesKey) {
			frontDoorWebApplicationFirewallPoliciesKey = key
			break
		}
	}
	if resourceId.FrontDoorWebApplicationFirewallPolicyName, err = id.PopSegment(frontDoorWebApplicationFirewallPoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
