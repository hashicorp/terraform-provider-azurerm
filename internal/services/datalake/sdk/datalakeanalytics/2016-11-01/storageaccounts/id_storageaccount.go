package storageaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StorageAccountId struct {
	SubscriptionId string
	ResourceGroup  string
	AccountName    string
	Name           string
}

func NewStorageAccountID(subscriptionId, resourceGroup, accountName, name string) StorageAccountId {
	return StorageAccountId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		AccountName:    accountName,
		Name:           name,
	}
}

func (id StorageAccountId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Account", segmentsStr)
}

func (id StorageAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeAnalytics/accounts/%s/storageAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.Name)
}

// ParseStorageAccountID parses a StorageAccount ID into an StorageAccountId struct
func ParseStorageAccountID(input string) (*StorageAccountId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageAccountId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AccountName, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseStorageAccountIDInsensitively parses an StorageAccount ID into an StorageAccountId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseStorageAccountID method should be used instead for validation etc.
func ParseStorageAccountIDInsensitively(input string) (*StorageAccountId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StorageAccountId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'accounts' segment
	accountsKey := "accounts"
	for key := range id.Path {
		if strings.EqualFold(key, accountsKey) {
			accountsKey = key
			break
		}
	}
	if resourceId.AccountName, err = id.PopSegment(accountsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'storageAccounts' segment
	storageAccountsKey := "storageAccounts"
	for key := range id.Path {
		if strings.EqualFold(key, storageAccountsKey) {
			storageAccountsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(storageAccountsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
