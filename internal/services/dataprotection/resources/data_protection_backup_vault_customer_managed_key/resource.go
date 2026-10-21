// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault_customer_managed_key

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id "data_protection_backup_vault_id" -test-name complete

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &backupvaultresources.BackupVaultId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_data_protection_backup_vault_customer_managed_key"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupvaultresources.ValidateBackupVaultID
}
