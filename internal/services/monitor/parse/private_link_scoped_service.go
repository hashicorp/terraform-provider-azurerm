package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type PrivateLinkScopedServiceId struct {
	SubscriptionId       string
	ResourceGroup        string
	PrivateLinkScopeName string
	ScopedResourceName   string
}

func NewPrivateLinkScopedServiceID(subscriptionId, resourceGroup, privateLinkScopeName, scopedResourceName string) PrivateLinkScopedServiceId {
	return PrivateLinkScopedServiceId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		PrivateLinkScopeName: privateLinkScopeName,
		ScopedResourceName:   scopedResourceName,
	}
}

func (id PrivateLinkScopedServiceId) String() string {
	segments := []string{
		fmt.Sprintf("Scoped Resource Name %q", id.ScopedResourceName),
		fmt.Sprintf("Private Link Scope Name %q", id.PrivateLinkScopeName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Link Scoped Service", segmentsStr)
}

func (id PrivateLinkScopedServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/privateLinkScopes/%s/scopedResources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateLinkScopeName, id.ScopedResourceName)
}

// PrivateLinkScopedServiceID parses a PrivateLinkScopedService ID into an PrivateLinkScopedServiceId struct
func PrivateLinkScopedServiceID(input string) (*PrivateLinkScopedServiceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateLinkScopedServiceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PrivateLinkScopeName, err = id.PopSegment("privateLinkScopes"); err != nil {
		return nil, err
	}
	if resourceId.ScopedResourceName, err = id.PopSegment("scopedResources"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
