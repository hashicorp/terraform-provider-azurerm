package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSshPublicKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSshPublicKeyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmSshPublicKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SSHPublicKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: SSH Public Key %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("[ERROR] Error making Read request on Azure SSH Public Key %q (Resource Group %q): %s", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	if props := resp.SSHPublicKeyResourceProperties; props != nil {
		d.Set("public_key", props.PublicKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
