// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_scheduled

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Sentinel Scheduled Alert Rule.",
		Create:      resourceSentinelAlertRuleScheduledCreateUpdate,
		Read:        resourceSentinelAlertRuleScheduledRead,
		Update:      resourceSentinelAlertRuleScheduledCreateUpdate,
		Delete:      resourceSentinelAlertRuleScheduledDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := alertrules.ParseAlertRuleID(id)
			return err
		}, helpers.ImportSentinelAlertRule(alertrules.AlertRuleKindScheduled)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: sentinelAlertRuleScheduledSchema(),
	}
}
