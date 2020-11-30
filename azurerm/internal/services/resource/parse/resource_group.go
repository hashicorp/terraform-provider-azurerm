package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ResourceGroupId struct {
	SubscriptionId string
	ResourceGroup  string
}

func (id ResourceGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup)
}

func NewResourceGroupID(subscriptionId, name string) ResourceGroupId {
	return ResourceGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  name,
	}
}

func ResourceGroupID(input string) (*ResourceGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Resource Group ID %q: %+v", input, err)
	}

	group := ResourceGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if group.ResourceGroup == "" {
		return nil, fmt.Errorf("ID contained no `resourceGroups` segment!")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &group, nil
}
