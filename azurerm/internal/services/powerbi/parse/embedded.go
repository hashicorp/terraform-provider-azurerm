package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EmbeddedId struct {
	SubscriptionId string
	ResourceGroup  string
	CapacityName   string
}

func NewEmbeddedID(subscriptionId, resourceGroup, capacityName string) EmbeddedId {
	return EmbeddedId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		CapacityName:   capacityName,
	}
}

func (id EmbeddedId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Capacity Name %q", id.CapacityName),
	}
	return strings.Join(segments, " / ")
}

func (id EmbeddedId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.PowerBIDedicated/capacities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CapacityName)
}

// EmbeddedID parses a Embedded ID into an EmbeddedId struct
func EmbeddedID(input string) (*EmbeddedId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EmbeddedId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CapacityName, err = id.PopSegment("capacities"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
