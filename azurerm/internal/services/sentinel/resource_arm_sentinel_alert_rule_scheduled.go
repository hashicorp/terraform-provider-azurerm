package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/rickb777/date/period"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSentinelAlertRuleScheduled() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSentinelAlertRuleScheduledCreateUpdate,
		Read:   resourceArmSentinelAlertRuleScheduledRead,
		Update: resourceArmSentinelAlertRuleScheduledCreateUpdate,
		Delete: resourceArmSentinelAlertRuleScheduledDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.SentinelAlertRuleID(id)
			return err
		}, importSentinelAlertRule(securityinsight.Scheduled)),

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

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tactics": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(securityinsight.Collection),
						string(securityinsight.CommandAndControl),
						string(securityinsight.CredentialAccess),
						string(securityinsight.DefenseEvasion),
						string(securityinsight.Discovery),
						string(securityinsight.Execution),
						string(securityinsight.Exfiltration),
						string(securityinsight.Impact),
						string(securityinsight.InitialAccess),
						string(securityinsight.LateralMovement),
						string(securityinsight.Persistence),
						string(securityinsight.PrivilegeEscalation),
					}, false),
				},
			},

			"severity": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityinsight.High),
					string(securityinsight.Medium),
					string(securityinsight.Low),
					string(securityinsight.Informational),
				}, false),
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"query": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"query_frequency": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "PT24H"),
			},

			"query_period": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "P14D"),
			},

			"trigger_operator": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(securityinsight.GreaterThan),
				ValidateFunc: validation.StringInSlice([]string{
					string(securityinsight.GreaterThan),
					string(securityinsight.LessThan),
					string(securityinsight.Equal),
					string(securityinsight.NotEqual),
				}, false),
			},

			"trigger_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"suppression_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"suppression_duration": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "PT24H"),
			},
		},
	}
}

func resourceArmSentinelAlertRuleScheduledCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Scheduled %q (Resource Group %q / Workspace %q): %+v", name, workspaceID.ResourceGroup, workspaceID.Name, err)
			}
		}

		id := alertRuleID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_scheduled", *id)
		}
	}

	// Sanity checks

	// query frequency must <= query period: ensure there is no gaps in the overall query coverage.
	queryFreq := d.Get("query_frequency").(string)
	queryFreqDuration := period.MustParse(queryFreq).DurationApprox()

	queryPeriod := d.Get("query_period").(string)
	queryPeriodDuration := period.MustParse(queryPeriod).DurationApprox()
	if queryFreqDuration > queryPeriodDuration {
		return fmt.Errorf("`query_frequency`(%v) should not be larger than `query period`(%v), which introduce gaps in the overall query coverage", queryFreq, queryPeriod)
	}

	// query frequency must <= suppression duration: otherwise suppression has no effect.
	suppressionDuration := d.Get("suppression_duration").(string)
	suppressionEnabled := d.Get("suppression_enabled").(bool)
	if suppressionEnabled {
		suppressionDurationDuration := period.MustParse(suppressionDuration).DurationApprox()
		if queryFreqDuration > suppressionDurationDuration {
			return fmt.Errorf("`query_frequency`(%v) should not be larger than `suppression_duration`(%v), which makes suppression pointless", queryFreq, suppressionDuration)
		}
	}

	param := securityinsight.ScheduledAlertRule{
		Kind: securityinsight.KindScheduled,
		ScheduledAlertRuleProperties: &securityinsight.ScheduledAlertRuleProperties{
			Description:         utils.String(d.Get("description").(string)),
			DisplayName:         utils.String(d.Get("display_name").(string)),
			Tactics:             expandAlertRuleScheduledTactics(d.Get("tactics").(*schema.Set).List()),
			Severity:            securityinsight.AlertSeverity(d.Get("severity").(string)),
			Enabled:             utils.Bool(d.Get("enabled").(bool)),
			Query:               utils.String(d.Get("query").(string)),
			QueryFrequency:      &queryFreq,
			QueryPeriod:         &queryPeriod,
			SuppressionEnabled:  &suppressionEnabled,
			SuppressionDuration: &suppressionDuration,
			TriggerOperator:     securityinsight.TriggerOperator(d.Get("trigger_operator").(string)),
			TriggerThreshold:    utils.Int32(int32(d.Get("trigger_threshold").(int))),
		},
	}

	// Service avoid concurrent update of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.Name, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q (Resource Group %q / Workspace %q): %+v", name, workspaceID.ResourceGroup, workspaceID.Name, err)
		}

		if err := assertAlertRuleKind(resp.Value, securityinsight.Scheduled); err != nil {
			return fmt.Errorf("asserting alert rule of %q (Resource Group %q / Workspace %q): %+v", name, workspaceID.ResourceGroup, workspaceID.Name, err)
		}
		param.Etag = resp.Value.(securityinsight.ScheduledAlertRule).Etag
	}

	if _, err := client.CreateOrUpdate(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.Name, name, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Scheduled %q (Resource Group %q / Workspace %q): %+v", name, workspaceID.ResourceGroup, workspaceID.Name, err)
	}

	resp, err := client.Get(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q (Resource Group %q / Workspace %q): %+v", name, workspaceID.ResourceGroup, workspaceID.Name, err)
	}

	id := alertRuleID(resp.Value)
	if id == nil || *id == "" {
		return fmt.Errorf("empty or nil ID returned for Sentinel Alert Rule Scheduled %q (Resource Group %q / Workspace %q) ID", name, workspaceID.ResourceGroup, workspaceID.Name)
	}
	d.SetId(*id)

	return resourceArmSentinelAlertRuleScheduledRead(d, meta)
}

func resourceArmSentinelAlertRuleScheduledRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	workspaceClient := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SentinelAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.Workspace, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Sentinel Alert Rule Scheduled %q was not found in Workspace %q in Resource Group %q - removing from state!", id.Name, id.Workspace, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q (Resource Group %q / Workspace %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	if err := assertAlertRuleKind(resp.Value, securityinsight.Scheduled); err != nil {
		return fmt.Errorf("asserting alert rule of %q (Resource Group %q / Workspace %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}
	rule := resp.Value.(securityinsight.ScheduledAlertRule)

	d.Set("name", id.Name)

	workspaceResp, err := workspaceClient.Get(ctx, id.ResourceGroup, id.Workspace)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Workspace %q (Resource Group: %q) where this Alert Rule belongs to: %+v", id.Workspace, id.ResourceGroup, err)
	}
	d.Set("log_analytics_workspace_id", workspaceResp.ID)

	if prop := rule.ScheduledAlertRuleProperties; prop != nil {
		d.Set("description", prop.Description)
		d.Set("display_name", prop.DisplayName)
		if err := d.Set("tactics", flattenAlertRuleScheduledTactics(prop.Tactics)); err != nil {
			return fmt.Errorf("setting `tactics`: %+v", err)
		}
		d.Set("severity", string(prop.Severity))
		d.Set("enabled", prop.Enabled)
		d.Set("query", prop.Query)
		d.Set("query_frequency", prop.QueryFrequency)
		d.Set("query_period", prop.QueryPeriod)
		d.Set("trigger_operator", string(prop.TriggerOperator))

		var threshold int32
		if prop.TriggerThreshold != nil {
			threshold = *prop.TriggerThreshold
		}

		d.Set("trigger_threshold", int(threshold))
		d.Set("suppression_enabled", prop.SuppressionEnabled)
		d.Set("suppression_duration", prop.SuppressionDuration)
	}

	return nil
}

func resourceArmSentinelAlertRuleScheduledDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SentinelAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.Workspace, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Scheduled %q (Resource Group %q / Workspace %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	return nil
}

func expandAlertRuleScheduledTactics(input []interface{}) *[]securityinsight.AttackTactic {
	result := make([]securityinsight.AttackTactic, 0)

	for _, e := range input {
		result = append(result, securityinsight.AttackTactic(e.(string)))
	}

	return &result
}

func flattenAlertRuleScheduledTactics(input *[]securityinsight.AttackTactic) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, string(e))
	}

	return output
}
