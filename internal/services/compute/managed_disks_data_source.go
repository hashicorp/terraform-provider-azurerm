// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//
// type ManagedDisksDataSource struct{}
// type ManagedDisksDataSourceModel struct {
// 	ResourceGroupName string `tfschema:"resource_group_name"`
// 	Disks             []Disk `tfschema:"disk"`
// }
//
// type Disk struct {
// 	Name                string              `tfschema:"name"`
// 	CreateOption        string              `tfschema:"create_option"`
// 	DiskAccessId        string              `tfschema:"disk_access_id"`
// 	DiskEncryptionSetId string              `tfschema:"disk_encryption_set_id"`
// 	DiskIOPSReadWrite   int64               `tfschema:"disk_iops_read_write"`
// 	DiskMBPSReadWrite   int64               `tfschema:"disk_mbps_read_write"`
// 	DiskSizeMB          int64               `tfschema:"disk_size_gb"`
// 	EncryptionSettings  []EncryptionSetting `tfschema:"encryption_settings"`
// 	Location            string              `tfschema:"location"`
// 	ImageReferenceID    string              `tfschema:"image_reference_id"`
// 	NetworkAccessPolicy string              `tfschema:"network_access_policy"`
// 	OSType              string              `tfschema:"os_type"`
// 	SourceResourceID    string              `tfschema:"source_resource_id"`
// 	SourceURI           string              `tfschema:"source_uri"`
// 	StorageAccountID    string              `tfschema:"storage_account_id"`
// 	StorageAccountType  string              `tfschema:"storage_account_type"`
// 	Tags                map[string]string   `tfschema:"tags"`
// 	Zones               []string            `tfschema:"zones"`
// }
//
// type EncryptionSetting struct {
// 	Enabled bool `tfschema:"enabled"`
// 	DiskEncryptionKey []DiskEncryptionKey `tfschema:"disk_encryption_key"`
// }

func dataSourceManagedDisks() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceManagedDisksRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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

						"disk_size_gb": {
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
		},
	}
}

func dataSourceManagedDisksRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewResourceGroupID(subscriptionId, d.Get("resource_group_name").(string))

	resp, err := client.ListByResourceGroup(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	d.SetId(id.ID())

	disks := make([]interface{}, 0)
	if model := resp.Model; model != nil {
		for _, disk := range *model {
			r := map[string]interface{}{
				"name":     pointer.From(disk.Name),
				"zones":    zones.FlattenUntyped(disk.Zones),
				"location": location.Normalize(disk.Location),
				"id":       commonids.NewManagedDiskID(id.SubscriptionId, id.ResourceGroupName, pointer.From(disk.Name)).ID(),
			}

			if props := disk.Properties; props != nil {
				r["create_option"] = props.CreationData.CreateOption
				r["source_uri"] = pointer.From(props.CreationData.SourceUri)
				r["source_resource_id"] = pointer.From(props.CreationData.SourceResourceId)
				r["storage_account_id"] = pointer.From(props.CreationData.StorageAccountId)

				if props.CreationData.ImageReference != nil {
					r["image_reference_id"] = pointer.From(props.CreationData.ImageReference.Id)
				}

				r["disk_access_id"] = pointer.From(props.DiskAccessId)
				r["network_access_policy"] = pointer.From(props.NetworkAccessPolicy)
				r["disk_size_gb"] = pointer.From(props.DiskSizeGB)
				r["disk_iops_read_write"] = pointer.From(props.DiskIOPSReadWrite)
				r["disk_mbps_read_write"] = pointer.From(props.DiskMBpsReadWrite)
				r["os_type"] = pointer.From(props.OsType)

				if enc := props.Encryption; enc != nil {
					r["disk_encryption_set_id"] = pointer.From(enc.DiskEncryptionSetId)
				}

				r["encryption_settings"] = flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)
			}

			if sku := disk.Sku; sku != nil {
				r["storage_account_type"] = pointer.From(sku.Name)
			}

			r["tags"] = tags.Flatten(disk.Tags)

			disks = append(disks, r)
		}
	}

	d.Set("disk", disks)

	d.Set("resource_group_name", id.ResourceGroupName)

	return nil
}
