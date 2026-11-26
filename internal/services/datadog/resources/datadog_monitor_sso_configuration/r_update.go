// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_sso_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogsinglesignonresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogSingleSignOnConfigurationsUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.DatadogSingleSignOnResources
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datadogsinglesignonresources.ParseSingleSignOnConfigurationID(d.Id())
	if err != nil {
		return err
	}

	payload := datadogsinglesignonresources.DatadogSingleSignOnResource{
		Properties: &datadogsinglesignonresources.DatadogSingleSignOnProperties{
			EnterpriseAppId: pointer.To(d.Get("enterprise_application_id").(string)),
		},
	}

	if v, ok := d.GetOk("single_sign_on"); ok && d.HasChange("single_sign_on") {
		payload.Properties.SingleSignOnState = pointer.ToEnum[datadogsinglesignonresources.SingleSignOnStates](v.(string))
	}

	if err := client.SingleSignOnConfigurationsCreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDatadogSingleSignOnConfigurationsRead(d, meta)
}
