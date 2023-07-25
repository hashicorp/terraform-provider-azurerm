// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSnapshot() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSnapshotRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			// Computed
			"os_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"disk_size_gb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"time_created": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"creation_option": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"source_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"source_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"storage_account_id": {
				Type:     pluginsdk.TypeString,
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

			"trusted_launch_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSnapshotRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := snapshots.NewSnapshotID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("loading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			osType := ""
			if props.OsType != nil {
				osType = string(*props.OsType)
			}
			d.Set("os_type", osType)

			timeCreated := ""
			if props.TimeCreated != nil {
				t, err := time.Parse(time.RFC3339, *props.TimeCreated)
				if err != nil {
					return fmt.Errorf("converting `time_reated`: %+v", err)
				}
				timeCreated = t.Format(time.RFC3339)
			}
			d.Set("time_created", timeCreated)

			diskSizeGb := 0
			if props.DiskSizeGB != nil {
				diskSizeGb = int(*props.DiskSizeGB)
			}
			d.Set("disk_size_gb", diskSizeGb)

			if err := d.Set("encryption_settings", flattenSnapshotDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
				return fmt.Errorf("setting `encryption_settings`: %+v", err)
			}

			trustedLaunchEnabled := false
			if securityProfile := props.SecurityProfile; securityProfile != nil && securityProfile.SecurityType != nil {
				trustedLaunchEnabled = *securityProfile.SecurityType == snapshots.DiskSecurityTypesTrustedLaunch
			}
			d.Set("trusted_launch_enabled", trustedLaunchEnabled)

			data := props.CreationData
			d.Set("creation_option", string(data.CreateOption))
			d.Set("source_uri", data.SourceUri)
			d.Set("source_resource_id", data.SourceResourceId)
			d.Set("storage_account_id", data.StorageAccountId)
		}
	}

	return nil
}
