package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CostManagementExportResourceGroupId struct {
	ResourceId string
	Name       string
}

func CostManagementExportResourceGroupID(input string) (*CostManagementExportResourceGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Cost Management Export Resource Group %q: %+v", input, err)
	}

	service := CostManagementExportResourceGroupId{
		ResourceId: fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", id.SubscriptionID, id.ResourceGroup),
	}

	if service.Name, err = id.PopSegment("exports"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
