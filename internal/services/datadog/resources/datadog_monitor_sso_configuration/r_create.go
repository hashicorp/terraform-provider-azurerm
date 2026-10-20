// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_sso_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogSingleSignOnConfigurationsCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.SingleSignOn
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := monitorsresource.ParseMonitorID(d.Get("datadog_monitor_id").(string))
	if err != nil {
		return err
	}

	id := singlesignon.NewSingleSignOnConfigurationID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.ConfigurationsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_datadog_monitor_sso_configuration", id.ID())
		}
	}

	payload := singlesignon.DatadogSingleSignOnResource{
		Properties: &singlesignon.DatadogSingleSignOnProperties{
			EnterpriseAppId: pointer.To(d.Get("enterprise_application_id").(string)),
		},
	}

	if v, ok := d.GetOk("single_sign_on"); ok {
		payload.Properties.SingleSignOnState = pointer.ToEnum[singlesignon.SingleSignOnStates](v.(string))
	}

	if err := client.ConfigurationsCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatadogSingleSignOnConfigurationsRead(d, meta)
}
