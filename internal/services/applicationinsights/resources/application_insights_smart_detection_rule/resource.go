// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_smart_detection_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_smart_detection_rule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsSmartDetectionRuleUpdate,
		Read:   resourceApplicationInsightsSmartDetectionRuleRead,
		Update: resourceApplicationInsightsSmartDetectionRuleUpdate,
		Delete: resourceApplicationInsightsSmartDetectionRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := componentproactivedetectionapis.ParseProactiveDetectionConfigID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SmartDetectionRuleUpgradeV0ToV1{},
			1: migration.SmartDetectionRuleUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: applicationInsightsSmartDetectionRuleSchema(),
	}
}
