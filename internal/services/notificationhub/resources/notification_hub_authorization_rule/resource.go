// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub_authorization_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/resources/notification_hub_authorization_rule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var notificationHubResourceName = "azurerm_notification_hub"

var notificationHubNamespaceResourceName = "azurerm_notification_hub_namespace"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNotificationHubAuthorizationRuleCreateUpdate,
		Read:   resourceNotificationHubAuthorizationRuleRead,
		Update: resourceNotificationHubAuthorizationRuleCreateUpdate,
		Delete: resourceNotificationHubAuthorizationRuleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := hubs.ParseNotificationHubAuthorizationRuleID(id)
			return err
		}),
		// TODO: customizeDiff for send+listen when manage selected

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NotificationHubAuthorizationRuleResourceV0ToV1{},
		}),

		Schema: notificationHubAuthorizationRuleSchema(),
	}
}
