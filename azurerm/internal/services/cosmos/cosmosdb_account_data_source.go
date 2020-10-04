package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmCosmosDbAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmCosmosDbAccountRead,

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

			"tags": tags.SchemaDataSource(),

			"offer_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_range_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_free_tier": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_automatic_failover": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"consistency_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"consistency_level": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"max_interval_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"max_staleness_prefix": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"geo_location": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"failover_priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"capabilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"is_virtual_network_filter_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"virtual_network_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"enable_multiple_write_locations": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"read_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"write_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_readonly_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_readonly_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_master_key": {
				Type:       schema.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `primary_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},

			"secondary_master_key": {
				Type:       schema.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `secondary_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},

			"primary_readonly_master_key": {
				Type:       schema.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `primary_readonly_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},

			"secondary_readonly_master_key": {
				Type:       schema.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `secondary_readonly_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},
		},
	}
}

func dataSourceArmCosmosDbAccountRead(d *schema.ResourceData, meta interface{}) error {
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
