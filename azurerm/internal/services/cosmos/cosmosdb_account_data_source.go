package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceCosmosDbAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCosmosDbAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"tags": tags.SchemaDataSource(),

			"offer_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ip_range_filter": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enable_free_tier": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"enable_automatic_failover": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"consistency_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"consistency_level": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"max_interval_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"max_staleness_prefix": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"geo_location": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"location": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"failover_priority": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"capabilities": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"is_virtual_network_filter_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"virtual_network_rule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"key_vault_key_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enable_multiple_write_locations": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"read_endpoints": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"write_endpoints": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_readonly_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_readonly_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_master_key": {
				Type:       pluginsdk.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `primary_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},

			"secondary_master_key": {
				Type:       pluginsdk.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `secondary_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},

			"primary_readonly_master_key": {
				Type:       pluginsdk.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `primary_readonly_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},

			"secondary_readonly_master_key": {
				Type:       pluginsdk.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `secondary_readonly_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},
		},
	}
}

func dataSourceCosmosDbAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: CosmosDB Account %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on AzureRM CosmosDB Account %s (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("kind", string(resp.Kind))

	if props := resp.DatabaseAccountGetProperties; props != nil {
		d.Set("offer_type", string(props.DatabaseAccountOfferType))
		d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilter(props.IPRules))
		d.Set("endpoint", props.DocumentEndpoint)
		d.Set("is_virtual_network_filter_enabled", resp.IsVirtualNetworkFilterEnabled)
		d.Set("enable_free_tier", resp.EnableFreeTier)
		d.Set("enable_automatic_failover", resp.EnableAutomaticFailover)

		if v := props.KeyVaultKeyURI; v != nil {
			d.Set("key_vault_key_id", resp.KeyVaultKeyURI)
		}

		if err = d.Set("consistency_policy", flattenAzureRmCosmosDBAccountConsistencyPolicy(resp.ConsistencyPolicy)); err != nil {
			return fmt.Errorf("Error setting `consistency_policy`: %+v", err)
		}

		// sort `geo_locations` by fail over priority
		locations := make([]map[string]interface{}, len(*props.FailoverPolicies))
		for _, l := range *props.FailoverPolicies {
			locations[*l.FailoverPriority] = map[string]interface{}{
				"id":                *l.ID,
				"location":          azure.NormalizeLocation(*l.LocationName),
				"failover_priority": int(*l.FailoverPriority),
			}
		}
		if err = d.Set("geo_location", locations); err != nil {
			return fmt.Errorf("Error setting `geo_location`: %+v", err)
		}

		if err = d.Set("capabilities", flattenAzureRmCosmosDBAccountCapabilitiesAsList(resp.Capabilities)); err != nil {
			return fmt.Errorf("Error setting `capabilities`: %+v", err)
		}

		if err = d.Set("virtual_network_rule", flattenAzureRmCosmosDBAccountVirtualNetworkRulesAsList(props.VirtualNetworkRules)); err != nil {
			return fmt.Errorf("Error setting `virtual_network_rule`: %+v", err)
		}

		readEndpoints := make([]string, 0)
		if locations := props.ReadLocations; locations != nil {
			for _, l := range *locations {
				if l.DocumentEndpoint == nil {
					continue
				}

				readEndpoints = append(readEndpoints, *l.DocumentEndpoint)
			}
		}
		if err := d.Set("read_endpoints", readEndpoints); err != nil {
			return fmt.Errorf("Error setting `read_endpoints`: %s", err)
		}

		writeEndpoints := make([]string, 0)
		if locations := props.WriteLocations; locations != nil {
			for _, l := range *locations {
				if l.DocumentEndpoint == nil {
					continue
				}

				writeEndpoints = append(writeEndpoints, *l.DocumentEndpoint)
			}
		}
		if err := d.Set("write_endpoints", writeEndpoints); err != nil {
			return fmt.Errorf("Error setting `write_endpoints`: %s", err)
		}

		d.Set("enable_multiple_write_locations", resp.EnableMultipleWriteLocations)
	}

	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Write keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_key", keys.PrimaryMasterKey)
		d.Set("secondary_key", keys.SecondaryMasterKey)
		d.Set("primary_master_key", keys.PrimaryMasterKey)
		d.Set("secondary_master_key", keys.SecondaryMasterKey)
	}

	readonlyKeys, err := client.ListReadOnlyKeys(ctx, resourceGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List read-only keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_readonly_key", readonlyKeys.PrimaryReadonlyMasterKey)
		d.Set("secondary_readonly_key", readonlyKeys.SecondaryReadonlyMasterKey)
		d.Set("primary_readonly_master_key", readonlyKeys.PrimaryReadonlyMasterKey)
		d.Set("secondary_readonly_master_key", readonlyKeys.SecondaryReadonlyMasterKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenAzureRmCosmosDBAccountCapabilitiesAsList(capabilities *[]documentdb.Capability) *[]map[string]interface{} {
	slice := make([]map[string]interface{}, 0)

	for _, c := range *capabilities {
		if v := c.Name; v != nil {
			e := map[string]interface{}{
				"name": *v,
			}
			slice = append(slice, e)
		}
	}

	return &slice
}

func flattenAzureRmCosmosDBAccountVirtualNetworkRulesAsList(rules *[]documentdb.VirtualNetworkRule) []map[string]interface{} {
	if rules == nil {
		return []map[string]interface{}{}
	}

	virtualNetworkRules := make([]map[string]interface{}, len(*rules))
	for i, r := range *rules {
		virtualNetworkRules[i] = map[string]interface{}{
			"id": *r.ID,
		}
	}
	return virtualNetworkRules
}
