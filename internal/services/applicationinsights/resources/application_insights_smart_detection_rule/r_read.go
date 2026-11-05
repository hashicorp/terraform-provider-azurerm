// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_smart_detection_rule

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsSmartDetectionRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := componentproactivedetectionapis.ParseProactiveDetectionConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ProactiveDetectionConfigurationsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found, removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("application_insights_id", componentproactivedetectionapis.NewComponentID(id.SubscriptionId, id.ResourceGroupName, id.ComponentName).ID())

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
		d.Set("enabled", model.Enabled)
		d.Set("send_emails_to_subscription_owners", model.SendEmailsToSubscriptionOwners)
		d.Set("additional_email_recipients", pluginsdk.FlattenSlice(model.CustomEmails))
	}
	return nil
}

// The Smart Detection Rule name from the UI doesn't match what the API accepts.
// This Diff checks if the old and new name match when converted to the API version of the name
func smartDetectionRuleNameDiff(_, old string, new string, _ *pluginsdk.ResourceData) bool {
	apiNew := convertUiNameToApiName(new)

	return strings.EqualFold(old, apiNew)
}

func convertUiNameToApiName(uiName any) string {
	apiName := uiName.(string)
	switch uiName.(string) {
	case "Slow page load time":
		apiName = "slowpageloadtime"
	case "Slow server response time":
		apiName = "slowserverresponsetime"
	case "Long dependency duration":
		apiName = "longdependencyduration"
	case "Degradation in server response time":
		apiName = "degradationinserverresponsetime"
	case "Degradation in dependency duration":
		apiName = "degradationindependencyduration"
	case "Degradation in trace severity ratio":
		apiName = "extension_traceseveritydetector"
	case "Abnormal rise in exception volume":
		apiName = "extension_exceptionchangeextension"
	case "Potential memory leak detected":
		apiName = "extension_memoryleakextension"
	case "Potential security issue detected":
		apiName = "extension_securityextensionspackage"
	case "Abnormal rise in daily data volume":
		apiName = "extension_billingdatavolumedailyspikeextension"
	}
	return apiName
}
