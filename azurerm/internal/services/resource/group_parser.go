package resource

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ResourceGroupResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseResourceGroupID(input string) (*ResourceGroupResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Resource Group ID %q: %+v", input, err)
	}

	group := ResourceGroupResourceID{
		Base: *id,
		Name: id.ResourceGroup,
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID was missing the `resourceGroups` element")
	}

	pathWithoutSubs := group.Base.Path
	delete(pathWithoutSubs, "subscriptions")
	if len(pathWithoutSubs) != 0 {
		return nil, fmt.Errorf("ID contained more segments than a Resource ID requires: %q", input)
	}

	return &group, nil
}
