// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogTagRulesCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.TagRules
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := tagrules.ParseMonitorID(d.Get("datadog_monitor_id").(string))
	if err != nil {
		return err
	}

	id := tagrules.NewTagRuleID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for an existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) && !isDefaultSettings(existing.Model) {
			return tf.ImportAsExistsError("azurerm_datadog_monitor_tag_rule", id.ID())
		}
	}

	payload := tagrules.MonitoringTagRules{
		Properties: &tagrules.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]any)),
			MetricRules: expandMetricRules(d.Get("metric").([]any)),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatadogTagRulesRead(d, meta)
}
