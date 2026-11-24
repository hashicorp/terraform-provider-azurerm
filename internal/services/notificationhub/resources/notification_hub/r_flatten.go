// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
)

func flattenNotificationHubsAPNSCredentials(input *hubs.ApnsCredential) []any {
	if input == nil {
		return make([]any, 0)
	}

	output := make(map[string]any)

	if bundleId := input.Properties.AppName; bundleId != nil {
		output["bundle_id"] = *bundleId
	}

	applicationEndpoints := map[string]string{
		apnsProductionEndpoint: apnsProductionName,
		apnsSandboxEndpoint:    apnsSandboxName,
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

func flattenNotificationHubsBrowserCredentials(input *hubs.BrowserCredential) []any {
	if input == nil {
		return make([]any, 0)
	}

	output := make(map[string]any)

	output["subject"] = input.Properties.Subject
	output["vapid_private_key"] = input.Properties.VapidPrivateKey
	output["vapid_public_key"] = input.Properties.VapidPublicKey

	return []any{output}
}

func flattenNotificationHubsGCMCredentials(input *hubs.GcmCredential) []any {
	if input == nil {
		return []any{}
	}

	output := make(map[string]any)
	output["api_key"] = input.Properties.GoogleApiKey

	return []any{output}
}
