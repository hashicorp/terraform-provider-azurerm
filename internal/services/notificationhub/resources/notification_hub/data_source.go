// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package notification_hub

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

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Notification Hub within a Notification Hub Namespace.",
		Read:        dataSourceNotificationHubRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The Name of the Notification Hub.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"namespace_name": {
				Description: "The Name of the Notification Hub Namespace which contains the Notification Hub.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"apns_credential": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"application_mode": {
							Description: "The Application Mode which defines which server the APNS Messages should be sent to.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"bundle_id": {
							Description: "The Bundle ID of the iOS/macOS application to send push notifications for, such as `com.hashicorp.example`.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"key_id": {
							Description: "The Apple Push Notifications Service (APNS) Key.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						// Team ID (within Apple & the Portal) == "AppID" (within the API)
						"team_id": {
							Description: "The ID of the team the Token belongs to.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"token": {
							Description: "The Push Token associated with the Apple Developer Account.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
							Sensitive:   true,
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
							Description: "The API Key associated with the Google Cloud Messaging service.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceNotificationHubRead(d *pluginsdk.ResourceData, meta any) error {
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
			if setErr := d.Set("apns_credential", flattenNotificationHubsDataSourceAPNSCredentials(props.ApnsCredential)); setErr != nil {
				return fmt.Errorf("setting `apns_credential`: %+v", setErr)
			}

			if setErr := d.Set("gcm_credential", flattenNotificationHubsDataSourceGCMCredentials(props.GcmCredential)); setErr != nil {
				return fmt.Errorf("setting `gcm_credential`: %+v", setErr)
			}
		}
	}

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))

		return d.Set("tags", tags.Flatten(model.Tags))
	}

	return nil
}

func flattenNotificationHubsDataSourceAPNSCredentials(input *hubs.ApnsCredential) []any {
	if input == nil {
		return make([]any, 0)
	}

	output := make(map[string]any)

	if bundleId := input.Properties.AppName; bundleId != nil {
		output["bundle_id"] = *bundleId
	}

	applicationEndpoints := map[string]string{
		"https://api.push.apple.com:443/3/device":             "Production",
		"https://api.development.push.apple.com:443/3/device": "Sandbox",
	}
	output["application_mode"] = applicationEndpoints[input.Properties.Endpoint]

	if keyId := input.Properties.KeyId; keyId != nil {
		output["key_id"] = *keyId
	}

	if teamId := input.Properties.AppId; teamId != nil {
		output["team_id"] = *teamId
	}

	if token := input.Properties.Token; token != nil {
		output["token"] = *token
	}

	return []any{output}
}

func flattenNotificationHubsDataSourceGCMCredentials(input *hubs.GcmCredential) []any {
	if input == nil {
		return []any{}
	}

	output := make(map[string]any)
	output["api_key"] = input.Properties.GoogleApiKey

	return []any{output}
}
