// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func notificationHubSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name to use for this Notification Hub.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"namespace_name": {
			Description: "The name of the Notification Hub Namespace in which to create this Notification Hub.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"apns_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// NOTE: APNS supports two modes, certificate auth (v1) and token auth (v2)
					// certificate authentication/v1 is marked for deprecation; as such we're not
					// supporting it at this time.
					"application_mode": {
						Description: "The Application Mode which defines which server the APNS Messages should be sent to.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							apnsProductionName,
							apnsSandboxName,
						}, false),
					},
					"bundle_id": {
						Description: "The Bundle ID of the iOS/macOS application to send push notifications for, such as `com.hashicorp.example`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"key_id": {
						Description: "The Apple Push Notifications Service (APNS) Key.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					// Team ID (within Apple & the Portal) == "AppID" (within the API)
					"team_id": {
						Description: "The ID of the team the Token belongs to.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"token": {
						Description: "The Push Token associated with the Apple Developer Account. This is the contents of the `key` downloaded from [the Apple Developer Portal](https://developer.apple.com/account/ios/authkey/) between the `-----BEGIN PRIVATE KEY-----` and `-----END PRIVATE KEY-----` blocks.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"browser_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subject": {
						Description:  "The subject name of web push.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"vapid_private_key": {
						Description:  "The Voluntary Application Server Identification (VAPID) private key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"vapid_public_key": {
						Description:  "The Voluntary Application Server Identification (VAPID) public key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"gcm_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"api_key": {
						Description: "The API Key associated with the Google Cloud Messaging service.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
