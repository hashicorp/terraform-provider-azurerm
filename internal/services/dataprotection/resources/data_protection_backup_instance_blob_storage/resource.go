// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_blob_storage -service-package-name dataprotection -properties "name" -compare-values "subscription_id:vault_id,resource_group_name:vault_id,backup_vault_name:vault_id"

type BackupDatasourceParameters struct {
	IncludedNamespaces          []string `tfschema:"included_namespaces"`
	IncludedResourceTypes       []string `tfschema:"included_resource_types"`
	ExcludedNamespaces          []string `tfschema:"excluded_namespaces"`
	ExcludedResourceTypes       []string `tfschema:"excluded_resource_types"`
	LabelSelectors              []string `tfschema:"label_selectors"`
	VolumeSnapshotEnabled       bool     `tfschema:"volume_snapshot_enabled"`
	ClusterScopeResourceEnabled bool     `tfschema:"cluster_scoped_resources_enabled"`
}

func RegisterResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupInstanceBlobStorageCreateUpdate,
		Read:   resourceDataProtectionBackupInstanceBlobStorageRead,
		Update: resourceDataProtectionBackupInstanceBlobStorageCreateUpdate,
		Delete: resourceDataProtectionBackupInstanceBlobStorageDelete,

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

		Schema: dataProtectionBackupInstanceBlobStorageSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// The `storage_account_container_names` can not be removed once specified.
			pluginsdk.ForceNewIfChange("storage_account_container_names", func(ctx context.Context, old, new, _ any) bool {
				return len(old.([]any)) > 0 && len(new.([]any)) == 0
			}),
		),
	}
}
