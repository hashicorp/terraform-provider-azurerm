// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_disk

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_disk -service-package-name dataprotection -properties "name" -compare-values "subscription_id:vault_id,resource_group_name:vault_id,backup_vault_name:vault_id"

func RegisterResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupInstanceDiskCreateUpdate,
		Read:   resourceDataProtectionBackupInstanceDiskRead,
		Update: resourceDataProtectionBackupInstanceDiskCreateUpdate,
		Delete: resourceDataProtectionBackupInstanceDiskDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&backupinstanceresources.BackupInstanceId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&backupinstanceresources.BackupInstanceId{}),
		},

		Schema: dataProtectionBackupInstanceDiskSchema(),
	}
}
