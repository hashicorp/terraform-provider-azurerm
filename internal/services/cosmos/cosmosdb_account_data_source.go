// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

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

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_free_tier": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			// TODO 4.0: change this from enable_* to *_enabled
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

			// TODO 4.0: change this from enable_* to *_enabled
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
		},
	}
}

func dataSourceCosmosDbAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.DatabaseClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDatabaseAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("kind", string(resp.Kind))

	if props := resp.DatabaseAccountGetProperties; props != nil {
		d.Set("offer_type", string(props.DatabaseAccountOfferType))
		d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilterThreePointOh(props.IPRules))
		d.Set("endpoint", props.DocumentEndpoint)
		d.Set("is_virtual_network_filter_enabled", resp.IsVirtualNetworkFilterEnabled)
		d.Set("enable_free_tier", resp.EnableFreeTier)
		d.Set("enable_automatic_failover", resp.EnableAutomaticFailover)

		if v := props.KeyVaultKeyURI; v != nil {
			d.Set("key_vault_key_id", resp.KeyVaultKeyURI)
		}

		if err = d.Set("consistency_policy", flattenAzureRmCosmosDBAccountConsistencyPolicy(resp.ConsistencyPolicy)); err != nil {
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
					"id":                *l.ID,
					"location":          location.NormalizeNilable(l.LocationName),
					"failover_priority": int(*l.FailoverPriority),
				}
			}
		} else {
			for _, l := range *props.FailoverPolicies {
				locations[*l.FailoverPriority] = map[string]interface{}{
					"id":                *l.ID,
					"location":          location.NormalizeNilable(l.LocationName),
					"failover_priority": int(*l.FailoverPriority),
				}
			}
		}
		if err = d.Set("geo_location", locations); err != nil {
			return fmt.Errorf("setting `geo_location`: %+v", err)
		}

		if err = d.Set("capabilities", flattenAzureRmCosmosDBAccountCapabilitiesAsList(resp.Capabilities)); err != nil {
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

		d.Set("enable_multiple_write_locations", resp.EnableMultipleWriteLocations)
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Write keys for %s: %s", id, err)
	} else {
		d.Set("primary_key", keys.PrimaryMasterKey)
		d.Set("secondary_key", keys.SecondaryMasterKey)
	}

	readonlyKeys, err := client.ListReadOnlyKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		log.Printf("[ERROR] Unable to List read-only keys for %s: %s", id, err)
	} else {
		d.Set("primary_readonly_key", readonlyKeys.PrimaryReadonlyMasterKey)
		d.Set("secondary_readonly_key", readonlyKeys.SecondaryReadonlyMasterKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func anyUnexpectedFailoverPriority(failoverPolicies []documentdb.FailoverPolicy) bool {
	size := len(failoverPolicies)
	for _, policy := range failoverPolicies {
		if int(*policy.FailoverPriority) > size-1 {
			return true
		}
	}
	return false
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
