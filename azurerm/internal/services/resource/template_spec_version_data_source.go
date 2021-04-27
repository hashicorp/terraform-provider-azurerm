package resource

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

		//lintignore:S033
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateSpecName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateSpecVersionName,
			},

			"template_body": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceTemplateSpecVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecsVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewTemplateSpecVersionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string), d.Get("version").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.TemplateSpecName, id.VersionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Templatespec %q, with version name %q (resource group %q) was not found: %+v", id.TemplateSpecName, id.VersionName, id.ResourceGroup, err)
		}
		return fmt.Errorf("reading Templatespec %q, with version name %q (resource group %q): %+v", id.TemplateSpecName, id.VersionName, id.ResourceGroup, err)
	}

	templateBody := "{}"
	if props := resp.VersionProperties; props != nil && props.Template != nil {
		templateBodyRaw, err := flattenTemplateDeploymentBody(props.Template)
		if err != nil {
			return err
		}

		templateBody = *templateBodyRaw
	}
	d.Set("template_body", templateBody)

	d.SetId(id.ID())

	return tags.FlattenAndSet(d, resp.Tags)
}
