package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmFunctionAppHostKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmFunctionAppHostKeysRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"function_keys": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Sensitive: true,
			},

			"system_keys": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Sensitive: true,
			},
		},
	}
}

func dataSourceArmFunctionAppHostKeysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	res, err := client.ListHostKeys(ctx, resourceGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(res.Response) {
			return fmt.Errorf("Error: AzureRM Function App %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on AzureRM Function App Hostkeys %q: %+v", name, err)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("master_key", res.MasterKey); err != nil {
		return err
	}
	if err = d.Set("function_keys", res.FunctionKeys); err != nil {
		return err
	}
	if err = d.Set("system_keys", res.SystemKeys); err != nil {
		return err
	}

	return nil
}
