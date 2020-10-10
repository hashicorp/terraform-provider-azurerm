package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceFabricMeshApplicationId struct {
	ResourceGroup string
	Name          string
}

func ServiceFabricMeshApplicationID(input string) (*ServiceFabricMeshApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Service Fabric Mesh Application ID %q: %+v", input, err)
	}

	cluster := ServiceFabricMeshApplicationId{
		ResourceGroup: id.ResourceGroup,
	}

	if cluster.Name, err = id.PopSegment("applications"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cluster, nil
}
