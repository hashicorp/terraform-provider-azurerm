// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/resources/notification_hub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNotificationHubCreateUpdate,
		Read:   resourceNotificationHubRead,
		Update: resourceNotificationHubCreateUpdate,
		Delete: resourceNotificationHubDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := hubs.ParseNotificationHubID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NotificationHubResourceV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
			// NOTE: the ForceNew is to workaround a bug in the Azure SDK where nil-values aren't sent to the API.
			// Bug: https://github.com/Azure/azure-sdk-for-go/issues/2246

			oAPNS, nAPNS := diff.GetChange("apns_credential.#")
			oAPNSi := oAPNS.(int)
			nAPNSi := nAPNS.(int)
			if nAPNSi < oAPNSi {
				diff.ForceNew("apns_credential")
			}

			oGCM, nGCM := diff.GetChange("gcm_credential.#")
			oGCMi := oGCM.(int)
			nGCMi := nGCM.(int)
			if nGCMi < oGCMi {
				diff.ForceNew("gcm_credential")
			}

			return nil
		}),

		Schema: notificationHubSchema(),
	}
}
