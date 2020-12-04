package cognitive

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

func dataSourceCognitiveAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCognitiveAccountRead,

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

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"qna_runtime_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceCognitiveAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.GetProperties(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Cognitive Services Account %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading the state of AzureRM Cognitive Services Account %q: %+v", name, err)
	}

	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Keys for Cognitive Services Account %q (Resource Group %q) were not found", name, resourceGroup)
		}
		return fmt.Errorf("Error obtaining keys for Cognitive Services Account %q in Resource Group %q: %v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("kind", resp.Kind)

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if props := resp.Properties; props != nil {
		if apiProps := props.APIProperties; apiProps != nil {
			d.Set("qna_runtime_endpoint", apiProps.QnaRuntimeEndpoint)
		}
		d.Set("endpoint", props.Endpoint)
	}

	d.Set("primary_access_key", keys.Key1)
	d.Set("secondary_access_key", keys.Key2)

	return tags.FlattenAndSet(d, resp.Tags)
}
