package notificationhub

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceNotificationHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Notification Hub %q was not found in Namespace %q / Resource Group %q", name, namespaceName, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
	}

	credentials, err := client.GetPnsCredentials(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Credentials for Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := credentials.PnsCredentialsProperties; props != nil {
		apns := flattenNotificationHubsDataSourceAPNSCredentials(props.ApnsCredential)
		if setErr := d.Set("apns_credential", apns); setErr != nil {
			return fmt.Errorf("Error setting `apns_credential`: %+v", err)
		}

		gcm := flattenNotificationHubsDataSourceGCMCredentials(props.GcmCredential)
		if setErr := d.Set("gcm_credential", gcm); setErr != nil {
			return fmt.Errorf("Error setting `gcm_credential`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenNotificationHubsDataSourceAPNSCredentials(input *notificationhubs.ApnsCredential) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	output := make(map[string]interface{})

	if bundleId := input.AppName; bundleId != nil {
		output["bundle_id"] = *bundleId
	}

	if endpoint := input.Endpoint; endpoint != nil {
		applicationEndpoints := map[string]string{
			"https://api.push.apple.com:443/3/device":             "Production",
			"https://api.development.push.apple.com:443/3/device": "Sandbox",
		}
		applicationMode := applicationEndpoints[*endpoint]
		output["application_mode"] = applicationMode
	}

	if keyId := input.KeyID; keyId != nil {
		output["key_id"] = *keyId
	}

	if teamId := input.AppID; teamId != nil {
		output["team_id"] = *teamId
	}

	if token := input.Token; token != nil {
		output["token"] = *token
	}

	return []interface{}{output}
}

func flattenNotificationHubsDataSourceGCMCredentials(input *notificationhubs.GcmCredential) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	if props := input.GcmCredentialProperties; props != nil {
		if apiKey := props.GoogleAPIKey; apiKey != nil {
			output["api_key"] = *apiKey
		}
	}

	return []interface{}{output}
}
