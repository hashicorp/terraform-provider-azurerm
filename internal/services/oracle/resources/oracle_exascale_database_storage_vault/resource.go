// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exascale_database_storage_vault

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exascaledbstoragevaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_oracle_exascale_database_storage_vault"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return exascaledbstoragevaults.ValidateExascaleDbStorageVaultID
}
