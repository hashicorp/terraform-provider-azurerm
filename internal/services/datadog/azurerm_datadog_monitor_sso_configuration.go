// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// @tombuildsstuff: in 4.0 consider inlining this within the `azurerm_datadog_monitors` resource
// since this appears to be a 1:1 with it (given the name defaults to `default`)

func resourceDatadogSingleSignOnConfigurations() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatadogSingleSignOnConfigurationsCreate,
		Read:   resourceDatadogSingleSignOnConfigurationsRead,
		Update: resourceDatadogSingleSignOnConfigurationsUpdate,
		Delete: resourceDatadogSingleSignOnConfigurationsDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := singlesignon.ParseSingleSignOnConfigurationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"datadog_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: monitorsresource.ValidateMonitorID,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "default",
			},

			"enterprise_application_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DatadogEnterpriseApplicationID,
			},

			"single_sign_on_enabled": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					// @tombuildsstuff: other options are available, but the Create handles this as a boolean for now
					// should the field be a boolean? one to consider for 4.0 when this resource is inlined
					string(singlesignon.SingleSignOnStatesEnable),
					string(singlesignon.SingleSignOnStatesDisable),
				}, false),
			},

			"login_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDatadogSingleSignOnConfigurationsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.SingleSignOn
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := monitorsresource.ParseMonitorID(d.Get("datadog_monitor_id").(string))
	if err != nil {
		return err
	}

	id := singlesignon.NewSingleSignOnConfigurationID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, d.Get("name").(string))
	existing, err := client.ConfigurationsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_datadog_monitor_sso_configuration", id.ID())
	}

	payload := singlesignon.DatadogSingleSignOnResource{
		Properties: &singlesignon.DatadogSingleSignOnProperties{
			SingleSignOnState: pointer.To(singlesignon.SingleSignOnStates(d.Get("single_sign_on_enabled").(string))),
			EnterpriseAppId:   utils.String(d.Get("enterprise_application_id").(string)),
		},
	}

	if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatadogSingleSignOnConfigurationsRead(d, meta)
}

func resourceDatadogSingleSignOnConfigurationsRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			singleSignOnEnabled := props.SingleSignOnState != nil && *props.SingleSignOnState == singlesignon.SingleSignOnStatesEnable
			d.Set("single_sign_on_enabled", singleSignOnEnabled)
			d.Set("login_url", props.SingleSignOnURL)
			d.Set("enterprise_application_id", props.EnterpriseAppId)
		}
	}

	return nil
}

func resourceDatadogSingleSignOnConfigurationsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.SingleSignOn
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := singlesignon.ParseSingleSignOnConfigurationID(d.Id())
	if err != nil {
		return err
	}

	payload := singlesignon.DatadogSingleSignOnResource{
		Properties: &singlesignon.DatadogSingleSignOnProperties{
			SingleSignOnState: pointer.To(singlesignon.SingleSignOnStates(d.Get("single_sign_on_enabled").(string))),
			EnterpriseAppId:   utils.String(d.Get("enterprise_application_id").(string)),
		},
	}

	if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDatadogSingleSignOnConfigurationsRead(d, meta)
}

func resourceDatadogSingleSignOnConfigurationsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
			EnterpriseAppId:   utils.String(""),
		},
	}

	if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("removing %s: %+v", id, err)
	}

	return nil
}
