// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_sso_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogSingleSignOnConfigurationsDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.SingleSignOn
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := singlesignon.ParseSingleSignOnConfigurationID(d.Id())
	if err != nil {
		return err
	}

	// SingleSignOnConfigurations can't be removed, but can be disabled/reset, which is what we do here
	payload := singlesignon.DatadogSingleSignOnResource{
		Properties: &singlesignon.DatadogSingleSignOnProperties{
			SingleSignOnState: pointer.To(singlesignon.SingleSignOnStatesDisable),
			EnterpriseAppId:   pointer.To(""),
		},
	}

	if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("removing %s: %+v", id, err)
	}

	return nil
}
