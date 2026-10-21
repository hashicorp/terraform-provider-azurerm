// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_sso_configuration

import (
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogSingleSignOnConfigurationsRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.SingleSignOn
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := singlesignon.ParseSingleSignOnConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ConfigurationsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
	}

	d.Set("name", id.SingleSignOnConfigurationName)
	d.Set("datadog_monitor_id", monitorsresource.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			// per the create func
			d.Set("single_sign_on", pointer.FromEnum(props.SingleSignOnState))
			d.Set("login_url", props.SingleSignOnURL)
			d.Set("enterprise_application_id", props.EnterpriseAppId)
		}
	}

	return nil
}
