package applicationinsights

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights/parse"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApplicationInsightsSmartDetectionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationInsightsSmartDetectionRuleUpdate,
		Read:   resourceApplicationInsightsSmartDetectionRuleRead,
		Update: resourceApplicationInsightsSmartDetectionRuleUpdate,
		Delete: resourceApplicationInsightsSmartDetectionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Slow page load time",
					"Slow server response time",
					"Long dependency duration",
					"Degredation in server response time",
					"Degredation in dependency duration",
					"Degradation in trace severity ratio",
					"Abnormal rise in exception volume",
					"Potential memory leak detected",
					"Potential security issue detected",
					"Abnormal rise in daily data volume",
				}, false),
			},

			"application_insights_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"additional_emails": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceApplicationInsightsSmartDetectionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights Samrt Detection Rule update.")

	name := strings.ToLower(strings.Join(strings.Split(d.Get("name").(string), " "), ""))
	appInsightsID := d.Get("application_insights_id").(string)

	id, err := parse.ComponentID(appInsightsID)
	if err != nil {
		return err
	}

	smartDetectionRuleProperties := insights.ApplicationInsightsComponentProactiveDetectionConfiguration{
		Name:         &name,
		Enabled:      utils.Bool(d.Get("enabled").(bool)),
		CustomEmails: utils.ExpandStringSlice(d.Get("additional_emails").(*schema.Set).List()),
	}

	_, err = client.Update(ctx, id.ResourceGroup, id.Name, name, smartDetectionRuleProperties)
	if err != nil {
		return fmt.Errorf("updating Application Insights Smart Detection Rule %q (Application Insights %q): %+v", name, id.String(), err)
	}

	d.SetId(fmt.Sprintf("%s/SmartDetectionRule/%s", id, name))

	return resourceApplicationInsightsSmartDetectionRuleRead(d, meta)
}

func resourceApplicationInsightsSmartDetectionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SmartDetectionRuleID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights Smart Detection Rule %q", id.String())

	result, err := client.Get(ctx, id.ResourceGroup, id.ComponentName, id.SmartDetectionRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			log.Printf("[WARN] AzureRM Application Insights Smart Detection Rule  %q not found, removing from state", id.String())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Application Insights Smart Detection Rule %q: %+v", id.String(), err)
	}

	d.Set("application_insights_id", parse.NewComponentID(id.SubscriptionId, id.ResourceGroup, id.ComponentName))

	d.Set("name", result.Name)
	d.Set("enabled", result.Enabled)
	d.Set("additional_emails", utils.FlattenStringSlice(result.CustomEmails))
	return nil
}

func resourceApplicationInsightsSmartDetectionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.SmartDetectionRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SmartDetectionRuleID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] reseting AzureRM Application Insights Smart Detection Rule %q", id.String())

	result, err := client.Get(ctx, id.ResourceGroup, id.ComponentName, id.SmartDetectionRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			log.Printf("[WARN] AzureRM Application Insights Smart Detection Rule %q not found, removing from state", id.String())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Application Insights Smart Detection Rule %q: %+v", id.String(), err)
	}

	smartDetectionRuleProperties := insights.ApplicationInsightsComponentProactiveDetectionConfiguration{
		Name:                           utils.String(id.SmartDetectionRuleName),
		Enabled:                        result.RuleDefinitions.IsEnabledByDefault,
		SendEmailsToSubscriptionOwners: utils.Bool(true),
		CustomEmails:                   utils.ExpandStringSlice([]interface{}{}),
	}

	_, err = client.Update(ctx, id.ResourceGroup, id.ComponentName, id.SmartDetectionRuleName, smartDetectionRuleProperties)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			return nil
		}
		return fmt.Errorf("issuing AzureRM reset update request for Application Insights Smart Detection Rule %q: %+v", id.String(), err)
	}

	return nil
}
