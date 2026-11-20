// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
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
			pluginsdk.CustomizeDiffShim(resourceDataProtectionBackupInstanceBlobStorageCustomizeDiff),
		),
	}
}

func resourceDataProtectionBackupInstanceBlobStorageCustomizeDiff(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	if !d.NewValueKnown("auto_protection_enabled") {
		return nil
	}

	oldAutoProtection, newAutoProtection := d.GetChange("auto_protection_enabled")
	autoProtectionEnabled := newAutoProtection.(bool)

	// Azure does not allow switching a Backup Instance back from auto protection to a container list (or to no
	// container selection at all) once auto protection has been enabled - the only way out is to re-create the
	// Backup Instance, which also removes its vaulted recovery points. This is deliberately an error rather than
	// `ForceNew` so that the re-creation has to be requested explicitly.
	if d.Id() != "" && oldAutoProtection.(bool) && !autoProtectionEnabled {
		return fmt.Errorf("`auto_protection_enabled` cannot be changed from `true` to `false`: enabling auto protection for a Backup Instance is irreversible in Azure. To switch back to a container list the Backup Instance has to be re-created (e.g. via `terraform apply -replace=<resource address>`), which deletes its vaulted recovery points")
	}

	if !d.NewValueKnown("storage_account_container_names") || !d.NewValueKnown("excluded_container_name_prefixes") {
		return nil
	}

	oldContainerNames, newContainerNames := d.GetChange("storage_account_container_names")

	if autoProtectionEnabled && len(newContainerNames.([]any)) > 0 {
		return fmt.Errorf("`storage_account_container_names` cannot be set when `auto_protection_enabled` is `true`: auto protection covers all present and future containers of the Storage Account")
	}

	if !autoProtectionEnabled && d.Get("excluded_container_name_prefixes").(*pluginsdk.Set).Len() > 0 {
		return fmt.Errorf("`excluded_container_name_prefixes` can only be set when `auto_protection_enabled` is `true`")
	}

	// The `storage_account_container_names` can not be removed once specified - unless the Backup Instance is being
	// migrated to auto protection, which Azure supports as an in-place update.
	if len(oldContainerNames.([]any)) > 0 && len(newContainerNames.([]any)) == 0 && !autoProtectionEnabled {
		if err := d.ForceNew("storage_account_container_names"); err != nil {
			return err
		}
	}

	return nil
}
