// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceFunctionAppHostKeys() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceFunctionAppHostKeysRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_function_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"event_grid_extension_config_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"event_grid_extension_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"signalr_extension_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"durabletask_extension_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"webpubsub_extension_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"blobs_extension_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceFunctionAppHostKeysRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewAppServiceID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	functionSettings, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(functionSettings.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		resp, err := client.ListHostKeys(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return pluginsdk.NonRetryableError(fmt.Errorf("%s was not found", id))
			}

			return pluginsdk.RetryableError(fmt.Errorf("retrieving %s: %+v", id, err))
		}

		if resp.Model == nil {
			return pluginsdk.RetryableError(fmt.Errorf("retrieving %s: model was nil", id))
		}

		d.Set("primary_key", resp.Model.MasterKey)

		defaultFunctionKey := ""
		if v, ok := pointer.From(resp.Model.FunctionKeys)["default"]; ok {
			defaultFunctionKey = v
		}
		d.Set("default_function_key", defaultFunctionKey)

		// The name of the EventGrid System Key has changed from version 1.x to version 2.x:
		// https://learn.microsoft.com/en-us/azure/azure-functions/event-grid-how-tos?tabs=v2%2Cportal#system-key
		// This block accommodates both keys.
		eventGridExtensionConfigKey := ""
		for _, key := range []string{"eventgridextensionconfig_extension", "eventgrid_extension"} {
			if v, ok := pointer.From(resp.Model.SystemKeys)[key]; ok {
				eventGridExtensionConfigKey = v
				break
			}
		}
		d.Set("event_grid_extension_config_key", eventGridExtensionConfigKey)

		signalrExtensionKey := ""
		if v, ok := pointer.From(resp.Model.SystemKeys)["signalr_extension"]; ok {
			signalrExtensionKey = v
		}
		d.Set("signalr_extension_key", signalrExtensionKey)

		durableTaskExtensionKey := ""
		if v, ok := pointer.From(resp.Model.SystemKeys)["durabletask_extension"]; ok {
			durableTaskExtensionKey = v
		}
		d.Set("durabletask_extension_key", durableTaskExtensionKey)

		webPubSubExtensionKey := ""
		if v, ok := pointer.From(resp.Model.SystemKeys)["webpubsub_extension"]; ok {
			webPubSubExtensionKey = v
		}
		d.Set("webpubsub_extension_key", webPubSubExtensionKey)

		blobsExtensionKey := ""
		if v, ok := pointer.From(resp.Model.SystemKeys)["blobs_extension"]; ok {
			blobsExtensionKey = v
		}
		d.Set("blobs_extension_key", blobsExtensionKey)

		return nil
	})
}
