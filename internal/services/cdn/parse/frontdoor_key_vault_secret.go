package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorKeyVaultSecretId struct {
	SubscriptionId string
	ResourceGroup  string
	VaultName      string
	SecretName     string
}

func NewFrontdoorKeyVaultSecretID(subscriptionId, resourceGroup, vaultName, secretName string) FrontdoorKeyVaultSecretId {
	return FrontdoorKeyVaultSecretId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VaultName:      vaultName,
		SecretName:     secretName,
	}
}

func (id FrontdoorKeyVaultSecretId) String() string {
	segments := []string{
		fmt.Sprintf("Secret Name %q", id.SecretName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Key Vault Secret", segmentsStr)
}

func (id FrontdoorKeyVaultSecretId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/secrets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.SecretName)
}

// FrontdoorKeyVaultSecretID parses a FrontdoorKeyVaultSecret ID into an FrontdoorKeyVaultSecretId struct
func FrontdoorKeyVaultSecretID(input string) (*FrontdoorKeyVaultSecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorKeyVaultSecretId{
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

// FrontdoorKeyVaultSecretIDInsensitively parses an FrontdoorKeyVaultSecret ID into an FrontdoorKeyVaultSecretId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontdoorKeyVaultSecretID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontdoorKeyVaultSecretIDInsensitively(input string) (*FrontdoorKeyVaultSecretId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorKeyVaultSecretId{
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
