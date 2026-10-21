// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogTagRulesUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.Rules
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	payload := rules.MonitoringTagRules{
		Properties: &rules.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]any)),
			MetricRules: expandMetricRules(d.Get("metric").([]any)),
		},
	}
	if _, err := client.TagRulesCreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDatadogTagRulesRead(d, meta)
}
