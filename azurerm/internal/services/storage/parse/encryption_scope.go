package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EncryptionScopeId struct {
	Name           string
	AccountName    string
	ResourceGroup  string
	SubscriptionId string
}

// the subscriptionId isn't used here, this is just to comply with the interface for now..
func (id EncryptionScopeId) ID(_ string) string {
	fmtString := "%s/encryptionScopes/%s"
	accountId := NewAccountId(id.SubscriptionId, id.ResourceGroup, id.Name).ID("")
	return fmt.Sprintf(fmtString, accountId, id.Name)
}

func NewEncryptionScopeId(storageAccount AccountId, name string) EncryptionScopeId {
	return EncryptionScopeId{
		Name:           name,
		AccountName:    storageAccount.Name,
		ResourceGroup:  storageAccount.ResourceGroup,
		SubscriptionId: storageAccount.SubscriptionId,
	}
}

func EncryptionScopeID(input string) (*EncryptionScopeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	es := EncryptionScopeId{
		ResourceGroup: id.ResourceGroup,
	}

	if es.AccountName, err = id.PopSegment("storageAccounts"); err != nil {
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
