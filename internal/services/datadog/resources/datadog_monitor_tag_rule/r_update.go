// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	datadogClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogTagRulesUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := datadogClient.From(meta).TagRules
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := tagrules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	payload := tagrules.MonitoringTagRules{
		Properties: &tagrules.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]any)),
			MetricRules: expandMetricRules(d.Get("metric").([]any)),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDatadogTagRulesRead(d, meta)
}
