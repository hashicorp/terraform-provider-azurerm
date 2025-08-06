// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package notificationhub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceNotificationHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNotificationHubRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"namespace_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"apns_credential": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"application_mode": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"bundle_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"key_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						// Team ID (within Apple & the Portal) == "AppID" (within the API)
						"team_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"token": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"gcm_credential": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"api_key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceNotificationHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := hubs.NewNotificationHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))
	resp, err := client.NotificationHubsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	credentials, err := client.NotificationHubsGetPnsCredentials(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving credentials for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.NotificationHubName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if credentialsModel := credentials.Model; credentialsModel != nil {
		if props := credentialsModel.Properties; props != nil {
			apns := flattenNotificationHubsDataSourceAPNSCredentials(props.ApnsCredential)
			if setErr := d.Set("apns_credential", apns); setErr != nil {
				return fmt.Errorf("setting `apns_credential`: %+v", err)
			}

			gcm := flattenNotificationHubsDataSourceGCMCredentials(props.GcmCredential)
			if setErr := d.Set("gcm_credential", gcm); setErr != nil {
				return fmt.Errorf("setting `gcm_credential`: %+v", err)
			}
		}
	}

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))

		return d.Set("tags", tags.Flatten(model.Tags))
	}

	return nil
}

func flattenNotificationHubsDataSourceAPNSCredentials(input *hubs.ApnsCredential) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	output := make(map[string]interface{})

	if bundleId := input.Properties.AppName; bundleId != nil {
		output["bundle_id"] = *bundleId
	}

	applicationEndpoints := map[string]string{
		"https://api.push.apple.com:443/3/device":             "Production",
		"https://api.development.push.apple.com:443/3/device": "Sandbox",
	}
	applicationMode := applicationEndpoints[input.Properties.Endpoint]
	output["application_mode"] = applicationMode

	if keyId := input.Properties.KeyId; keyId != nil {
		output["key_id"] = *keyId
	}

	if teamId := input.Properties.AppId; teamId != nil {
		output["team_id"] = *teamId
	}

	if token := input.Properties.Token; token != nil {
		output["token"] = *token
	}

	return []interface{}{output}
}

func flattenNotificationHubsDataSourceGCMCredentials(input *hubs.GcmCredential) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	output["api_key"] = input.Properties.GoogleApiKey

	return []interface{}{output}
}
