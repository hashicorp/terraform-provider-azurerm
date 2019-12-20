package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmHealthcareService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmHealthcareServiceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
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
				Type:     schema.TypeSet,
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
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allowed_headers": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"max_age_in_seconds": {
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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmHealthcareServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
	if kind := resp.Kind; string(kind) != "" {
		d.Set("kind", kind)
	}

	if properties := resp.Properties; properties != nil {
		if accessPolicies := properties.AccessPolicies; accessPolicies != nil {
			d.Set("access_policy_object_ids", flattenHealthcareAccessPolicies(accessPolicies))
		}

		if config := properties.CosmosDbConfiguration; config != nil {
			d.Set("cosmosdb_throughput", config.OfferThroughput)
		}

		if authConfig := properties.AuthenticationConfiguration; authConfig != nil {
			if err := d.Set("authentication_configuration", flattenHealthcareAuthConfig(authConfig)); err != nil {
				return fmt.Errorf("Error setting `authentication_configuration`: %+v", flattenHealthcareAuthConfig(authConfig))
			}
		}

		if corsConfig := properties.CorsConfiguration; corsConfig != nil {
			if err := d.Set("cors_configuration", flattenHealthcareCorsConfig(corsConfig)); err != nil {
				return fmt.Errorf("Error setting `cors_configuration`: %+v", flattenHealthcareCorsConfig(corsConfig))
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
