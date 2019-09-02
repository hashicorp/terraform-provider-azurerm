package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDataLakeStoreAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDateLakeStoreAccountRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"tier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"encryption_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"encryption_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"firewall_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"firewall_allow_azure_ips": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDateLakeStoreAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).datalake.StoreAccountsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] DataLakeStoreAccount %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Data Lake %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if properties := resp.DataLakeStoreAccountProperties; properties != nil {
		d.Set("tier", string(properties.CurrentTier))

		d.Set("encryption_state", string(properties.EncryptionState))
		d.Set("firewall_allow_azure_ips", string(properties.FirewallAllowAzureIps))
		d.Set("firewall_state", string(properties.FirewallState))

		if config := properties.EncryptionConfig; config != nil {
			d.Set("encryption_type", string(config.Type))
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
