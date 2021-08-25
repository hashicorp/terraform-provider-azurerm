package eventhubsclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceGroupId struct {
	SubscriptionId string
	ResourceGroup  string
}

func NewResourceGroupID(subscriptionId, resourceGroup string) ResourceGroupId {
	return ResourceGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
	}
}

func (id ResourceGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Group", segmentsStr)
}

func (id ResourceGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup)
}

// ResourceGroupID parses a ResourceGroup ID into an ResourceGroupId struct
func ResourceGroupID(input string) (*ResourceGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ResourceGroupIDInsensitively parses an ResourceGroup ID into an ResourceGroupId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ResourceGroupID method should be used instead for validation etc.
func ResourceGroupIDInsensitively(input string) (*ResourceGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
