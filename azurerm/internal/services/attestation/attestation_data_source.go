package attestation

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAttestation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAttestationRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmAttestationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Attestation.ProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Attestation %q does not exist", name)
		}
		return fmt.Errorf("retrieving Attestation %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Attestation %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	return tags.FlattenAndSet(d, resp.Tags)
}
