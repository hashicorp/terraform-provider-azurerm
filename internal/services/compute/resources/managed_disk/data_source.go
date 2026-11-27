// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_disk

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Managed Disk.",
		Read:        dataSourceManagedDiskRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the Managed Disk.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"create_option": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"disk_access_id": {
				Description: "The ID of the Disk Access resource for using Private Endpoints on disks.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"disk_encryption_set_id": {
				Description: "The ID of the Disk Encryption Set used to encrypt this Managed Disk.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"disk_iops_read_write": {
				Description: "The number of IOPS allowed for this disk, where one operation can transfer between 4k and 256k bytes.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"disk_mbps_read_write": {
				Description: "The bandwidth allowed for this disk.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"disk_size_gb": {
				Description: "The size of the Managed Disk in gigabytes.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
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
										Description: "The URL to the Key Vault Secret used as the Disk Encryption Key.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},

									"source_vault_id": {
										Description: "The ID of the source Key Vault.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
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
										Description: "The URL to the Key Vault Key used as the Key Encryption Key.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},

									"source_vault_id": {
										Description: "The ID of the source Key Vault.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},

			"location": {
				Description: "The Azure location of the Managed Disk.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"image_reference_id": {
				Description: "The ID of the source image used for creating this Managed Disk.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"network_access_policy": {
				Description: "Policy for accessing the disk via network.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"os_type": {
				Description: "The operating system used for this Managed Disk.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"source_resource_id": {
				Description: "The ID of an existing Managed Disk which this Disk was created from.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"source_uri": {
				Description: "The Source URI for this Managed Disk.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage_account_id": {
				Description: "The ID of the Storage Account where the `source_uri` is located.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage_account_type": {
				Description: "The Storage Account type for the Managed Disk.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),

			"zones": commonschema.ZonesMultipleComputed(),
		},
	}
}

func dataSourceManagedDiskRead(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).DisksClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := commonids.NewManagedDiskID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.DiskName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		d.Set("location", location.Normalize(model.Location))

		storageAccountType := ""
		if sku := model.Sku; sku != nil {
			storageAccountType = string(*sku.Name)
		}
		d.Set("storage_account_type", storageAccountType)

		if props := model.Properties; props != nil {
			creationData := props.CreationData
			d.Set("create_option", string(creationData.CreateOption))

			imageReferenceID := ""
			if creationData.ImageReference != nil && creationData.ImageReference.Id != nil {
				imageReferenceID = *creationData.ImageReference.Id
			}
			d.Set("image_reference_id", imageReferenceID)

			d.Set("source_resource_id", creationData.SourceResourceId)
			d.Set("source_uri", creationData.SourceUri)
			d.Set("storage_account_id", creationData.StorageAccountId)

			d.Set("disk_access_id", pointer.From(props.DiskAccessId))

			d.Set("network_access_policy", string(*props.NetworkAccessPolicy))
			d.Set("disk_size_gb", props.DiskSizeGB)
			d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
			d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
			d.Set("os_type", pointer.FromEnum(props.OsType))

			diskEncryptionSetId := ""
			if props.Encryption != nil && props.Encryption.DiskEncryptionSetId != nil {
				diskEncryptionSetId = *props.Encryption.DiskEncryptionSetId
			}
			d.Set("disk_encryption_set_id", diskEncryptionSetId)

			if err := d.Set("encryption_settings", helpers.FlattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
				return fmt.Errorf("setting `encryption_settings`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
