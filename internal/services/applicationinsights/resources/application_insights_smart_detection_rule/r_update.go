// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_smart_detection_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsSmartDetectionRuleUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// The Smart Detection Rule name from the UI doesn't match what the API accepts.
	// We'll have the user submit what the name looks like in the UI and convert it behind the scenes to match what the API accepts
	name := convertUiNameToApiName(d.Get("name"))

	appInsightsId, err := componentproactivedetectionapis.ParseComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	id := componentproactivedetectionapis.NewProactiveDetectionConfigID(appInsightsId.SubscriptionId, appInsightsId.ResourceGroupName, appInsightsId.ComponentName, name)

	smartDetectionRuleProperties := componentproactivedetectionapis.ApplicationInsightsComponentProactiveDetectionConfiguration{
		Name:                           &name,
		Enabled:                        pointer.To(d.Get("enabled").(bool)),
		SendEmailsToSubscriptionOwners: pointer.To(d.Get("send_emails_to_subscription_owners").(bool)),
		CustomEmails:                   pluginsdk.ExpandStringSlice(d.Get("additional_email_recipients").(*pluginsdk.Set).List()),
	}

	if _, err = client.ProactiveDetectionConfigurationsUpdate(ctx, id, smartDetectionRuleProperties); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApplicationInsightsSmartDetectionRuleRead(d, meta)
}
