package storageaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ContainerId struct {
	SubscriptionId     string
	ResourceGroup      string
	AccountName        string
	StorageAccountName string
	Name               string
}

func NewContainerID(subscriptionId, resourceGroup, accountName, storageAccountName, name string) ContainerId {
	return ContainerId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		AccountName:        accountName,
		StorageAccountName: storageAccountName,
		Name:               name,
	}
}

func (id ContainerId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container", segmentsStr)
}

func (id ContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeAnalytics/accounts/%s/storageAccounts/%s/containers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.StorageAccountName, id.Name)
}

// ParseContainerID parses a Container ID into an ContainerId struct
func ParseContainerID(input string) (*ContainerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ContainerId{
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
	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseContainerIDInsensitively parses an Container ID into an ContainerId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseContainerID method should be used instead for validation etc.
func ParseContainerIDInsensitively(input string) (*ContainerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ContainerId{
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
	if resourceId.StorageAccountName, err = id.PopSegment(storageAccountsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'containers' segment
	containersKey := "containers"
	for key := range id.Path {
		if strings.EqualFold(key, containersKey) {
			containersKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(containersKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
