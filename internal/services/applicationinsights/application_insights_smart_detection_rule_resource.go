// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	smartdetection "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApplicationInsightsSmartDetectionRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsSmartDetectionRuleUpdate,
		Read:   resourceApplicationInsightsSmartDetectionRuleRead,
		Update: resourceApplicationInsightsSmartDetectionRuleUpdate,
		Delete: resourceApplicationInsightsSmartDetectionRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := smartdetection.ParseProactiveDetectionConfigID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SmartDetectionRuleUpgradeV0ToV1{},
			1: migration.SmartDetectionRuleUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Slow page load time",
					"Slow server response time",
					"Long dependency duration",
					"Degradation in server response time",
					"Degradation in dependency duration",
					// The below rules are currently preview and may change in future
					"Degradation in trace severity ratio",
					"Abnormal rise in exception volume",
					"Potential memory leak detected",
					"Potential security issue detected",
					"Abnormal rise in daily data volume",
				}, false),
				DiffSuppressFunc: smartDetectionRuleNameDiff,
			},

			"application_insights_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: components.ValidateComponentID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"send_emails_to_subscription_owners": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"additional_email_recipients": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},
		},
	}
}

func resourceApplicationInsightsSmartDetectionRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights Smart Detection Rule update.")

	// The Smart Detection Rule name from the UI doesn't match what the API accepts.
	// We'll have the user submit what the name looks like in the UI and convert it behind the scenes to match what the API accepts
	name := convertUiNameToApiName(d.Get("name"))

	appInsightsId, err := smartdetection.ParseComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	id := smartdetection.NewProactiveDetectionConfigID(appInsightsId.SubscriptionId, appInsightsId.ResourceGroupName, appInsightsId.ComponentName, name)

	smartDetectionRuleProperties := smartdetection.ApplicationInsightsComponentProactiveDetectionConfiguration{
		Name:                           &name,
		Enabled:                        pointer.To(d.Get("enabled").(bool)),
		SendEmailsToSubscriptionOwners: pointer.To(d.Get("send_emails_to_subscription_owners").(bool)),
		CustomEmails:                   utils.ExpandStringSlice(d.Get("additional_email_recipients").(*pluginsdk.Set).List()),
	}

	_, err = client.ProactiveDetectionConfigurationsUpdate(ctx, id, smartDetectionRuleProperties)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApplicationInsightsSmartDetectionRuleRead(d, meta)
}

func resourceApplicationInsightsSmartDetectionRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := smartdetection.ParseProactiveDetectionConfigID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights Smart Detection Rule %s", id)

	resp, err := client.ProactiveDetectionConfigurationsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found, removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("application_insights_id", smartdetection.NewComponentID(id.SubscriptionId, id.ResourceGroupName, id.ComponentName).ID())

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
		d.Set("enabled", model.Enabled)
		d.Set("send_emails_to_subscription_owners", model.SendEmailsToSubscriptionOwners)
		d.Set("additional_email_recipients", utils.FlattenStringSlice(model.CustomEmails))

	}
	return nil
}

func resourceApplicationInsightsSmartDetectionRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := smartdetection.ParseProactiveDetectionConfigID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] reseting AzureRM Application Insights Smart Detection Rule %s", id)

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
	smartDetectionRuleProperties := smartdetection.ApplicationInsightsComponentProactiveDetectionConfiguration{
		Name:                           pointer.To(id.ConfigurationId),
		Enabled:                        resp.Model.RuleDefinitions.IsEnabledByDefault,
		SendEmailsToSubscriptionOwners: resp.Model.RuleDefinitions.SupportsEmailNotifications,
		CustomEmails:                   utils.ExpandStringSlice([]interface{}{}),
	}

	// Application Insights defaults all the Smart Detection Rules so if a user wants to delete a rule, we'll update it back to it's default values.
	_, err = client.ProactiveDetectionConfigurationsUpdate(ctx, *id, smartDetectionRuleProperties)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("resetting %s: %+v", id, err)
	}

	return nil
}

// The Smart Detection Rule name from the UI doesn't match what the API accepts.
// This Diff checks if the old and new name match when converted to the API version of the name
func smartDetectionRuleNameDiff(_, old string, new string, _ *pluginsdk.ResourceData) bool {
	apiNew := convertUiNameToApiName(new)

	return strings.EqualFold(old, apiNew)
}

func convertUiNameToApiName(uiName interface{}) string {
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
