package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type StorageEncryptionScopeId struct {
	Name           string
	StorageAccName string
	ResourceGroup  string
}

func StorageEncryptionScopeID(input string) (*StorageEncryptionScopeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	es := StorageEncryptionScopeId{
		ResourceGroup: id.ResourceGroup,
	}

	if es.StorageAccName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}

	if es.Name, err = id.PopSegment("encryptionScopes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &es, nil
}
