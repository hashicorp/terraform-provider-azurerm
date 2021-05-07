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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSentinelAlertRuleScheduled() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentinelAlertRuleScheduledCreateUpdate,
		Read:   resourceSentinelAlertRuleScheduledRead,
		Update: resourceSentinelAlertRuleScheduledCreateUpdate,
		Delete: resourceSentinelAlertRuleScheduledDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.AlertRuleID(id)
			return err
		}, importSentinelAlertRule(securityinsight.AlertRuleKindScheduled)),

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

			"alert_rule_template_guid": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"event_grouping": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aggregation_method": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(securityinsight.EventGroupingAggregationKindAlertPerResult),
								string(securityinsight.EventGroupingAggregationKindSingleAlert),
							}, false),
						},
					},
				},
			},

			"tactics": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(securityinsight.AttackTacticCollection),
						string(securityinsight.AttackTacticCommandAndControl),
						string(securityinsight.AttackTacticCredentialAccess),
						string(securityinsight.AttackTacticDefenseEvasion),
						string(securityinsight.AttackTacticDiscovery),
						string(securityinsight.AttackTacticExecution),
						string(securityinsight.AttackTacticExfiltration),
						string(securityinsight.AttackTacticImpact),
						string(securityinsight.AttackTacticInitialAccess),
						string(securityinsight.AttackTacticLateralMovement),
						string(securityinsight.AttackTacticPersistence),
						string(securityinsight.AttackTacticPrivilegeEscalation),
					}, false),
				},
			},

			"incident_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_incident": {
							Required: true,
							Type:     schema.TypeBool,
						},
						"grouping": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"lookback_duration": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.ISO8601Duration,
										Default:      "PT5M",
									},
									"reopen_closed_incidents": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"entity_matching_method": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  securityinsight.EntitiesMatchingMethodNone,
										ValidateFunc: validation.StringInSlice([]string{
											string(securityinsight.EntitiesMatchingMethodAll),
											string(securityinsight.EntitiesMatchingMethodCustom),
											string(securityinsight.EntitiesMatchingMethodNone),
										}, false),
									},
									"group_by": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(securityinsight.GroupingEntityTypeAccount),
												string(securityinsight.GroupingEntityTypeHost),
												string(securityinsight.GroupingEntityTypeIP),
												string(securityinsight.GroupingEntityTypeURL),
											}, false),
										},
									},
								},
							},
						},
					},
				},
			},

			"severity": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityinsight.AlertSeverityHigh),
					string(securityinsight.AlertSeverityMedium),
					string(securityinsight.AlertSeverityLow),
					string(securityinsight.AlertSeverityInformational),
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
				Default:  string(securityinsight.TriggerOperatorGreaterThan),
				ValidateFunc: validation.StringInSlice([]string{
					string(securityinsight.TriggerOperatorGreaterThan),
					string(securityinsight.TriggerOperatorLessThan),
					string(securityinsight.TriggerOperatorEqual),
					string(securityinsight.TriggerOperatorNotEqual),
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

func resourceSentinelAlertRuleScheduledCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroup, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Scheduled %q: %+v", id, err)
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
		Kind: securityinsight.KindBasicAlertRuleKindScheduled,
		ScheduledAlertRuleProperties: &securityinsight.ScheduledAlertRuleProperties{
			Description:           utils.String(d.Get("description").(string)),
			DisplayName:           utils.String(d.Get("display_name").(string)),
			Tactics:               expandAlertRuleScheduledTactics(d.Get("tactics").(*schema.Set).List()),
			IncidentConfiguration: expandAlertRuleScheduledIncidentConfiguration(d.Get("incident_configuration").([]interface{})),
			Severity:              securityinsight.AlertSeverity(d.Get("severity").(string)),
			Enabled:               utils.Bool(d.Get("enabled").(bool)),
			Query:                 utils.String(d.Get("query").(string)),
			QueryFrequency:        &queryFreq,
			QueryPeriod:           &queryPeriod,
			SuppressionEnabled:    &suppressionEnabled,
			SuppressionDuration:   &suppressionDuration,
			TriggerOperator:       securityinsight.TriggerOperator(d.Get("trigger_operator").(string)),
			TriggerThreshold:      utils.Int32(int32(d.Get("trigger_threshold").(int))),
		},
	}

	if v, ok := d.GetOk("alert_rule_template_guid"); ok {
		param.ScheduledAlertRuleProperties.AlertRuleTemplateName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("event_grouping"); ok {
		param.ScheduledAlertRuleProperties.EventGroupingSettings = expandAlertRuleScheduledEventGroupingSetting(v.([]interface{}))
	}

	// Service avoid concurrent update of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindScheduled); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
		param.Etag = resp.Value.(securityinsight.ScheduledAlertRule).Etag
	}

	if _, err := client.CreateOrUpdate(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName, name, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleScheduledRead(d, meta)
}

func resourceSentinelAlertRuleScheduledRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Sentinel Alert Rule Scheduled %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindScheduled); err != nil {
		return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
	}
	rule := resp.Value.(securityinsight.ScheduledAlertRule)

	d.Set("name", id.Name)

	workspaceId := loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if prop := rule.ScheduledAlertRuleProperties; prop != nil {
		d.Set("description", prop.Description)
		d.Set("display_name", prop.DisplayName)
		if err := d.Set("tactics", flattenAlertRuleScheduledTactics(prop.Tactics)); err != nil {
			return fmt.Errorf("setting `tactics`: %+v", err)
		}
		if err := d.Set("incident_configuration", flattenAlertRuleScheduledIncidentConfiguration(prop.IncidentConfiguration)); err != nil {
			return fmt.Errorf("setting `incident_configuration`: %+v", err)
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
		d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)

		if err := d.Set("event_grouping", flattenAlertRuleScheduledEventGroupingSetting(prop.EventGroupingSettings)); err != nil {
			return fmt.Errorf("setting `event_grouping`: %+v", err)
		}
	}

	return nil
}

func resourceSentinelAlertRuleScheduledDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Scheduled %q: %+v", id, err)
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

func expandAlertRuleScheduledIncidentConfiguration(input []interface{}) *securityinsight.IncidentConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &securityinsight.IncidentConfiguration{
		CreateIncident:        utils.Bool(raw["create_incident"].(bool)),
		GroupingConfiguration: expandAlertRuleScheduledGrouping(raw["grouping"].([]interface{})),
	}

	return output
}

func flattenAlertRuleScheduledIncidentConfiguration(input *securityinsight.IncidentConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	createIncident := false
	if input.CreateIncident != nil {
		createIncident = *input.CreateIncident
	}

	return []interface{}{
		map[string]interface{}{
			"create_incident": createIncident,
			"grouping":        flattenAlertRuleScheduledGrouping(input.GroupingConfiguration),
		},
	}
}

func expandAlertRuleScheduledGrouping(input []interface{}) *securityinsight.GroupingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &securityinsight.GroupingConfiguration{
		Enabled:                utils.Bool(raw["enabled"].(bool)),
		ReopenClosedIncident:   utils.Bool(raw["reopen_closed_incidents"].(bool)),
		LookbackDuration:       utils.String(raw["lookback_duration"].(string)),
		EntitiesMatchingMethod: securityinsight.EntitiesMatchingMethod(raw["entity_matching_method"].(string)),
	}

	groupByEntitiesSet := raw["group_by"].(*schema.Set).List()
	groupByEntities := make([]securityinsight.GroupingEntityType, len(groupByEntitiesSet))
	for idx, t := range groupByEntitiesSet {
		groupByEntities[idx] = securityinsight.GroupingEntityType(t.(string))
	}
	output.GroupByEntities = &groupByEntities

	return output
}

func expandAlertRuleScheduledEventGroupingSetting(input []interface{}) *securityinsight.EventGroupingSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	result := securityinsight.EventGroupingSettings{}

	if aggregationKind := v["aggregation_method"].(string); aggregationKind != "" {
		result.AggregationKind = securityinsight.EventGroupingAggregationKind(aggregationKind)
	}

	return &result
}

func flattenAlertRuleScheduledGrouping(input *securityinsight.GroupingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	enabled := false
	if input.Enabled != nil {
		enabled = *input.Enabled
	}

	lookbackDuration := ""
	if input.LookbackDuration != nil {
		lookbackDuration = *input.LookbackDuration
	}

	reopenClosedIncidents := false
	if input.ReopenClosedIncident != nil {
		reopenClosedIncidents = *input.ReopenClosedIncident
	}

	var groupByEntities []interface{}
	if input.GroupByEntities != nil {
		for _, entity := range *input.GroupByEntities {
			groupByEntities = append(groupByEntities, string(entity))
		}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":                 enabled,
			"lookback_duration":       lookbackDuration,
			"reopen_closed_incidents": reopenClosedIncidents,
			"entity_matching_method":  string(input.EntitiesMatchingMethod),
			"group_by":                groupByEntities,
		},
	}
}

func flattenAlertRuleScheduledEventGroupingSetting(input *securityinsight.EventGroupingSettings) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var aggregationKind string
	if input.AggregationKind != "" {
		aggregationKind = string(input.AggregationKind)
	}

	return []interface{}{
		map[string]interface{}{
			"aggregation_method": aggregationKind,
		},
	}
}
