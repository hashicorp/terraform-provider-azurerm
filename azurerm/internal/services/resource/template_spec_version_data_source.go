package resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceTemplateSpecVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTemplateSpecVersionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty, //TODO - check validation for names
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty, //TODO - Check validation for version string
			},
			// Should only need ID from here, but we can potentially surface JSON body data etc if needed
		},
	}
}

func dataSourceTemplateSpecVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecsVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewTemplateSpecVersionID(subscriptionId, resourceGroup, name, version)

	resp, err := client.Get(ctx, id.ResourceGroup, id.TemplateSpecName, id.VersionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Templatespec %q, with version name %q (resource group %q) was not found: %+v", id.TemplateSpecName, id.VersionName, id.ResourceGroup, err)
		}
		return fmt.Errorf("reading Templatespec %q, with version name %q (resource group %q): %+v", id.TemplateSpecName, id.VersionName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	return nil
}
