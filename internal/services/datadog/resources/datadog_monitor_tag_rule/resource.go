// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatadogTagRulesCreate,
		Read:   resourceDatadogTagRulesRead,
		Update: resourceDatadogTagRulesUpdate,
		Delete: resourceDatadogTagRulesDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := tagrules.ParseTagRuleID(id)
			return err
		}),

		Schema: datadogMonitorTagRuleSchema(),
	}
}
