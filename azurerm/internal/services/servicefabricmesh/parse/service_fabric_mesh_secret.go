package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceFabricMeshSecretId struct {
	ResourceGroup string
	Name          string
}

func ServiceFabricMeshSecretID(input string) (*ServiceFabricMeshSecretId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Service Fabric Mesh Secret ID %q: %+v", input, err)
	}

	secret := ServiceFabricMeshSecretId{
		ResourceGroup: id.ResourceGroup,
	}

	if secret.Name, err = id.PopSegment("secrets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &secret, nil
}
