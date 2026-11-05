// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.ComponentsClient
	ruleClient := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := componentsapis.ParseComponentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ComponentsDelete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	// if disable_generated_rule=true, the generated rule is not automatically deleted.
	if meta.(*clients.Client).Features.ApplicationInsights.DisableGeneratedRule {
		ruleName := fmt.Sprintf("Failure Anomalies - %s", id.ComponentName)
		ruleId := smartdetectoralertrules.NewSmartDetectorAlertRuleID(id.SubscriptionId, id.ResourceGroupName, ruleName)
		deleteResp, deleteErr := ruleClient.Delete(ctx, ruleId)
		if deleteErr != nil && !response.WasNotFound(deleteResp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", ruleId, deleteErr)
		}
	}

	return err
}
