package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatadogMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DatadogMonitorsName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"configuration_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  utils.String("default"),
			},

			"enterpriseapp_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"singlesignon_state": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Enable",
					"Disable",
				}, false),
			},

			"singlesignon_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDatadogSingleSignOnConfigurationsCreateorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Datadog.SingleSignOnConfigurationsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	configurationName := d.Get("configuration_name").(string)
	enterpriseAppID := d.Get("enterpriseapp_id").(string)

	id := parse.NewDatadogSingleSignOnConfigurationsID(subscriptionId, resourceGroup, name, configurationName).ID()

	existing, err := client.Get(ctx, resourceGroup, name, configurationName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing Datadog Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	body := datadog.SingleSignOnResource{
		Properties: &datadog.SingleSignOnProperties{
			EnterpriseAppID: utils.String(enterpriseAppID),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, configurationName, &body); err != nil {
		return fmt.Errorf("Configuring SingleSignOn on Datadog Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id)
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
			log.Printf("[INFO] datadog %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	if props := resp.Properties; props != nil {
		d.Set("singlesignon_state", props.SingleSignOnState)
		d.Set("provisioning_state", props.ProvisioningState)
		d.Set("singlesignon_url", props.SingleSignOnURL)
	}

	d.Set("type", resp.Type)
	d.Set("id", resp.ID)

	return nil
}

func resourceDatadogSingleSignOnConfigurationsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
