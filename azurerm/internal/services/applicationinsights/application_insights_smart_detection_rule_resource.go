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
				}, false),
				DiffSuppressFunc: smartDetectionRuleNameDiff,
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

			"send_emails_to_subscription_owners": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"additional_email_recipients": {
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

	// The Smart Detection Rule name from the UI doesn't match what the API accepts.
	// We'll have the user submit what the name looks like in the UI and trim it behind the scenes to match what the API accepts
	name := strings.ToLower(strings.Join(strings.Split(d.Get("name").(string), " "), ""))
	appInsightsID := d.Get("application_insights_id").(string)

	id, err := parse.ComponentID(appInsightsID)
	if err != nil {
		return err
	}

	smartDetectionRuleProperties := insights.ApplicationInsightsComponentProactiveDetectionConfiguration{
		Name:                           &name,
		Enabled:                        utils.Bool(d.Get("enabled").(bool)),
		SendEmailsToSubscriptionOwners: utils.Bool(d.Get("send_emails_to_subscription_owners").(bool)),
		CustomEmails:                   utils.ExpandStringSlice(d.Get("additional_email_recipients").(*schema.Set).List()),
	}

	_, err = client.Update(ctx, id.ResourceGroup, id.Name, name, smartDetectionRuleProperties)
	if err != nil {
		return fmt.Errorf("updating Application Insights Smart Detection Rule %q (Application Insights %q): %+v", name, id.String(), err)
	}

	d.SetId(fmt.Sprintf("%s/SmartDetectionRule/%s", id.ID(), name))

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

	d.Set("name", result.Name)
	d.Set("application_insights_id", parse.NewComponentID(id.SubscriptionId, id.ResourceGroup, id.ComponentName).ID())
	d.Set("enabled", result.Enabled)
	d.Set("send_emails_to_subscription_owners", result.SendEmailsToSubscriptionOwners)
	d.Set("additional_email_recipients", utils.FlattenStringSlice(result.CustomEmails))
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
		SendEmailsToSubscriptionOwners: result.RuleDefinitions.SupportsEmailNotifications,
		CustomEmails:                   utils.ExpandStringSlice([]interface{}{}),
	}

	// Application Insights defaults all the Smart Detection Rules so if a user wants to delete a rule, we'll update it back to it's default values.
	_, err = client.Update(ctx, id.ResourceGroup, id.ComponentName, id.SmartDetectionRuleName, smartDetectionRuleProperties)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			return nil
		}
		return fmt.Errorf("issuing AzureRM reset update request for Application Insights Smart Detection Rule %q: %+v", id.String(), err)
	}

	return nil
}

// The Smart Detection Rule name from the UI doesn't match what the API accepts.
// This Diff checks that the name UI name matches the API name when spaces are removed
func smartDetectionRuleNameDiff(_, old string, new string, _ *schema.ResourceData) bool {
	trimmedNew := strings.Join(strings.Split(strings.ToLower(new), " "), "")

	return strings.EqualFold(old, trimmedNew)
}
