package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmHealthcareService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmHealthcareServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cosmosdb_throughput": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"access_policy_object_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			/*
				"cors_configuration": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"origins": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"headers": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"methods": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"max_age": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"allow_credentials": {
								Type:     schema.TypeBool,
								Computed: true,
							},
						},
					},
				},
			*/
			"tags": tagsSchema(),
		},
	}
}

func dataSourceArmHealthcareServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).healthcare.HealthcareServiceClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Healthcare Service %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return fmt.Errorf("HealthCare Service %q was not found in Resource Group %q", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on Azure Healthcare Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	flattenAndSetTags(d, resp.Tags)

	return nil
}
