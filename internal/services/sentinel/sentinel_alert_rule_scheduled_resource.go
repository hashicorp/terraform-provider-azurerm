package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/rickb777/date/period"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func resourceSentinelAlertRuleScheduled() *pluginsdk.Resource {
	var entityMappingTypes = []string{
		string(securityinsight.EntityMappingTypeAccount),
		string(securityinsight.EntityMappingTypeAzureResource),
		string(securityinsight.EntityMappingTypeCloudApplication),
		string(securityinsight.EntityMappingTypeDNS),
		string(securityinsight.EntityMappingTypeFile),
		string(securityinsight.EntityMappingTypeFileHash),
		string(securityinsight.EntityMappingTypeHost),
		string(securityinsight.EntityMappingTypeIP),
		string(securityinsight.EntityMappingTypeMailbox),
		string(securityinsight.EntityMappingTypeMailCluster),
		string(securityinsight.EntityMappingTypeMailMessage),
		string(securityinsight.EntityMappingTypeMalware),
		string(securityinsight.EntityMappingTypeProcess),
		string(securityinsight.EntityMappingTypeRegistryKey),
		string(securityinsight.EntityMappingTypeRegistryValue),
		string(securityinsight.EntityMappingTypeSecurityGroup),
		string(securityinsight.EntityMappingTypeSubmissionMail),
		string(securityinsight.EntityMappingTypeURL),
	}
	return &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleScheduledCreateUpdate,
		Read:   resourceSentinelAlertRuleScheduledRead,
		Update: resourceSentinelAlertRuleScheduledCreateUpdate,
		Delete: resourceSentinelAlertRuleScheduledDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.AlertRuleID(id)
			return err
		}, importSentinelAlertRule(securityinsight.AlertRuleKindScheduled)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"alert_rule_template_guid": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"alert_rule_template_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"event_grouping": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"aggregation_method": {
							Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(securityinsight.AttackTacticCollection),
						string(securityinsight.AttackTacticCommandAndControl),
						string(securityinsight.AttackTacticCredentialAccess),
						string(securityinsight.AttackTacticDefenseEvasion),
						string(securityinsight.AttackTacticDiscovery),
						string(securityinsight.AttackTacticExecution),
						string(securityinsight.AttackTacticExfiltration),
						string(securityinsight.AttackTacticImpairProcessControl),
						string(securityinsight.AttackTacticInhibitResponseFunction),
						string(securityinsight.AttackTacticImpact),
						string(securityinsight.AttackTacticInitialAccess),
						string(securityinsight.AttackTacticLateralMovement),
						string(securityinsight.AttackTacticPersistence),
						string(securityinsight.AttackTacticPrivilegeEscalation),
						string(securityinsight.AttackTacticPreAttack),
						string(securityinsight.AttackTacticReconnaissance),
						string(securityinsight.AttackTacticResourceDevelopment),
					}, false),
				},
			},

			"techniques": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			// TODO 4.0 - rename this to "incident"
			"incident_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// TODO 4.0 - rename this to "create_incident_enabled"
						"create_incident": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
						"grouping": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
									},
									"lookback_duration": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validate.ISO8601Duration,
										Default:      "PT5M",
									},
									"reopen_closed_incidents": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
									"entity_matching_method": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  securityinsight.MatchingMethodAnyAlert,
										ValidateFunc: validation.StringInSlice([]string{
											string(securityinsight.MatchingMethodAnyAlert),
											string(securityinsight.MatchingMethodSelected),
											string(securityinsight.MatchingMethodAllEntities),
										}, false),
									},
									// TODO 4.0 - rename this to "by_entities"
									"group_by_entities": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(entityMappingTypes, false),
										},
									},
									// TODO 4.0 - rename this to "by_alert_details"
									"group_by_alert_details": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(securityinsight.AlertDetailDisplayName),
												string(securityinsight.AlertDetailSeverity),
											},
												false),
										},
									},
									// TODO 4.0 - rename this to "by_custom_details"
									"group_by_custom_details": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},
			},

			"severity": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityinsight.AlertSeverityHigh),
					string(securityinsight.AlertSeverityMedium),
					string(securityinsight.AlertSeverityLow),
					string(securityinsight.AlertSeverityInformational),
				}, false),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"query": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"query_frequency": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "P14D"),
			},

			"query_period": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "P14D"),
			},

			"trigger_operator": {
				Type:     pluginsdk.TypeString,
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
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"suppression_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
			"suppression_duration": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "PT24H"),
			},
			"alert_details_override": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description_format": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"display_name_format": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"severity_column_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"tactics_column_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"dynamic_property": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{
												string(securityinsight.AlertPropertyAlertLink),
												string(securityinsight.AlertPropertyConfidenceLevel),
												string(securityinsight.AlertPropertyConfidenceScore),
												string(securityinsight.AlertPropertyExtendedLinks),
												string(securityinsight.AlertPropertyProductComponentName),
												string(securityinsight.AlertPropertyProductName),
												string(securityinsight.AlertPropertyProviderName),
												string(securityinsight.AlertPropertyRemediationSteps),
												string(securityinsight.AlertPropertyTechniques),
											}, false),
									},
									"value": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
			"custom_details": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"entity_mapping": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 5,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"entity_type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(entityMappingTypes, false),
						},
						"field_mapping": {
							Type:     pluginsdk.TypeList,
							MaxItems: 3,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"identifier": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"column_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
			"sentinel_entity_mapping": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 5,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"column_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceSentinelAlertRuleScheduledCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
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
			Tactics:               expandAlertRuleTactics(d.Get("tactics").(*pluginsdk.Set).List()),
			Techniques:            expandAlertRuleTechnicals(d.Get("techniques").(*pluginsdk.Set).List()),
			IncidentConfiguration: expandAlertRuleIncidentConfiguration(d.Get("incident_configuration").([]interface{}), "create_incident", true),
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
	if v, ok := d.GetOk("alert_rule_template_version"); ok {
		param.ScheduledAlertRuleProperties.TemplateVersion = utils.String(v.(string))
	}
	if v, ok := d.GetOk("event_grouping"); ok {
		param.ScheduledAlertRuleProperties.EventGroupingSettings = expandAlertRuleScheduledEventGroupingSetting(v.([]interface{}))
	}
	if v, ok := d.GetOk("alert_details_override"); ok {
		param.ScheduledAlertRuleProperties.AlertDetailsOverride = expandAlertRuleAlertDetailsOverride(v.([]interface{}))
	}
	if v, ok := d.GetOk("custom_details"); ok {
		param.ScheduledAlertRuleProperties.CustomDetails = utils.ExpandMapStringPtrString(v.(map[string]interface{}))
	}

	entityMappingCount := 0
	sentinelEntityMappingCount := 0
	if v, ok := d.GetOk("entity_mapping"); ok {
		param.ScheduledAlertRuleProperties.EntityMappings = expandAlertRuleEntityMapping(v.([]interface{}))
		entityMappingCount = len(*param.ScheduledAlertRuleProperties.EntityMappings)
	}
	if v, ok := d.GetOk("sentinel_entity_mapping"); ok {
		param.ScheduledAlertRuleProperties.SentinelEntitiesMappings = expandAlertRuleSentinelEntityMapping(v.([]interface{}))
		sentinelEntityMappingCount = len(*param.ScheduledAlertRuleProperties.SentinelEntitiesMappings)
	}

	// the max number of `sentinel_entity_mapping` and `entity_mapping` together is 5
	if entityMappingCount+sentinelEntityMappingCount > 5 {
		return fmt.Errorf("`entity_mapping` and `sentinel_entity_mapping` together can't exceed 5")
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindScheduled); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.CreateOrUpdate(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleScheduledRead(d, meta)
}

func resourceSentinelAlertRuleScheduledRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
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

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if prop := rule.ScheduledAlertRuleProperties; prop != nil {
		d.Set("description", prop.Description)
		d.Set("display_name", prop.DisplayName)
		if err := d.Set("tactics", flattenAlertRuleTactics(prop.Tactics)); err != nil {
			return fmt.Errorf("setting `tactics`: %+v", err)
		}
		if err := d.Set("techniques", prop.Techniques); err != nil {
			return fmt.Errorf("setting `techniques`: %+v", err)
		}
		if err := d.Set("incident_configuration", flattenAlertRuleIncidentConfiguration(prop.IncidentConfiguration, "create_incident", true)); err != nil {
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
		d.Set("alert_rule_template_version", prop.TemplateVersion)

		if err := d.Set("event_grouping", flattenAlertRuleScheduledEventGroupingSetting(prop.EventGroupingSettings)); err != nil {
			return fmt.Errorf("setting `event_grouping`: %+v", err)
		}
		if err := d.Set("alert_details_override", flattenAlertRuleAlertDetailsOverride(prop.AlertDetailsOverride)); err != nil {
			return fmt.Errorf("setting `alert_details_override`: %+v", err)
		}
		if err := d.Set("custom_details", utils.FlattenMapStringPtrString(prop.CustomDetails)); err != nil {
			return fmt.Errorf("setting `custom_details`: %+v", err)
		}
		if err := d.Set("entity_mapping", flattenAlertRuleEntityMapping(prop.EntityMappings)); err != nil {
			return fmt.Errorf("setting `entity_mapping`: %+v", err)
		}
		if err := d.Set("sentinel_entity_mapping", flattenAlertRuleSentinelEntityMapping(prop.SentinelEntitiesMappings)); err != nil {
			return fmt.Errorf("setting `sentinel_entity_mapping`: %+v", err)
		}
	}

	return nil
}

func resourceSentinelAlertRuleScheduledDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	return nil
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
