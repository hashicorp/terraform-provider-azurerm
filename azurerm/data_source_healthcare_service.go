package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/healthcareapis/mgmt/2018-08-20-preview/healthcareapis"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"authentication_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"audience": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"smart_proxy_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"cors_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"allowed_headers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"max_age_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"allow_credentials": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

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
	if kind := resp.Kind; kind != "" {
		d.Set("kind", kind)
	}

	if properties := resp.Properties; properties != nil {
		if accessPolicies := properties.AccessPolicies; accessPolicies != nil {
			d.Set("access_policy_object_ids", flattenHealthcareAccessPolicies(accessPolicies))
		}

		if config := properties.CosmosDbConfiguration; config != nil {
			d.Set("cosmosdb_throughput", config.OfferThroughput)
		}

		authOutput := make([]interface{}, 0)
		if authConfig := properties.AuthenticationConfiguration; authConfig != nil {
			output := make(map[string]interface{})
			if authConfig.Authority != nil {
				output["authority"] = *authConfig.Authority
			}
			if authConfig.Audience != nil {
				output["audience"] = *authConfig.Audience
			}
			if authConfig.SmartProxyEnabled != nil {
				output["smart_proxy_enabled"] = *authConfig.SmartProxyEnabled
			}
			authOutput = append(authOutput, output)
		}

		if err := d.Set("authentication_configuration", authOutput); err != nil {
			return fmt.Errorf("Error setting `authentication_configuration`: %+v", authOutput)
		}

		corsOutput := make([]interface{}, 0)
		if corsConfig := properties.CorsConfiguration; corsConfig != nil {
			output := make(map[string]interface{})
			if corsConfig.Origins != nil {
				output["allowed_origins"] = *corsConfig.Origins
			}
			if corsConfig.Headers != nil {
				output["allowed_headers"] = *corsConfig.Headers
			}
			if corsConfig.Methods != nil {
				output["allowed_methods"] = *corsConfig.Methods
			}
			if corsConfig.MaxAge != nil {
				output["max_age_in_seconds"] = *corsConfig.MaxAge
			}
			if corsConfig.AllowCredentials != nil {
				output["allow_credentials"] = *corsConfig.AllowCredentials
			}
			corsOutput = append(corsOutput, output)
		}

		if err := d.Set("cors_configuration", corsOutput); err != nil {
			return fmt.Errorf("Error setting `cors_configuration`: %+v", corsOutput)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func flattenHealthcareAccessPolicies(policies *[]healthcareapis.ServiceAccessPolicyEntry) []string {
	result := make([]string, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		if objectId := policy.ObjectID; objectId != nil {
			result = append(result, *objectId)
		}
	}

	return result
}
