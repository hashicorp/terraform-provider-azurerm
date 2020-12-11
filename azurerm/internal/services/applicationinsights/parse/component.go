package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ComponentId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewComponentID(subscriptionId, resourceGroup, name string) ComponentId {
	return ComponentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ComponentId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Component", segmentsStr)
}

func (id ComponentId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/components/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ComponentID parses a Component ID into an ComponentId struct
func ComponentID(input string) (*ComponentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ComponentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("components"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
