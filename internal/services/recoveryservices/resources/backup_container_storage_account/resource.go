// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_container_storage_account

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages registration of a Storage Account with Azure Backup. Storage Accounts must be registered with an Azure Recovery Vault in order to backup file shares within the Storage Account. Registering a Storage Account with a vault creates what is known as a protection container within Azure Recovery Services. Once the container is created, Azure file shares within the Storage Account can be backed up using the `azurerm_backup_protected_file_share` resource.",
		Create:      resourceBackupProtectionContainerStorageAccountCreate,
		Read:        resourceBackupProtectionContainerStorageAccountRead,
		Delete:      resourceBackupProtectionContainerStorageAccountDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := protectioncontainers.ParseProtectionContainerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: backupContainerStorageAccountSchema(),
	}
}
