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

			"location": locationForDataSourceSchema(),

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"tags": tagsForDataSourceSchema(),

			//resource fields
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
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_master_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_readonly_master_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_readonly_master_key": {
				Type:     schema.TypeString,
				Computed: true,
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
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM CosmosDB Account '%s': %s", name, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resourceGroup)
	flattenAndSetTags(d, resp.Tags)

	d.Set("kind", string(resp.Kind))
	d.Set("offer_type", string(resp.DatabaseAccountOfferType))
	d.Set("ip_range_filter", resp.IPRangeFilter)
	d.Set("endpoint", resp.DocumentEndpoint)

	if v := resp.EnableAutomaticFailover; v != nil {
		d.Set("enable_automatic_failover", resp.EnableAutomaticFailover)
	}

	flattenAndSetAzureRmCosmosDBAccountConsistencyPolicy(d, resp.ConsistencyPolicy)

	//sort `geo_locations` by fail over priority
	locations := make([]map[string]interface{}, len(*resp.FailoverPolicies))
	for _, l := range *resp.FailoverPolicies {
		locations[*l.FailoverPriority] = map[string]interface{}{
			"id":                *l.ID,
			"location":          azureRMNormalizeLocation(*l.LocationName),
			"failover_priority": int(*l.FailoverPriority),
		}
	}
	d.Set("geo_location", &locations)

	readEndpoints := []string{}
	for _, l := range *resp.ReadLocations {
		readEndpoints = append(readEndpoints, *l.DocumentEndpoint)
	}
	d.Set("read_endpoints", readEndpoints)

	writeEndpoints := []string{}
	for _, l := range *resp.WriteLocations {
		writeEndpoints = append(writeEndpoints, *l.DocumentEndpoint)
	}
	d.Set("write_endpoints", writeEndpoints)

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
