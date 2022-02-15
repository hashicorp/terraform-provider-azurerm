package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

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
		},
	}
}

func dataSourceSnapshotRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSnapshotID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("loading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if props := resp.SnapshotProperties; props != nil {
		d.Set("os_type", string(props.OsType))
		d.Set("time_created", props.TimeCreated.String())

		if props.DiskSizeGB != nil {
			d.Set("disk_size_gb", int(*props.DiskSizeGB))
		}

		if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
			return fmt.Errorf("setting `encryption_settings`: %+v", err)
		}
	}

	if data := resp.CreationData; data != nil {
		d.Set("creation_option", string(data.CreateOption))
		d.Set("source_uri", data.SourceURI)
		d.Set("source_resource_id", data.SourceResourceID)
		d.Set("storage_account_id", data.StorageAccountID)
	}

	return nil
}
