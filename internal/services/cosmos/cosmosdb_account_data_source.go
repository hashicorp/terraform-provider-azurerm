// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceCosmosDbAccount() *pluginsdk.Resource {
	dataSource := &pluginsdk.Resource{
		Read: dataSourceCosmosDbAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"tags": commonschema.TagsDataSource(),

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

			"free_tier_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"automatic_failover_enabled": {
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

			"multiple_write_locations_enabled": {
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

			"primary_sql_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_sql_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_readonly_sql_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_readonly_sql_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_mongodb_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_mongodb_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_readonly_mongodb_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_readonly_mongodb_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}

	return dataSource
}

func dataSourceCosmosDbAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewDatabaseAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.DatabaseAccountsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
		d.Set("resource_group_name", id.ResourceGroupName)

		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("kind", string(pointer.From(model.Kind)))

		if props := model.Properties; props != nil {
			d.Set("offer_type", string(pointer.From(props.DatabaseAccountOfferType)))
			d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilterDataSource(props.IPRules))
			d.Set("endpoint", props.DocumentEndpoint)
			d.Set("is_virtual_network_filter_enabled", props.IsVirtualNetworkFilterEnabled)
			d.Set("free_tier_enabled", props.EnableFreeTier)
			d.Set("automatic_failover_enabled", props.EnableAutomaticFailover)
			d.Set("multiple_write_locations_enabled", props.EnableMultipleWriteLocations)

			if v := props.KeyVaultKeyUri; v != nil {
				d.Set("key_vault_key_id", v)
			}

			if err = d.Set("consistency_policy", flattenAzureRmCosmosDBAccountConsistencyPolicy(props.ConsistencyPolicy)); err != nil {
				return fmt.Errorf("setting `consistency_policy`: %+v", err)
			}

			locations := make([]map[string]interface{}, len(*props.FailoverPolicies))

			// the original procedure leads to a sorted locations slice by using failover priority as index
			// sort `geo_locations` by failover priority if we found priorities were not within limitation.
			if anyUnexpectedFailoverPriority(*props.FailoverPolicies) {
				policies := *props.FailoverPolicies
				sort.Slice(policies, func(i, j int) bool {
					return *policies[i].FailoverPriority < *policies[j].FailoverPriority
				})
				for i, l := range policies {
					locations[i] = map[string]interface{}{
						"id":                *l.Id,
						"location":          location.NormalizeNilable(l.LocationName),
						"failover_priority": int(*l.FailoverPriority),
					}
				}
			} else {
				for _, l := range *props.FailoverPolicies {
					locations[*l.FailoverPriority] = map[string]interface{}{
						"id":                *l.Id,
						"location":          location.NormalizeNilable(l.LocationName),
						"failover_priority": int(*l.FailoverPriority),
					}
				}
			}
			if err = d.Set("geo_location", locations); err != nil {
				return fmt.Errorf("setting `geo_location`: %+v", err)
			}

			if err = d.Set("capabilities", flattenAzureRmCosmosDBAccountCapabilitiesAsList(props.Capabilities)); err != nil {
				return fmt.Errorf("setting `capabilities`: %+v", err)
			}

			if err = d.Set("virtual_network_rule", flattenAzureRmCosmosDBAccountVirtualNetworkRulesAsList(props.VirtualNetworkRules)); err != nil {
				return fmt.Errorf("setting `virtual_network_rule`: %+v", err)
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
				return fmt.Errorf("setting `read_endpoints`: %s", err)
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
				return fmt.Errorf("setting `write_endpoints`: %s", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	keys, err := client.DatabaseAccountsListKeys(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Unable to List Write keys for %s: %s", id, err)
	} else {
		if model := keys.Model; model != nil {
			d.Set("primary_key", model.PrimaryMasterKey)
			d.Set("secondary_key", model.SecondaryMasterKey)
		}
	}

	readonlyKeys, err := client.DatabaseAccountsListReadOnlyKeys(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Unable to List read-only keys for %s: %s", id, err)
	} else {
		if model := readonlyKeys.Model; model != nil {
			d.Set("primary_readonly_key", model.PrimaryReadonlyMasterKey)
			d.Set("secondary_readonly_key", model.SecondaryReadonlyMasterKey)
		}
	}

	connStringResp, err := client.DatabaseAccountsListConnectionStrings(ctx, id)
	if err != nil {
		if response.WasNotFound(keys.HttpResponse) {
			log.Printf("[DEBUG] Connection Strings were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.DatabaseAccountName, id.ResourceGroupName)
		} else {
			log.Printf("[ERROR] Unable to List connection strings for CosmosDB Account %s: %s", id.DatabaseAccountName, err)
		}
	} else {
		if model := connStringResp.Model; model != nil {
			var connStrings []string
			if model.ConnectionStrings != nil {
				connStrings = make([]string, len(*model.ConnectionStrings))
				for i, v := range *model.ConnectionStrings {
					connStrings[i] = *v.ConnectionString
					if propertyName, propertyExists := connStringPropertyMap[*v.Description]; propertyExists {
						d.Set(propertyName, v.ConnectionString) // lintignore:R001
					}
				}
			}
		}
	}
	return nil
}

func anyUnexpectedFailoverPriority(failoverPolicies []cosmosdb.FailoverPolicy) bool {
	size := len(failoverPolicies)
	for _, policy := range failoverPolicies {
		if int(*policy.FailoverPriority) > size-1 {
			return true
		}
	}
	return false
}

func flattenAzureRmCosmosDBAccountCapabilitiesAsList(capabilities *[]cosmosdb.Capability) *[]map[string]interface{} {
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

func flattenAzureRmCosmosDBAccountVirtualNetworkRulesAsList(rules *[]cosmosdb.VirtualNetworkRule) []map[string]interface{} {
	if rules == nil {
		return []map[string]interface{}{}
	}

	virtualNetworkRules := make([]map[string]interface{}, len(*rules))
	for i, r := range *rules {
		virtualNetworkRules[i] = map[string]interface{}{
			"id": *r.Id,
		}
	}
	return virtualNetworkRules
}
