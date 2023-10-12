// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSentinelAlertRuleNrt() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleNrtCreateUpdate,
		Read:   resourceSentinelAlertRuleNrtRead,
		Update: resourceSentinelAlertRuleNrtCreateUpdate,
		Delete: resourceSentinelAlertRuleNrtDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := alertrules.ParseAlertRuleID(id)
			return err
		}, importSentinelAlertRule(alertrules.AlertRuleKindNRT)),

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
				ValidateFunc: alertrules.ValidateWorkspaceID,
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

			// lintignore:S013
			"event_grouping": {
				Type:     pluginsdk.TypeList,
				Required: features.FourPointOhBeta(),
				Optional: !features.FourPointOhBeta(),
				Computed: !features.FourPointOhBeta(), // the service will default it to `SingleAlert`.
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"aggregation_method": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForEventGroupingAggregationKind(), false),
						},
					},
				},
			},

			"tactics": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAttackTactic(), false),
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
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      alertrules.MatchingMethodAnyAlert,
										ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForMatchingMethod(), false),
									},
									"by_entities": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForEntityMappingType(), false),
										},
									},
									"by_alert_details": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertDetail(), false),
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertSeverity(), false),
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
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertProperty(), false),
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
							ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForEntityMappingType(), false),
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
	workspaceID, err := alertrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_nrt", id.ID())
		}
	}

	param := alertrules.NrtAlertRule{
		Properties: &alertrules.NrtAlertRuleProperties{
			Description:           utils.String(d.Get("description").(string)),
			DisplayName:           d.Get("display_name").(string),
			Techniques:            expandAlertRuleTechnicals(d.Get("techniques").(*pluginsdk.Set).List()),
			Tactics:               expandAlertRuleTactics(d.Get("tactics").(*pluginsdk.Set).List()),
			IncidentConfiguration: expandAlertRuleIncidentConfiguration(d.Get("incident").([]interface{}), "create_incident_enabled", false),
			Severity:              alertrules.AlertSeverity(d.Get("severity").(string)),
			Enabled:               d.Get("enabled").(bool),
			Query:                 d.Get("query").(string),
			SuppressionEnabled:    d.Get("suppression_enabled").(bool),
			SuppressionDuration:   d.Get("suppression_duration").(string),
		},
	}

	if v, ok := d.GetOk("alert_rule_template_guid"); ok {
		param.Properties.AlertRuleTemplateName = utils.String(v.(string))
	}
	if v, ok := d.GetOk("alert_rule_template_version"); ok {
		param.Properties.TemplateVersion = utils.String(v.(string))
	}
	if v, ok := d.GetOk("event_grouping"); ok {
		param.Properties.EventGroupingSettings = expandAlertRuleScheduledEventGroupingSetting(v.([]interface{}))
	}
	if v, ok := d.GetOk("alert_details_override"); ok {
		param.Properties.AlertDetailsOverride = expandAlertRuleAlertDetailsOverride(v.([]interface{}))
	}
	if v, ok := d.GetOk("custom_details"); ok {
		param.Properties.CustomDetails = utils.ExpandPtrMapStringString(v.(map[string]interface{}))
	}

	entityMappingCount := 0
	sentinelEntityMappingCount := 0
	if v, ok := d.GetOk("entity_mapping"); ok {
		param.Properties.EntityMappings = expandAlertRuleEntityMapping(v.([]interface{}))
		entityMappingCount = len(*param.Properties.EntityMappings)
	}
	if v, ok := d.GetOk("sentinel_entity_mapping"); ok {
		param.Properties.SentinelEntitiesMappings = expandAlertRuleSentinelEntityMapping(v.([]interface{}))
		sentinelEntityMappingCount = len(*param.Properties.SentinelEntitiesMappings)
	}

	// the max number of `sentinel_entity_mapping` and `entity_mapping` together is 5
	if entityMappingCount+sentinelEntityMappingCount > 5 {
		return fmt.Errorf("`entity_mapping` and `sentinel_entity_mapping` together can't exceed 5")
	}

	// Service avoid concurrent update of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindNRT); err != nil {
			return fmt.Errorf("asserting %q: %+v", id, err)
		}
		if model := resp.Model; model != nil {
			modelPtr := *model
			param.Etag = modelPtr.(alertrules.NrtAlertRule).Etag
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleNrtRead(d, meta)
}

func resourceSentinelAlertRuleNrtRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindNRT); err != nil {
		return fmt.Errorf("asserting %q: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		modelPtr := *model
		rule := modelPtr.(alertrules.NrtAlertRule)

		d.Set("name", id.RuleId)

		workspaceId := alertrules.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
		d.Set("log_analytics_workspace_id", workspaceId.ID())

		if prop := rule.Properties; prop != nil {
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
			if err := d.Set("custom_details", utils.FlattenPtrMapStringString(prop.CustomDetails)); err != nil {
				return fmt.Errorf("setting `custom_details`: %+v", err)
			}
			if err := d.Set("entity_mapping", flattenAlertRuleEntityMapping(prop.EntityMappings)); err != nil {
				return fmt.Errorf("setting `entity_mapping`: %+v", err)
			}
			if err := d.Set("sentinel_entity_mapping", flattenAlertRuleSentinelEntityMapping(prop.SentinelEntitiesMappings)); err != nil {
				return fmt.Errorf("setting `sentinel_entity_mapping`: %+v", err)
			}
		}
	}

	return nil
}

func resourceSentinelAlertRuleNrtDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Nrt %q: %+v", id, err)
	}

	return nil
}
