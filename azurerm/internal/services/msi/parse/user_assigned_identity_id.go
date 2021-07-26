package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/sdk/managedidentity"

// these are here primarily to enable migration over

func UserAssignedIdentityID(input string) (*managedidentity.UserAssignedIdentitiesId, error) {
	return managedidentity.ParseUserAssignedIdentitiesID(input)
}

func UserAssignedIdentityIDInsensitively(input string) (*managedidentity.UserAssignedIdentitiesId, error) {
	return managedidentity.ParseUserAssignedIdentitiesIDInsensitively(input)
}
