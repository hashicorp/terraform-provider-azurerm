package healthcare

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceHealthcareService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceHealthcareServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"cosmosdb_throughput": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"cosmosdb_key_vault_key_versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_policy_object_ids": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"authentication_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"audience": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"smart_proxy_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"cors_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"allowed_headers": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"allowed_methods": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"max_age_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"allow_credentials": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceHealthcareServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("HealthCare Service %q was not found in Resource Group %q", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on Azure Healthcare Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	if kind := resp.Kind; string(kind) != "" {
		d.Set("kind", kind)
	}

	if props := resp.Properties; props != nil {
		if err := d.Set("access_policy_object_ids", flattenHealthcareAccessPolicies(props.AccessPolicies)); err != nil {
			return fmt.Errorf("Error setting `access_policy_object_ids`: %+v", err)
		}

		cosmodDbKeyVaultKeyVersionlessId := ""
		cosmosDbThroughput := 0
		if cosmos := props.CosmosDbConfiguration; cosmos != nil {
			if v := cosmos.OfferThroughput; v != nil {
				cosmosDbThroughput = int(*v)
			}
			if v := cosmos.KeyVaultKeyURI; v != nil {
				cosmodDbKeyVaultKeyVersionlessId = *v
			}
		}
		d.Set("cosmosdb_key_vault_key_versionless_id", cosmodDbKeyVaultKeyVersionlessId)
		d.Set("cosmosdb_throughput", cosmosDbThroughput)

		if err := d.Set("authentication_configuration", flattenHealthcareAuthConfig(props.AuthenticationConfiguration)); err != nil {
			return fmt.Errorf("Error setting `authentication_configuration`: %+v", err)
		}

		if err := d.Set("cors_configuration", flattenHealthcareCorsConfig(props.CorsConfiguration)); err != nil {
			return fmt.Errorf("Error setting `cors_configuration`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
