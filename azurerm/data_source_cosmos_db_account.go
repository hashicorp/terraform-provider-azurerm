package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmCosmosDBAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmCosmosDBAccountRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"tags": tagsForDataSourceSchema(),

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

			"primary_master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_readonly_master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_readonly_master_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceArmCosmosDBAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cosmosDBClient
	ctx := meta.(*ArmClient).StopContext

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
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("kind", string(resp.Kind))
	flattenAndSetTags(d, resp.Tags)

	if props := resp.DatabaseAccountProperties; props != nil {
		d.Set("offer_type", string(props.DatabaseAccountOfferType))
		d.Set("ip_range_filter", props.IPRangeFilter)
		d.Set("endpoint", props.DocumentEndpoint)
		d.Set("enable_automatic_failover", resp.EnableAutomaticFailover)

		if err := d.Set("consistency_policy", flattenAzureRmCosmosDBAccountConsistencyPolicy(resp.ConsistencyPolicy)); err != nil {
			return fmt.Errorf("Error setting `consistency_policy`: %+v", err)
		}

		//sort `geo_locations` by fail over priority
		locations := make([]map[string]interface{}, len(*props.FailoverPolicies))
		for _, l := range *props.FailoverPolicies {
			locations[*l.FailoverPriority] = map[string]interface{}{
				"id":                *l.ID,
				"location":          azureRMNormalizeLocation(*l.LocationName),
				"failover_priority": int(*l.FailoverPriority),
			}
		}
		if err := d.Set("geo_location", locations); err != nil {
			return fmt.Errorf("Error setting `geo_location`: %+v", err)
		}

		readEndpoints := make([]string, 0)
		if locations := props.ReadLocations; locations != nil {
			for _, l := range *locations {
				readEndpoints = append(readEndpoints, *l.DocumentEndpoint)
			}
		}
		d.Set("read_endpoints", readEndpoints)

		writeEndpoints := make([]string, 0)
		if locations := props.WriteLocations; locations != nil {
			for _, l := range *locations {
				writeEndpoints = append(writeEndpoints, *l.DocumentEndpoint)
			}
		}
		d.Set("write_endpoints", writeEndpoints)
	}

	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Write keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_master_key", keys.PrimaryMasterKey)
		d.Set("secondary_master_key", keys.SecondaryMasterKey)
	}

	readonlyKeys, err := client.ListReadOnlyKeys(ctx, resourceGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List read-only keys for CosmosDB Account %s: %s", name, err)
	} else {
		d.Set("primary_readonly_master_key", readonlyKeys.PrimaryReadonlyMasterKey)
		d.Set("secondary_readonly_master_key", readonlyKeys.SecondaryReadonlyMasterKey)
	}

	return nil
}
