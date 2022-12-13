package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDatadogSingleSignOnConfigurations() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatadogSingleSignOnConfigurationsCreateorUpdate,
		Read:   resourceDatadogSingleSignOnConfigurationsRead,
		Update: resourceDatadogSingleSignOnConfigurationsCreateorUpdate,
		Delete: resourceDatadogSingleSignOnConfigurationsDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatadogSingleSignOnConfigurationsID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"datadog_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatadogMonitorID,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  utils.String("default"),
			},

			"enterprise_application_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DatadogEnterpriseApplicationID,
			},

			"single_sign_on_enabled": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"login_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDatadogSingleSignOnConfigurationsCreateorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.SingleSignOnConfigurationsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	datadogMonitorId := d.Get("datadog_monitor_id").(string)
	configurationName := d.Get("name").(string)
	enterpriseAppID := d.Get("enterprise_application_id").(string)
	id, err := parse.DatadogMonitorID(datadogMonitorId)
	if err != nil {
		return err
	}

	ssoId := parse.NewDatadogSingleSignOnConfigurationsID(id.SubscriptionId, id.ResourceGroup, id.MonitorName, configurationName).ID()

	existing, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, configurationName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing Datadog Monitor %q (Resource Group %q): %+v", id.ResourceGroup, id.MonitorName, err)
		}
	}

	singleSignOnState := datadog.SingleSignOnStatesEnable
	if d.Get("single_sign_on_enabled").(string) == "Disable" {
		singleSignOnState = datadog.SingleSignOnStatesDisable
	}

	body := datadog.SingleSignOnResource{
		Properties: &datadog.SingleSignOnProperties{
			SingleSignOnState: singleSignOnState,
			EnterpriseAppID:   utils.String(enterpriseAppID),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MonitorName, configurationName, &body); err != nil {
		return fmt.Errorf("configuring SingleSignOn on Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	d.SetId(ssoId)
	return resourceDatadogSingleSignOnConfigurationsRead(d, meta)
}

func resourceDatadogSingleSignOnConfigurationsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.SingleSignOnConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatadogSingleSignOnConfigurationsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.SingleSignOnConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Datadog monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}
	monitorId := parse.NewDatadogMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName)
	d.Set("datadog_monitor_id", monitorId.ID())
	d.Set("name", id.SingleSignOnConfigurationName)

	if props := resp.Properties; props != nil {
		d.Set("single_sign_on_enabled", props.SingleSignOnState)
		d.Set("login_url", props.SingleSignOnURL)
		d.Set("enterprise_application_id", props.EnterpriseAppID)
	}

	return nil
}

func resourceDatadogSingleSignOnConfigurationsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.SingleSignOnConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatadogSingleSignOnConfigurationsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.SingleSignOnConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Datadog monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	d.Set("enterprise_application_id", nil)
	enterpriseAppID := d.Get("enterprise_application_id").(string)

	body := datadog.SingleSignOnResource{
		Properties: &datadog.SingleSignOnProperties{
			SingleSignOnState: datadog.SingleSignOnStatesDisable,
			EnterpriseAppID:   utils.String(enterpriseAppID),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MonitorName, id.SingleSignOnConfigurationName, &body); err != nil {
		return fmt.Errorf("removing SingleSignOnConfiguration on Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	return nil
}
