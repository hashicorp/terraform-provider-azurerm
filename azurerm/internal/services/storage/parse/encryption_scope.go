package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EncryptionScopeId struct {
	Name               string
	StorageAccountName string
	ResourceGroup      string
	SubscriptionId     string
}

// the subscriptionId isn't used here, this is just to comply with the interface for now..
func (id EncryptionScopeId) ID(_ string) string {
	fmtString := "%s/encryptionScopes/%s"
	accountId := NewAccountId(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName).ID("")
	return fmt.Sprintf(fmtString, accountId, id.Name)
}

func NewEncryptionScopeId(resourceGroup, storageAccount, name string) EncryptionScopeId {
	return EncryptionScopeId{
		Name:               name,
		StorageAccountName: storageAccount,
		ResourceGroup:      resourceGroup,
		SubscriptionId:     "TODO",
	}
}

func EncryptionScopeID(input string) (*EncryptionScopeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	es := EncryptionScopeId{
		ResourceGroup:  id.ResourceGroup,
		SubscriptionId: id.SubscriptionID,
	}

	if es.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
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
