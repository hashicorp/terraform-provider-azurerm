package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecretValueId struct {
	ResourceGroup string
	SecretName    string
	ValueName     string
}

func SecretValueID(input string) (*SecretValueId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Service Fabric Mesh Secret ID %q: %+v", input, err)
	}

	value := SecretValueId{
		ResourceGroup: id.ResourceGroup,
	}

	if value.SecretName, err = id.PopSegment("secrets"); err != nil {
		return nil, err
	}

	if value.ValueName, err = id.PopSegment("values"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &value, nil
}
