// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
)

func expandNotificationHubsAPNSCredentials(inputs []any) *hubs.ApnsCredential {
	if len(inputs) == 0 {
		return nil
	}

	input := inputs[0].(map[string]any)
	applicationMode := input["application_mode"].(string)
	bundleId := input["bundle_id"].(string)
	keyId := input["key_id"].(string)
	teamId := input["team_id"].(string)
	token := input["token"].(string)

	applicationEndpoints := map[string]string{
		apnsProductionName: apnsProductionEndpoint,
		apnsSandboxName:    apnsSandboxEndpoint,
	}
	endpoint := applicationEndpoints[applicationMode]

	credentials := hubs.ApnsCredential{
		Properties: hubs.ApnsCredentialProperties{
			AppId:    pointer.To(teamId),
			AppName:  pointer.To(bundleId),
			Endpoint: endpoint,
			KeyId:    pointer.To(keyId),
			Token:    pointer.To(token),
		},
	}
	return &credentials
}

func expandNotificationHubsBrowserCredentials(inputs []any) *hubs.BrowserCredential {
	if len(inputs) == 0 {
		return nil
	}

	input := inputs[0].(map[string]any)
	credentials := hubs.BrowserCredential{
		Properties: hubs.BrowserCredentialProperties{
			Subject:         input["subject"].(string),
			VapidPrivateKey: input["vapid_private_key"].(string),
			VapidPublicKey:  input["vapid_public_key"].(string),
		},
	}
	return &credentials
}

func expandNotificationHubsGCMCredentials(inputs []any) *hubs.GcmCredential {
	if len(inputs) == 0 {
		return nil
	}

	input := inputs[0].(map[string]any)
	apiKey := input["api_key"].(string)
	credentials := hubs.GcmCredential{
		Properties: hubs.GcmCredentialProperties{
			GoogleApiKey: apiKey,
		},
	}
	return &credentials
}
