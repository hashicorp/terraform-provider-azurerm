// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceWebPubsubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webpubsub.NewWebPubSubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	keys, err := client.ListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("listing keys for %q: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.WebPubSubName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		skuName := ""
		skuCapacity := int64(0)
		if model.Sku != nil {
			skuName = model.Sku.Name
			skuCapacity = *model.Sku.Capacity
		}
		d.Set("sku", skuName)
		d.Set("capacity", skuCapacity)

		if props := model.Properties; props != nil {
			d.Set("external_ip", props.ExternalIP)
			d.Set("hostname", props.HostName)
			d.Set("public_port", props.PublicPort)
			d.Set("server_port", props.ServerPort)
			d.Set("version", props.Version)
			aadAuthEnabled := true
			if props.DisableAadAuth != nil {
				aadAuthEnabled = !(*props.DisableAadAuth)
			}
			d.Set("aad_auth_enabled", aadAuthEnabled)

			disableLocalAuth := false
			if props.DisableLocalAuth != nil {
				disableLocalAuth = !(*props.DisableLocalAuth)
			}
			d.Set("local_auth_enabled", disableLocalAuth)

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccessEnabled = strings.EqualFold(*props.PublicNetworkAccess, "Enabled")
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			tlsClientCertEnabled := false
			if props.Tls != nil && props.Tls.ClientCertEnabled != nil {
				tlsClientCertEnabled = *props.Tls.ClientCertEnabled
			}
			d.Set("tls_client_cert_enabled", tlsClientCertEnabled)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
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
