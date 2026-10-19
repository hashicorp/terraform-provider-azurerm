// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_disk

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/managed_disk/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Disk struct {
	Name                string              `tfschema:"name"`
	ID                  string              `tfschema:"id"`
	CreateOption        string              `tfschema:"create_option"`
	DiskAccessId        string              `tfschema:"disk_access_id"`
	DiskEncryptionSetId string              `tfschema:"disk_encryption_set_id"`
	DiskIOPSReadWrite   int64               `tfschema:"disk_iops_read_write"`
	DiskMBPSReadWrite   int64               `tfschema:"disk_mbps_read_write"`
	DiskSizeGB          int64               `tfschema:"disk_size_in_gb"`
	EncryptionSettings  []EncryptionSetting `tfschema:"encryption_settings"`
	Location            string              `tfschema:"location"`
	ImageReferenceID    string              `tfschema:"image_reference_id"`
	NetworkAccessPolicy string              `tfschema:"network_access_policy"`
	OSType              string              `tfschema:"os_type"`
	SourceResourceID    string              `tfschema:"source_resource_id"`
	SourceURI           string              `tfschema:"source_uri"`
	StorageAccountID    string              `tfschema:"storage_account_id"`
	StorageAccountType  string              `tfschema:"storage_account_type"`
	Tags                map[string]string   `tfschema:"tags"`
	Zones               []string            `tfschema:"zones"`
}

type Properties struct {
	Output Output `tfschema:"output"`
}

type EncryptionSetting struct {
	Enabled            bool                `tfschema:"enabled"`
	DiskEncryptionKeys []DiskEncryptionKey `tfschema:"disk_encryption_keys"`
	KeyEncryptionKeys  []KeyEncryptionKey  `tfschema:"key_encryption_keys"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

type DiskEncryptionKey struct {
	SecretURL     string `tfschema:"secret_url"`
	SourceVaultID string `tfschema:"source_vault_id"`
}

type KeyEncryptionKey struct {
	KeyURL        string `tfschema:"key_url"`
	SourceVaultID string `tfschema:"source_vault_id"`
}

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-name "empty"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedDiskCreate,
		Read:   resourceManagedDiskRead,
		Update: resourceManagedDiskUpdate,
		Delete: resourceManagedDiskDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ManagedDiskV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.ManagedDiskId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.ManagedDiskId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: managedDiskSchema(),

		// Encryption Settings cannot be disabled once enabled
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("encryption_settings", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) > 0 && len(new.([]any)) == 0
			}),
		),
	}
}
