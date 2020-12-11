package databox

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databox/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databox/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDataBoxJob() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataBoxJobRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DataBoxJobName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"destination_managed_disk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"share_password": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"staging_storage_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"destination_storage_account": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_password": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"storage_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"device_password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDataBoxJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBox.JobClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewDataBoxJobId(subscriptionId, name, resourceGroup)

	resp, err := client.Get(ctx, resourceGroup, name, "Details")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Data Box Job (DataBox Job Name %q / Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("reading Data Box Job (Data Box Job Name %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("sku_name", resp.Sku.Name)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.JobProperties; props != nil {
		if details := props.Details; details != nil {
			if v, ok := details.AsJobDetailsType(); ok && v != nil {
				d.Set("device_password", v.DevicePassword)

				destinationManagedDisk, destinationStorageAccount := flattenArmDataBoxJobDestinationAccount(v.DestinationAccountDetails)
				if err := d.Set("destination_managed_disk", destinationManagedDisk); err != nil {
					return fmt.Errorf("setting `destination_managed_disk`: %+v", err)
				}

				if err := d.Set("destination_storage_account", destinationStorageAccount); err != nil {
					return fmt.Errorf("setting `destination_storage_account`: %+v", err)
				}
			} else if v, ok := details.AsDiskJobDetails(); ok && v != nil {
				destinationManagedDisk, destinationStorageAccount := flattenArmDataBoxJobDestinationAccount(v.DestinationAccountDetails)
				if err := d.Set("destination_managed_disk", destinationManagedDisk); err != nil {
					return fmt.Errorf("setting `destination_managed_disk`: %+v", err)
				}

				if err := d.Set("destination_storage_account", destinationStorageAccount); err != nil {
					return fmt.Errorf("setting `destination_storage_account`: %+v", err)
				}
			} else if v, ok := details.AsHeavyJobDetails(); ok && v != nil {
				d.Set("device_password", v.DevicePassword)

				destinationManagedDisk, destinationStorageAccount := flattenArmDataBoxJobDestinationAccount(v.DestinationAccountDetails)
				if err := d.Set("destination_managed_disk", destinationManagedDisk); err != nil {
					return fmt.Errorf("setting `destination_managed_disk`: %+v", err)
				}

				if err := d.Set("destination_storage_account", destinationStorageAccount); err != nil {
					return fmt.Errorf("setting `destination_storage_account`: %+v", err)
				}
			}
		}
	}

	d.SetId(id.ID(""))

	return tags.FlattenAndSet(d, resp.Tags)
}
