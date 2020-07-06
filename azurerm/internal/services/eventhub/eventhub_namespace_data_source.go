package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceEventHubNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEventHubNamespaceRead,

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

			"default_primary_connection_string_alias": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string_alias": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"auto_inflate_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"dedicated_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"kafka_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"maximum_throughput_units": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"default_primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceEventHubNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: EventHub Namespace %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on EventHub Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("sku", string(resp.Sku.Name))
	d.Set("capacity", resp.Sku.Capacity)

	keys, err := client.ListKeys(ctx, resourceGroup, name, eventHubNamespaceDefaultAuthorizationRule)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for EventHub Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
	} else {
		d.Set("default_primary_connection_string_alias", keys.AliasPrimaryConnectionString)
		d.Set("default_secondary_connection_string_alias", keys.AliasSecondaryConnectionString)
		d.Set("default_primary_connection_string", keys.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", keys.SecondaryConnectionString)
		d.Set("default_primary_key", keys.PrimaryKey)
		d.Set("default_secondary_key", keys.SecondaryKey)
	}

	if props := resp.EHNamespaceProperties; props != nil {
		d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
		d.Set("kafka_enabled", props.KafkaEnabled)
		d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
		d.Set("zone_redundant", props.ZoneRedundant)
		d.Set("dedicated_cluster_id", props.ClusterArmID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
