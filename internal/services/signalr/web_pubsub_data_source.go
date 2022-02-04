package signalr

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceWebPubsub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceWebPubsubRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"capacity": {
				Type:     pluginsdk.TypeInt,
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

			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"aad_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tls_client_cert_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"external_ip": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceWebPubsubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewWebPubsubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Web Pubsub %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("listing keys for %q: %+v", id, err)
	}

	d.Set("primary_access_key", keys.PrimaryKey)
	d.Set("primary_connection_string", keys.PrimaryConnectionString)
	d.Set("secondary_access_key", keys.SecondaryKey)
	d.Set("secondary_connection_string", keys.SecondaryConnectionString)

	d.SetId(id.ID())

	d.Set("name", id.WebPubSubName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		if sku.Name != nil {
			d.Set("sku", sku.Name)
		}
		if sku.Capacity != nil {
			d.Set("capacity", sku.Capacity)
		}
	}

	if props := resp.Properties; props != nil {
		d.Set("external_ip", props.ExternalIP)
		d.Set("hostname", props.HostName)
		d.Set("public_port", props.PublicPort)
		d.Set("server_port", props.ServerPort)
		d.Set("version", props.Version)
		if props.DisableAadAuth != nil {
			d.Set("aad_auth_enabled", !(*props.DisableAadAuth))
		}
		if props.DisableLocalAuth != nil {
			d.Set("local_auth_enabled", !(*props.DisableLocalAuth))
		}
		if props.PublicNetworkAccess != nil {
			d.Set("public_network_access_enabled", strings.EqualFold(*props.PublicNetworkAccess, "Enabled"))
		}
		if props.TLS != nil {
			d.Set("tls_client_cert_enabled", props.TLS.ClientCertEnabled)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
