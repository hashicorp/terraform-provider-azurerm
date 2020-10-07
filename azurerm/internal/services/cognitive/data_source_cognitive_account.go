package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmCognitiveAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmCognitiveAccountRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tier": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmCognitiveAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cognitiveAccountsClient
	ctx := meta.(*ArmClient).StopContext

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
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("kind", resp.Kind)

	if err = d.Set("sku", flattenCognitiveAccountSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error reading `sku`: %+v", err)
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("endpoint", props.Endpoint)
	}

	d.Set("primary_access_key", keys.Key1)

	d.Set("secondary_access_key", keys.Key2)

	flattenAndSetTags(d, resp.Tags)

	return nil
}
