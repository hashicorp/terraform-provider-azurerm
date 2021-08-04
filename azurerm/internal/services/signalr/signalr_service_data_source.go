package signalr

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/sdk/signalr"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmSignalRService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmSignalRServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_port": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"server_port": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmSignalRServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := signalr.NewSignalRID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	keys, err := client.ListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.SignalRName)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("hostname", props.HostName)
			d.Set("ip_address", props.ExternalIP)
			d.Set("public_port", props.PublicPort)
			d.Set("server_port", props.ServerPort)
		}

		if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	if model := keys.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_access_key", model.SecondaryKey)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
	}

	return nil
}
