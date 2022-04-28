package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceManagerSecretId struct {
	SubscriptionId string
	ResourceGroup  string
	VaultName      string
	SecretName     string
}

func NewResourceManagerSecretID(subscriptionId, resourceGroup, vaultName, secretName string) ResourceManagerSecretId {
	return ResourceManagerSecretId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VaultName:      vaultName,
		SecretName:     secretName,
	}
}

func (id ResourceManagerSecretId) String() string {
	segments := []string{
		fmt.Sprintf("Secret Name %q", id.SecretName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Manager Secret", segmentsStr)
}

func (id ResourceManagerSecretId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/secrets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.SecretName)
}

// ResourceManagerSecretID parses a ResourceManagerSecret ID into an ResourceManagerSecretId struct
func ResourceManagerSecretID(input string) (*ResourceManagerSecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceManagerSecretId{
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
	if resourceId.SecretName, err = id.PopSegment("secrets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ResourceManagerSecretIDInsensitively parses an ResourceManagerSecret ID into an ResourceManagerSecretId struct, insensitively
// This should only be used to parse an ID for rewriting, the ResourceManagerSecretID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ResourceManagerSecretIDInsensitively(input string) (*ResourceManagerSecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceManagerSecretId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'vaults' segment
	vaultsKey := "vaults"
	for key := range id.Path {
		if strings.EqualFold(key, vaultsKey) {
			vaultsKey = key
			break
		}
	}
	if resourceId.VaultName, err = id.PopSegment(vaultsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'secrets' segment
	secretsKey := "secrets"
	for key := range id.Path {
		if strings.EqualFold(key, secretsKey) {
			secretsKey = key
			break
		}
	}
	if resourceId.SecretName, err = id.PopSegment(secretsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
