// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNotificationHubRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hubs.ParseNotificationHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.NotificationHubsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	credentials, err := client.NotificationHubsGetPnsCredentials(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving credentials for %s: %+v", *id, err)
	}

	d.Set("name", id.NotificationHubName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if credentialsModel := credentials.Model; credentialsModel != nil {
		if props := credentialsModel.Properties; props != nil {
			if setErr := d.Set("apns_credential", flattenNotificationHubsAPNSCredentials(props.ApnsCredential)); setErr != nil {
				return fmt.Errorf("setting `apns_credential`: %+v", setErr)
			}
			if setErr := d.Set("browser_credential", flattenNotificationHubsBrowserCredentials(props.BrowserCredential)); setErr != nil {
				return fmt.Errorf("setting `browser_credential`: %+v", setErr)
			}
			if setErr := d.Set("gcm_credential", flattenNotificationHubsGCMCredentials(props.GcmCredential)); setErr != nil {
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

const (
	apnsProductionName     = "Production"
	apnsProductionEndpoint = "https://api.push.apple.com:443/3/device"
	apnsSandboxName        = "Sandbox"
	apnsSandboxEndpoint    = "https://api.development.push.apple.com:443/3/device"
)
