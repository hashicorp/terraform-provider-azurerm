package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2017-08-01-preview/securityinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSentinelAlertRuleMsSecurityIncident() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSentinelAlertRuleMsSecurityIncidentCreateUpdate,
		Read:   resourceArmSentinelAlertRuleMsSecurityIncidentRead,
		Update: resourceArmSentinelAlertRuleMsSecurityIncidentCreateUpdate,
		Delete: resourceArmSentinelAlertRuleMsSecurityIncidentDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.SentinelAlertRuleID(id)
			return err
		}, importSentinelAlertRule(securityinsight.MicrosoftSecurityIncidentCreation)),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"product_filter": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityinsight.MicrosoftCloudAppSecurity),
					string(securityinsight.AzureSecurityCenter),
					string(securityinsight.AzureActiveDirectoryIdentityProtection),
					string(securityinsight.AzureSecurityCenterforIoT),
					string(securityinsight.AzureAdvancedThreatProtection),
				}, false),
			},

			"severity_filter": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(securityinsight.High),
						string(securityinsight.Medium),
						string(securityinsight.Low),
						string(securityinsight.Informational),
					}, false),
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"text_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmSentinelAlertRuleMsSecurityIncidentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Ms Security Incident %q (Resource Group %q): %+v", name, workspaceID.ResourceGroup, err)
			}
		}

		id := alertRuleID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_ms_security_incident", *id)
		}
	}

	param := securityinsight.MicrosoftSecurityIncidentCreationAlertRule{
		Kind: securityinsight.KindMicrosoftSecurityIncidentCreation,
		MicrosoftSecurityIncidentCreationAlertRuleProperties: &securityinsight.MicrosoftSecurityIncidentCreationAlertRuleProperties{
			ProductFilter:    securityinsight.MicrosoftSecurityProductName(d.Get("product_filter").(string)),
			DisplayName:      utils.String(d.Get("display_name").(string)),
			Description:      utils.String(d.Get("description").(string)),
			Enabled:          utils.Bool(true),
			SeveritiesFilter: expandSeverityFilter(d.Get("severity_filter").(*schema.Set).List()),
		},
	}

	// If `text_whitelist` is set, it has to at least contain one item.
	if whitelist := d.Get("text_whitelist").(*schema.Set).List(); len(whitelist) != 0 {
		param.DisplayNamesFilter = utils.ExpandStringSlice(whitelist)
	}

	// Service avoid concurrent update of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		param.Etag = utils.String(d.Get("etag").(string))
	}

	if _, err := client.CreateOrUpdate(ctx, workspaceID.ResourceGroup, workspaceID.Name, name, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Ms Security Incident %q (Resource Group %q / Workspace: %q): %+v", name, workspaceID.ResourceGroup, workspaceID.Name, err)
	}

	resp, err := client.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Sentinel Alert Rule Ms Security Incident %q (Resource Group %q / Workspace: %q): %+v", name, workspaceID.ResourceGroup, workspaceID.Name, err)
	}
	id := alertRuleID(resp.Value)
	if id == nil || *id == "" {
		return fmt.Errorf("empty or nil ID returned for Sentinel Alert Rule Ms Security Incident %q (Resource Group %q / Workspace: %q) ID", name, workspaceID.ResourceGroup, workspaceID.Name)
	}
	d.SetId(*id)

	return resourceArmSentinelAlertRuleMsSecurityIncidentRead(d, meta)
}

func resourceArmSentinelAlertRuleMsSecurityIncidentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	workspaceClient := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SentinelAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Sentinel Alert Rule Ms Security Incident %q was not found in Workspace: %q in Resource Group %q - removing from state!", id.Name, id.Workspace, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Ms Security Incident %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	if err := assertAlertRuleKind(resp.Value, securityinsight.MicrosoftSecurityIncidentCreation); err != nil {
		return fmt.Errorf("asserting alert rule of %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}
	rule := resp.Value.(securityinsight.MicrosoftSecurityIncidentCreationAlertRule)

	d.Set("name", id.Name)

	workspaceResp, err := workspaceClient.Get(ctx, id.ResourceGroup, id.Workspace)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Workspace %q (Resource Group: %q) where this Alert Rule belongs to: %+v", id.Workspace, id.ResourceGroup, err)
	}
	workspaceID := ""
	if workspaceResp.ID != nil {
		workspaceID = *workspaceResp.ID
	}
	d.Set("log_analytics_workspace_id", workspaceID)

	if prop := rule.MicrosoftSecurityIncidentCreationAlertRuleProperties; prop != nil {
		d.Set("product_filter", string(prop.ProductFilter))

		displayName := ""
		if prop.DisplayName != nil {
			displayName = *prop.DisplayName
		}
		d.Set("display_name", displayName)

		description := ""
		if prop.Description != nil {
			description = *prop.Description
		}
		d.Set("description", description)

		if err := d.Set("text_whitelist", utils.FlattenStringSlice(prop.DisplayNamesFilter)); err != nil {
			return fmt.Errorf(`setting "text_whitelist": %+v`, err)
		}
		if err := d.Set("severity_filter", flattenSeverityFilter(prop.SeveritiesFilter)); err != nil {
			return fmt.Errorf(`setting "severity_filter": %+v`, err)
		}
	}

	etag := ""
	if rule.Etag != nil {
		etag = *rule.Etag
	}
	d.Set("etag", etag)

	return nil
}

func resourceArmSentinelAlertRuleMsSecurityIncidentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SentinelAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Ms Security Incident %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	return nil
}

func expandSeverityFilter(input []interface{}) *[]securityinsight.AlertSeverity {
	result := make([]securityinsight.AlertSeverity, 0)

	for _, e := range input {
		result = append(result, securityinsight.AlertSeverity(e.(string)))
	}

	return &result
}

func flattenSeverityFilter(input *[]securityinsight.AlertSeverity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, string(e))
	}

	return output
}
