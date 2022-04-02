package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type KeyId struct {
	SubscriptionId string
	ResourceGroup  string
	VaultName      string
	Name           string
}

func NewKeyID(subscriptionId, resourceGroup, vaultName, name string) KeyId {
	return KeyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VaultName:      vaultName,
		Name:           name,
	}
}

func (id KeyId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Key", segmentsStr)
}

func (id KeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.Name)
}

// KeyID parses a Key ID into an KeyId struct
func KeyID(input string) (*KeyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := KeyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("keys"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
