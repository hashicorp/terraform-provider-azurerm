// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing EventHub Namespace.",
		Read:        EventHubNamespaceDataSourceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the EventHub Namespace.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"default_primary_connection_string_alias": {
				Description: "The alias of the primary connection string for the authorization rule `RootManageSharedAccessKey`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"default_secondary_connection_string_alias": {
				Description: "The alias of the secondary connection string for the authorization rule `RootManageSharedAccessKey`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"auto_inflate_enabled": {
				Description: "Whether Auto Inflate is enabled for the EventHub Namespace.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"dedicated_cluster_id": {
				Description: "The ID of the EventHub Dedicated Cluster where this Namespace exists.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"capacity": {
				Description: "The Capacity / Throughput Units for a `Standard` SKU namespace.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"kafka_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"maximum_throughput_units": {
				Description: "The maximum number of throughput units when Auto Inflate is Enabled.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"default_primary_connection_string": {
				Description: "The primary connection string for the authorization rule `RootManageSharedAccessKey`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"default_primary_key": {
				Description: "The primary access key for the authorization rule `RootManageSharedAccessKey`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"default_secondary_connection_string": {
				Description: "The secondary connection string for the authorization rule `RootManageSharedAccessKey`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"default_secondary_key": {
				Description: "The secondary access key for the authorization rule `RootManageSharedAccessKey`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"sku": {
				Description: "Which tier to use.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func EventHubNamespaceDataSourceRead(d *pluginsdk.ResourceData, meta any) error {
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

	d.Set("name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

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
			d.Set("dedicated_cluster_id", props.ClusterArmId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	defaultRuleId := authorizationrulesnamespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, eventHubNamespaceDefaultAuthorizationRule)
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
