package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type NamespaceId struct {
	Name          string
	ResourceGroup string
}

func NamespaceID(input string) (*NamespaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	rule := NamespaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if rule.Name, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &rule, nil
}
