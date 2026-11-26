// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogTagRulesDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.TagRules
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	// Tag Rules can't be removed on their own, they can only be nil'd out
	payload := tagrules.MonitoringTagRules{
		Properties: &tagrules.MonitoringTagRulesProperties{
			LogRules: &tagrules.LogRules{
				SendAadLogs:          pointer.To(false),
				SendSubscriptionLogs: pointer.To(false),
				SendResourceLogs:     pointer.To(false),
				FilteringTags:        &[]tagrules.FilteringTag{},
			},
			MetricRules: &tagrules.MetricRules{
				FilteringTags: &[]tagrules.FilteringTag{},
			},
		},
	}
	if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("removing %s: %+v", *id, err)
	}

	return nil
}
