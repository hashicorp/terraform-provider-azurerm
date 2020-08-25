package logic

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

func dataSourceLogicAppIntegrationAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmLogicAppIntegrationAccountRead,

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

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"integration_service_environment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmLogicAppIntegrationAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Logic App Integration Account Account %q does not exist in Resource Group %q", name, resourceGroup)
		}
		return fmt.Errorf("retrieving Logic App Integration Account Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading Logic App Integration Account Account %q (Resource Group %q): ID is empty or nil", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	sku := ""
	if resp.Sku != nil {
		sku = string(resp.Sku.Name)
	}
	d.Set("sku_name", sku)

	iseID := ""
	if resp.IntegrationAccountProperties != nil && resp.IntegrationAccountProperties.IntegrationServiceEnvironment != nil && resp.IntegrationAccountProperties.IntegrationServiceEnvironment.ID != nil {
		iseID = *resp.IntegrationAccountProperties.IntegrationServiceEnvironment.ID
	}
	d.Set("integration_service_environment_id", iseID)
	return tags.FlattenAndSet(d, resp.Tags)
}
