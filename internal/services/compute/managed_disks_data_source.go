// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedDisksDataSource struct{}

var _ sdk.DataSource = &ManagedDisksDataSource{}

func (m ManagedDisksDataSource) ModelObject() interface{} {
	return &ManagedDisksDataSourceModel{}
}

func (m ManagedDisksDataSource) ResourceType() string {
	return "azurerm_managed_disks"
}

type ManagedDisksDataSourceModel struct {
	ResourceGroupName string `tfschema:"resource_group_name"`
	Disks             []Disk `tfschema:"disk"`
}

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

type EncryptionSetting struct {
	Enabled            bool                `tfschema:"enabled"`
	DiskEncryptionKeys []DiskEncryptionKey `tfschema:"disk_encryption_keys"`
	KeyEncryptionKeys  []KeyEncryptionKey  `tfschema:"key_encryption_keys"`
}

type DiskEncryptionKey struct {
	SecretURL     string `tfschema:"secret_url"`
	SourceVaultID string `tfschema:"source_vault_id"`
}

type KeyEncryptionKey struct {
	KeyURL        string `tfschema:"key_url"`
	SourceVaultID string `tfschema:"source_vault_id"`
}

func (m ManagedDisksDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (m ManagedDisksDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disk": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"create_option": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"disk_access_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"disk_encryption_set_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"disk_iops_read_write": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"disk_mbps_read_write": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"disk_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"encryption_settings": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Computed: true,
								},

								"disk_encryption_key": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"secret_url": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"source_vault_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
								"key_encryption_key": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"key_url": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"source_vault_id": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"location": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"image_reference_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"network_access_policy": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"os_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"source_resource_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"source_uri": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"storage_account_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"storage_account_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"tags": commonschema.TagsDataSource(),

					"zones": commonschema.ZonesMultipleComputed(),

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (m ManagedDisksDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.DisksClient

			state := ManagedDisksDataSourceModel{}

			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := commonids.NewResourceGroupID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName)

			resp, err := client.ListByResourceGroup(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("making Read request on %s: %s", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				managedDisks := make([]Disk, 0)
				for _, d := range *model {
					disk := Disk{
						Name:     pointer.From(d.Name),
						ID:       pointer.From(d.Id),
						Location: location.Normalize(d.Location),
						Zones:    zones.Flatten(d.Zones),
						Tags:     pointer.From(d.Tags),
					}
					if props := d.Properties; props != nil {
						disk.CreateOption = string(props.CreationData.CreateOption)
						disk.SourceURI = pointer.From(props.CreationData.SourceUri)
						disk.SourceResourceID = pointer.From(props.CreationData.SourceResourceId)
						disk.StorageAccountID = pointer.From(props.CreationData.StorageAccountId)
						if props.CreationData.ImageReference != nil {
							disk.ImageReferenceID = pointer.From(props.CreationData.ImageReference.Id)
						}

						disk.DiskAccessId = pointer.From(props.DiskAccessId)
						disk.NetworkAccessPolicy = pointer.FromEnum(props.NetworkAccessPolicy)
						disk.DiskSizeGB = pointer.From(props.DiskSizeGB)
						disk.DiskIOPSReadWrite = pointer.From(props.DiskIOPSReadWrite)
						disk.DiskMBPSReadWrite = pointer.From(props.DiskMBpsReadWrite)
						disk.OSType = pointer.FromEnum(props.OsType)
						if enc := props.Encryption; enc != nil {
							disk.DiskEncryptionSetId = pointer.From(enc.DiskEncryptionSetId)
						}

						disk.EncryptionSettings = flattenManagedDiskEncryptionSettingsTyped(props.EncryptionSettingsCollection)
					}

					if sku := d.Sku; sku != nil {
						disk.StorageAccountType = pointer.FromEnum(sku.Name)
					}

					disk.Tags = pointer.From(d.Tags)

					managedDisks = append(managedDisks, disk)
				}
				state.Disks = managedDisks
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenManagedDiskEncryptionSettingsTyped(encryptionSettings *disks.EncryptionSettingsCollection) []EncryptionSetting {
	if encryptionSettings == nil {
		return []EncryptionSetting{}
	}

	diskEncryptionKeys := make([]DiskEncryptionKey, 0)
	keyEncryptionKeys := make([]KeyEncryptionKey, 0)
	if encryptionSettings.EncryptionSettings != nil && len(*encryptionSettings.EncryptionSettings) > 0 {
		settings := (*encryptionSettings.EncryptionSettings)[0]

		if key := settings.DiskEncryptionKey; key != nil {
			secretUrl := ""
			if key.SecretURL != "" {
				secretUrl = key.SecretURL
			}

			sourceVaultId := ""
			if key.SourceVault.Id != nil {
				sourceVaultId = *key.SourceVault.Id
			}

			diskEncryptionKeys = append(diskEncryptionKeys, DiskEncryptionKey{
				SecretURL:     secretUrl,
				SourceVaultID: sourceVaultId,
			})
		}

		if key := settings.KeyEncryptionKey; key != nil {
			keyUrl := ""
			if key.KeyURL != "" {
				keyUrl = key.KeyURL
			}

			sourceVaultId := ""
			if key.SourceVault.Id != nil {
				sourceVaultId = *key.SourceVault.Id
			}

			keyEncryptionKeys = append(keyEncryptionKeys, KeyEncryptionKey{
				KeyURL:        keyUrl,
				SourceVaultID: sourceVaultId,
			})
		}

		return []EncryptionSetting{
			{
				DiskEncryptionKeys: diskEncryptionKeys,
				KeyEncryptionKeys:  keyEncryptionKeys,
			},
		}
	}

	return []EncryptionSetting{}
}
