// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
	client := meta.(*clients.Client).Web.AppServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFunctionAppID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	functionSettings, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(functionSettings.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		res, err := client.ListHostKeys(ctx, id.ResourceGroup, id.SiteName)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return pluginsdk.NonRetryableError(fmt.Errorf("%s was not found", id))
			}

			return pluginsdk.RetryableError(fmt.Errorf("making Read request on %s: %+v", id, err))
		}

		d.Set("primary_key", res.MasterKey)

		defaultFunctionKey := ""
		if v, ok := res.FunctionKeys["default"]; ok {
			defaultFunctionKey = *v
		}
		d.Set("default_function_key", defaultFunctionKey)

		// The name of the EventGrid System Key has changed from version 1.x to version 2.x:
		// https://learn.microsoft.com/en-us/azure/azure-functions/event-grid-how-tos?tabs=v2%2Cportal#system-key
		// This block accommodates both keys.
		eventGridExtensionConfigKey := ""
		for _, key := range []string{"eventgridextensionconfig_extension", "eventgrid_extension"} {
			if v, ok := res.SystemKeys[key]; ok {
				eventGridExtensionConfigKey = *v
				break
			}
		}
		d.Set("event_grid_extension_config_key", eventGridExtensionConfigKey)

		signalrExtensionKey := ""
		if v, ok := res.SystemKeys["signalr_extension"]; ok {
			signalrExtensionKey = *v
		}
		d.Set("signalr_extension_key", signalrExtensionKey)

		durableTaskExtensionKey := ""
		if v, ok := res.SystemKeys["durabletask_extension"]; ok {
			durableTaskExtensionKey = *v
		}
		d.Set("durabletask_extension_key", durableTaskExtensionKey)

		webPubSubExtensionKey := ""
		if v, ok := res.SystemKeys["webpubsub_extension"]; ok {
			webPubSubExtensionKey = *v
		}
		d.Set("webpubsub_extension_key", webPubSubExtensionKey)

		blobsExtensionKey := ""
		if v, ok := res.SystemKeys["blobs_extension"]; ok {
			blobsExtensionKey = *v
		}
		d.Set("blobs_extension_key", blobsExtensionKey)

		return nil
	})
}
