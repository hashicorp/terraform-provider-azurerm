package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type JitNetworkAccessPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	LocationName   string
	Name           string
}

func NewJitNetworkAccessPolicyID(subscriptionId, resourceGroup, locationName, name string) JitNetworkAccessPolicyId {
	return JitNetworkAccessPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LocationName:   locationName,
		Name:           name,
	}
}

func (id JitNetworkAccessPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Location Name %q", id.LocationName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Jit Network Access Policy", segmentsStr)
}

func (id JitNetworkAccessPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Security/locations/%s/jitNetworkAccessPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LocationName, id.Name)
}

// JitNetworkAccessPolicyID parses a JitNetworkAccessPolicy ID into an JitNetworkAccessPolicyId struct
func JitNetworkAccessPolicyID(input string) (*JitNetworkAccessPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := JitNetworkAccessPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("jitNetworkAccessPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
