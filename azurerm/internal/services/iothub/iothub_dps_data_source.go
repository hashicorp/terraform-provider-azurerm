package iothub

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceIotHubDPS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIotHubDPSRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"allocation_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"device_provisioning_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"id_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_operations_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceIotHubDPSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, name, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: IoT Device Provisioning Service %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving IoT Device Provisioning Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.Properties; props != nil {
		d.Set("service_operations_host_name", props.ServiceOperationsHostName)
		d.Set("device_provisioning_host_name", props.DeviceProvisioningHostName)
		d.Set("id_scope", props.IDScope)
		d.Set("allocation_policy", props.AllocationPolicy)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
