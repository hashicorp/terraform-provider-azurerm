package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmFhirApiService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmFhirApiServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type: schema.TypeString,
				Default: "fhir",
				Required: true,
			},

			"cosmodb_throughput" : {
				Type: schema.TypeInt,
				Default: 1000,
				Required: true,
			},

			"accees_policy_object_ids": {
				Type: 		schema.TypeList,
				Required: 	true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ObjectID": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},

			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceArmFhirApiServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).fhirApiServiceClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] FHIR API Service %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure FHIR API Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	// ToDo Fix this here. replace with FHIR specific values. What do we need?
	/*
		if properties := resp.ServiceProperties; properties != nil {
			d.Set("tier", string(properties.CurrentTier))

			d.Set("encryption_state", string(properties.EncryptionState))
			d.Set("firewall_state", string(properties.FirewallState))
			d.Set("firewall_allow_azure_ips", string(properties.FirewallAllowAzureIps))

			if config := properties.EncryptionConfig; config != nil {
				d.Set("encryption_type", string(config.Type))
			}

			d.Set("endpoint", properties.Endpoint)
		} */

	flattenAndSetTags(d, resp.Tags)

	return nil
}
