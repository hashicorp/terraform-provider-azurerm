// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_smart_detection_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsSmartDetectionRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := componentproactivedetectionapis.ParseProactiveDetectionConfigID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] resetting AzureRM Application Insights Smart Detection Rule %s", id)

	resp, err := client.ProactiveDetectionConfigurationsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found, removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("model was nil for %s", id)
	}
	smartDetectionRuleProperties := componentproactivedetectionapis.ApplicationInsightsComponentProactiveDetectionConfiguration{
		Name:                           pointer.To(id.ConfigurationId),
		Enabled:                        resp.Model.RuleDefinitions.IsEnabledByDefault,
		SendEmailsToSubscriptionOwners: resp.Model.RuleDefinitions.SupportsEmailNotifications,
		CustomEmails:                   pluginsdk.ExpandStringSlice([]any{}),
	}

	// Application Insights defaults all the Smart Detection Rules so if a user wants to delete a rule, we'll update it back to it's default values.
	if _, err = client.ProactiveDetectionConfigurationsUpdate(ctx, *id, smartDetectionRuleProperties); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("resetting %s: %+v", id, err)
	}

	return nil
}
