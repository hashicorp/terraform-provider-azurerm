// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_fusion

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleFusionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := sentinelClient.From(meta).AlertRulesClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload, ok := resp.Model.(alertrules.FusionAlertRule)
	if !ok {
		return fmt.Errorf("retrieving %s: expected an alert rule of type `Fusion`, got %q", id, pointer.From(resp.Model.AlertRule().Type))
	}

	if payload.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if d.HasChange("alert_rule_template_guid") {
		payload.Properties.AlertRuleTemplateName = d.Get("alert_rule_template_guid").(string)
	}

	if d.HasChange("enabled") {
		payload.Properties.Enabled = d.Get("enabled").(bool)
	}

	if d.HasChange("source") {
		payload.Properties.SourceSettings = expandFusionSourceSettings(d.Get("source").([]any))
	}

	// The `Description` is read-only but not specified on the Swagger, tracked on: https://github.com/Azure/azure-rest-api-specs/issues/31330
	payload.Properties.Description = nil
	payload.Properties.DisplayName = nil
	payload.Properties.LastModifiedUtc = nil
	payload.Properties.Severity = nil
	payload.Properties.Tactics = nil

	if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceSentinelAlertRuleFusionRead(d, meta)
}
