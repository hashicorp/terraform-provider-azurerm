// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_webhook

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_webhook/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationWebhookCreate,
		Read:   resourceAutomationWebhookRead,
		Update: resourceAutomationWebhookUpdate,
		Delete: resourceAutomationWebhookDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webhook.ParseWebHookID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AutomationWebhookV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: automationWebhookSchema(),
	}
}
