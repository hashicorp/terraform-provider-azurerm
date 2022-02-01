package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorProfilePolicyId struct {
	SubscriptionId                      string
	ResourceGroup                       string
	CdnWebApplicationFirewallPolicyName string
}

func NewFrontdoorProfilePolicyID(subscriptionId, resourceGroup, cdnWebApplicationFirewallPolicyName string) FrontdoorProfilePolicyId {
	return FrontdoorProfilePolicyId{
		SubscriptionId:                      subscriptionId,
		ResourceGroup:                       resourceGroup,
		CdnWebApplicationFirewallPolicyName: cdnWebApplicationFirewallPolicyName,
	}
}

func (id FrontdoorProfilePolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Cdn Web Application Firewall Policy Name %q", id.CdnWebApplicationFirewallPolicyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Profile Policy", segmentsStr)
}

func (id FrontdoorProfilePolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CdnWebApplicationFirewallPolicyName)
}

// FrontdoorProfilePolicyID parses a FrontdoorProfilePolicy ID into an FrontdoorProfilePolicyId struct
func FrontdoorProfilePolicyID(input string) (*FrontdoorProfilePolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorProfilePolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CdnWebApplicationFirewallPolicyName, err = id.PopSegment("cdnWebApplicationFirewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
