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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func resourceSentinelAlertRuleNrt() *pluginsdk.Resource {
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
		Create: resourceSentinelAlertRuleNrtCreateUpdate,
		Read:   resourceSentinelAlertRuleNrtRead,
		Update: resourceSentinelAlertRuleNrtCreateUpdate,
		Delete: resourceSentinelAlertRuleNrtDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.AlertRuleID(id)
			return err
		}, importSentinelAlertRule(securityinsight.AlertRuleKindNRT)),

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
				Required: features.FourPointOhBeta(),
				Optional: !features.FourPointOhBeta(),
				Computed: !features.FourPointOhBeta(), // the service will default it to `SingleAlert`.
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
						string(securityinsight.AttackTacticImpact),
						string(securityinsight.AttackTacticInitialAccess),
						string(securityinsight.AttackTacticLateralMovement),
						string(securityinsight.AttackTacticPersistence),
						string(securityinsight.AttackTacticPrivilegeEscalation),
						string(securityinsight.AttackTacticPreAttack),
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

			"incident": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"create_incident_enabled": {
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
									"by_entities": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(entityMappingTypes, false),
										},
									},
									"by_alert_details": {
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
									"by_custom_details": {
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

func resourceSentinelAlertRuleNrtCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}

		id := alertRuleID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_nrt", *id)
		}
	}

	// query frequency must <= suppression duration: otherwise suppression has no effect.
	suppressionDuration := d.Get("suppression_duration").(string)
	suppressionEnabled := d.Get("suppression_enabled").(bool)

	param := securityinsight.NrtAlertRule{
		Kind: securityinsight.KindBasicAlertRuleKindNRT,
		NrtAlertRuleProperties: &securityinsight.NrtAlertRuleProperties{
			Description:           utils.String(d.Get("description").(string)),
			DisplayName:           utils.String(d.Get("display_name").(string)),
			Techniques:            expandAlertRuleTechnicals(d.Get("techniques").(*pluginsdk.Set).List()),
			Tactics:               expandAlertRuleTactics(d.Get("tactics").(*pluginsdk.Set).List()),
			IncidentConfiguration: expandAlertRuleIncidentConfiguration(d.Get("incident").([]interface{}), "create_incident_enabled", false),
			Severity:              securityinsight.AlertSeverity(d.Get("severity").(string)),
			Enabled:               utils.Bool(d.Get("enabled").(bool)),
			Query:                 utils.String(d.Get("query").(string)),
			SuppressionEnabled:    &suppressionEnabled,
			SuppressionDuration:   &suppressionDuration,
		},
	}

	if v, ok := d.GetOk("alert_rule_template_guid"); ok {
		param.NrtAlertRuleProperties.AlertRuleTemplateName = utils.String(v.(string))
	}
	if v, ok := d.GetOk("alert_rule_template_version"); ok {
		param.NrtAlertRuleProperties.TemplateVersion = utils.String(v.(string))
	}
	if v, ok := d.GetOk("event_grouping"); ok {
		param.NrtAlertRuleProperties.EventGroupingSettings = expandAlertRuleScheduledEventGroupingSetting(v.([]interface{}))
	}
	if v, ok := d.GetOk("alert_details_override"); ok {
		param.NrtAlertRuleProperties.AlertDetailsOverride = expandAlertRuleAlertDetailsOverride(v.([]interface{}))
	}
	if v, ok := d.GetOk("custom_details"); ok {
		param.NrtAlertRuleProperties.CustomDetails = utils.ExpandMapStringPtrString(v.(map[string]interface{}))
	}

	entityMappingCount := 0
	sentinelEntityMappingCount := 0
	if v, ok := d.GetOk("entity_mapping"); ok {
		param.NrtAlertRuleProperties.EntityMappings = expandAlertRuleEntityMapping(v.([]interface{}))
		entityMappingCount = len(*param.NrtAlertRuleProperties.EntityMappings)
	}
	if v, ok := d.GetOk("sentinel_entity_mapping"); ok {
		param.NrtAlertRuleProperties.SentinelEntitiesMappings = expandAlertRuleSentinelEntityMapping(v.([]interface{}))
		sentinelEntityMappingCount = len(*param.NrtAlertRuleProperties.SentinelEntitiesMappings)
	}

	// the max number of `sentinel_entity_mapping` and `entity_mapping` together is 5
	if entityMappingCount+sentinelEntityMappingCount > 5 {
		return fmt.Errorf("`entity_mapping` and `sentinel_entity_mapping` together can't exceed 5")
	}

	// Service avoid concurrent update of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindNRT); err != nil {
			return fmt.Errorf("asserting %q: %+v", id, err)
		}
		param.Etag = resp.Value.(securityinsight.NrtAlertRule).Etag
	}

	if _, err := client.CreateOrUpdate(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name, param); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleNrtRead(d, meta)
}

func resourceSentinelAlertRuleNrtRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindNRT); err != nil {
		return fmt.Errorf("asserting %q: %+v", id, err)
	}
	rule := resp.Value.(securityinsight.NrtAlertRule)

	d.Set("name", id.Name)

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if prop := rule.NrtAlertRuleProperties; prop != nil {
		d.Set("description", prop.Description)
		d.Set("display_name", prop.DisplayName)
		if err := d.Set("tactics", flattenAlertRuleTactics(prop.Tactics)); err != nil {
			return fmt.Errorf("setting `tactics`: %+v", err)
		}
		if err := d.Set("techniques", prop.Techniques); err != nil {
			return fmt.Errorf("setting `techniques`: %+v", err)
		}
		if err := d.Set("incident", flattenAlertRuleIncidentConfiguration(prop.IncidentConfiguration, "create_incident_enabled", false)); err != nil {
			return fmt.Errorf("setting `incident`: %+v", err)
		}
		d.Set("severity", string(prop.Severity))
		d.Set("enabled", prop.Enabled)
		d.Set("query", prop.Query)

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

func resourceSentinelAlertRuleNrtDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Nrt %q: %+v", id, err)
	}

	return nil
}
