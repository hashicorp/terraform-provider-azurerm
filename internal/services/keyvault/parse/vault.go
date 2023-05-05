package parse

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type VaultId = commonids.KeyVaultId

func NewVaultID(subscriptionId, resourceGroup, name string) commonids.KeyVaultId {
	return commonids.NewKeyVaultID(subscriptionId, resourceGroup, name)
}

// VaultID parses a Vault ID into an VaultId struct
func VaultID(input string) (*VaultId, error) {
	return commonids.ParseKeyVaultID(input)
}

// VaultIDInsensitively parses an Vault ID into an VaultId struct, insensitively
// This should only be used to parse an ID for rewriting, the VaultID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func VaultIDInsensitively(input string) (*VaultId, error) {
	return commonids.ParseKeyVaultIDInsensitively(input)
}
