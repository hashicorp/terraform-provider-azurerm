package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/authorizationrulesnamespaces"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/namespaces"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func EventHubNamespaceDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: EventHubNamespaceDataSourceRead,

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

			"default_primary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"auto_inflate_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"dedicated_cluster_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"kafka_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"maximum_throughput_units": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"default_primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func EventHubNamespaceDataSourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	authorizationRulesClient := meta.(*clients.Client).Eventhub.NamespaceAuthorizationRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if sku := model.Sku; sku != nil {
			d.Set("sku", string(sku.Name))
			d.Set("capacity", sku.Capacity)
		}

		if props := model.Properties; props != nil {
			d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
			d.Set("kafka_enabled", props.KafkaEnabled)
			d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
			d.Set("zone_redundant", props.ZoneRedundant)
			d.Set("dedicated_cluster_id", props.ClusterArmId)
		}

		if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	defaultRuleId := authorizationrulesnamespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroup, id.Name, eventHubNamespaceDefaultAuthorizationRule)
	keys, err := authorizationRulesClient.NamespacesListKeys(ctx, defaultRuleId)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for %s: %+v", id, err)
	}
	if model := keys.Model; model != nil {
		d.Set("default_primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("default_secondary_connection_string_alias", model.AliasSecondaryConnectionString)
		d.Set("default_primary_connection_string", model.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", model.SecondaryConnectionString)
		d.Set("default_primary_key", model.PrimaryKey)
		d.Set("default_secondary_key", model.SecondaryKey)
	}

	return nil
}
