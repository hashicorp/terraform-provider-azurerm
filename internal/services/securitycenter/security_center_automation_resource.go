// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: 4.0 - remove these and use the SDK constants instead
const (
	typeLogicApp     = "logicapp"
	typeEventHub     = "eventhub"
	typeLogAnalytics = "loganalytics"
)

func resourceSecurityCenterAutomation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterAutomationCreateUpdate,
		Read:   resourceSecurityCenterAutomationRead,
		Update: resourceSecurityCenterAutomationCreateUpdate,
		Delete: resourceSecurityCenterAutomationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomationID(id)
			return err
		}),

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

			"location": {
				Type:      pluginsdk.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azure.NormalizeLocation,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Default:  true,
				Optional: true,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"scopes": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: commonids.ValidateScopeID,
				},
			},

			"action": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								// TODO: 4.0 - remove these and use the SDK constants instead
								typeLogicApp,
								typeLogAnalytics,
								typeEventHub,
							}, false),
						},

						"resource_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"trigger_url": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},

						"connection_string": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"source": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"event_source": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(automations.EventSourceAlerts),
								string(automations.EventSourceAssessments),
								string(automations.EventSourceAssessmentsSnapshot),
								string(automations.EventSourceRegulatoryComplianceAssessment),
								string(automations.EventSourceRegulatoryComplianceAssessmentSnapshot),
								string(automations.EventSourceSecureScoreControls),
								string(automations.EventSourceSecureScoreControlsSnapshot),
								string(automations.EventSourceSecureScores),
								string(automations.EventSourceSecureScoresSnapshot),
								string(automations.EventSourceSubAssessments),
								string(automations.EventSourceSubAssessmentsSnapshot),
							}, false),
						},

						"rule_set": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"rule": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"property_path": {
													Type:     pluginsdk.TypeString,
													Required: true,
												},
												"expected_value": {
													Type:     pluginsdk.TypeString,
													Required: true,
												},
												"operator": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(automations.OperatorContains),
														string(automations.OperatorEndsWith),
														string(automations.OperatorEquals),
														string(automations.OperatorGreaterThan),
														string(automations.OperatorGreaterThanOrEqualTo),
														string(automations.OperatorLesserThan),
														string(automations.OperatorLesserThanOrEqualTo),
														string(automations.OperatorNotEquals),
														string(automations.OperatorStartsWith),
													}, false),
												},
												"property_type": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(automations.PropertyTypeInteger),
														string(automations.PropertyTypeString),
														string(automations.PropertyTypeBoolean),
														string(automations.PropertyTypeNumber),
													}, false),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceSecurityCenterAutomationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutomationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAutomationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	automationId := automations.NewAutomationID(id.SubscriptionId, id.ResourceGroup, id.Name)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, automationId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_security_center_automation", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	enabled := d.Get("enabled").(bool)

	// Build automation struct
	automation := automations.Automation{
		Location: &location,
		Properties: &automations.AutomationProperties{
			Description: utils.String(d.Get("description").(string)),
			IsEnabled:   &enabled,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	automation.Properties.Scopes = expandSecurityCenterAutomationScopes(d.Get("scopes").([]interface{}))

	var err error
	automation.Properties.Actions, err = expandSecurityCenterAutomationActions(d.Get("action").([]interface{}))
	if err != nil {
		return err
	}

	automation.Properties.Sources, err = expandSecurityCenterAutomationSources(d.Get("source").([]interface{}))
	if err != nil {
		return err
	}

	if _, err := client.CreateOrUpdate(ctx, automationId, automation); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSecurityCenterAutomationRead(d, meta)
}

func resourceSecurityCenterAutomationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutomationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomationID(d.Id())
	if err != nil {
		return err
	}
	automationId := automations.NewAutomationID(id.SubscriptionId, id.ResourceGroup, id.Name)

	resp, err := client.Get(ctx, automationId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Model.Location))

	if properties := resp.Model.Properties; properties != nil {
		d.Set("description", properties.Description)
		d.Set("enabled", properties.IsEnabled)

		flatScopes, err := flattenSecurityCenterAutomationScopes(properties.Scopes)
		if err != nil {
			return err
		}
		if err := d.Set("scopes", flatScopes); err != nil {
			return fmt.Errorf("reading Security Center automation scopes: %+v", err)
		}

		flatActions, err := flattenSecurityCenterAutomationActions(properties.Actions, d)
		if err != nil {
			return err
		}
		if err = d.Set("action", flatActions); err != nil {
			return fmt.Errorf("reading Security Center automation actions: %+v", err)
		}

		flatSources, err := flattenSecurityCenterAutomationSources(properties.Sources)
		if err != nil {
			return err
		}
		if err = d.Set("source", flatSources); err != nil {
			return fmt.Errorf("reading Security Center automation sources: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Model.Tags)
}

func resourceSecurityCenterAutomationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutomationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomationID(d.Id())
	if err != nil {
		return err
	}
	automationId := automations.NewAutomationID(id.SubscriptionId, id.ResourceGroup, id.Name)

	if _, err := client.Delete(ctx, automationId); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandSecurityCenterAutomationSources(sourcesRaw []interface{}) (*[]automations.AutomationSource, error) {
	if len(sourcesRaw) == 0 {
		return &[]automations.AutomationSource{}, nil
	}

	// Output is an array of AutomationSource
	output := make([]automations.AutomationSource, 0)

	// Top level loop over sources array
	for _, sourceRaw := range sourcesRaw {
		sourceMap, ok := sourceRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Security Center automation, unable to decode sources")
		}

		// Build and parse array of RuleSets
		ruleSets := make([]automations.AutomationRuleSet, 0)
		ruleSetsRaw := sourceMap["rule_set"].([]interface{})
		for _, ruleSetRaw := range ruleSetsRaw {
			ruleSetMap := ruleSetRaw.(map[string]interface{})
			rulesRaw := ruleSetMap["rule"].([]interface{})

			// Build and parse array of Rules in each RuleSet
			rules := make([]automations.AutomationTriggeringRule, 0)
			for _, ruleRaw := range rulesRaw {
				// Parse the rule fields
				ruleMap := ruleRaw.(map[string]interface{})
				rulePath := ruleMap["property_path"].(string)
				ruleType := automations.PropertyType(ruleMap["property_type"].(string))
				ruleValue := ruleMap["expected_value"].(string)
				ruleOperator := automations.Operator(ruleMap["operator"].(string))

				// Create AutomationTriggeringRule struct and push into array
				rule := automations.AutomationTriggeringRule{
					PropertyJPath: &rulePath,
					PropertyType:  &ruleType,
					ExpectedValue: &ruleValue,
					Operator:      &ruleOperator,
				}
				rules = append(rules, rule)
			}

			// Create AutomationRuleSet struct and push into array
			ruleSet := automations.AutomationRuleSet{
				Rules: &rules,
			}
			ruleSets = append(ruleSets, ruleSet)
		}

		// Finally create AutomationSource struct holding our list of RuleSets
		eventSource := automations.EventSource(sourceMap["event_source"].(string))
		source := automations.AutomationSource{
			EventSource: &eventSource,
			RuleSets:    &ruleSets,
		}

		// Finally (no really this time), push AutomationSource into output
		output = append(output, source)
	}

	return &output, nil
}

func expandSecurityCenterAutomationScopes(scopePathsRaw []interface{}) *[]automations.AutomationScope {
	scopes := make([]automations.AutomationScope, 0)

	for _, scopePathRaw := range scopePathsRaw {
		if path, ok := scopePathRaw.(string); ok {
			desc := fmt.Sprintf("scope for %s", path)
			scope := automations.AutomationScope{
				ScopePath:   &path,
				Description: &desc,
			}
			scopes = append(scopes, scope)
		}
	}

	return &scopes
}

func expandSecurityCenterAutomationActions(actionsRaw []interface{}) (*[]automations.AutomationAction, error) {
	if len(actionsRaw) == 0 {
		return &[]automations.AutomationAction{}, nil
	}

	output := make([]automations.AutomationAction, 0)

	for _, actionRaw := range actionsRaw {
		actionMap := actionRaw.(map[string]interface{})

		var autoAction automations.AutomationAction
		var resourceID string
		var actionType string
		var ok bool

		// No checking, as fields are enforced by resource schema
		resourceID = actionMap["resource_id"].(string)
		actionType = actionMap["type"].(string)

		// Ignore case on type field
		switch strings.ToLower(actionType) {
		// Handle LogicApp action type
		case typeLogicApp:
			var triggerURL string
			if triggerURL, ok = actionMap["trigger_url"].(string); !ok || triggerURL == "" {
				return nil, fmt.Errorf("Security Center automation, trigger_url is required for LogicApp action")
			}
			autoAction = automations.AutomationActionLogicApp{
				LogicAppResourceId: &resourceID,
				Uri:                &triggerURL,
			}

			// Handle LogAnalytics action type
		case typeLogAnalytics:
			autoAction = automations.AutomationActionWorkspace{
				WorkspaceResourceId: &resourceID,
			}

			// Handle EventHub action type
		case typeEventHub:
			var connString string
			if connString, ok = actionMap["connection_string"].(string); !ok || connString == "" {
				return nil, fmt.Errorf("Security Center automation, connection_string is required for EventHub action")
			}
			autoAction = automations.AutomationActionEventHub{
				EventHubResourceId: &resourceID,
				ConnectionString:   &connString,
			}
		default:
			return nil, fmt.Errorf("Security Center automation, expected action type to be one of: %s, %s or %s", automations.ActionTypeEventHub, automations.ActionTypeWorkspace, automations.ActionTypeLogicApp)
		}
		output = append(output, autoAction)
	}

	return &output, nil
}

func flattenSecurityCenterAutomationSources(sources *[]automations.AutomationSource) ([]map[string]interface{}, error) {
	if sources == nil {
		return make([]map[string]interface{}, 0), nil
	}

	resultSlice := make([]map[string]interface{}, 0)
	for _, source := range *sources {
		ruleSetSlice := make([]interface{}, 0)

		// RuleSets is an optional field need check for nil
		if source.RuleSets != nil {
			for _, ruleSet := range *source.RuleSets {
				ruleSlice := make([]map[string]string, 0)

				for _, rule := range *ruleSet.Rules {
					if rule.PropertyJPath == nil {
						return nil, fmt.Errorf("Security Center automation, API returned a rule with an empty PropertyJPath")
					}
					if rule.ExpectedValue == nil {
						return nil, fmt.Errorf("Security Center automation, API returned a rule with empty ExpectedValue")
					}
					ruleMap := map[string]string{
						"property_path":  *rule.PropertyJPath,
						"expected_value": *rule.ExpectedValue,
						"operator":       string(*rule.Operator),
						"property_type":  string(*rule.PropertyType),
					}
					ruleSlice = append(ruleSlice, ruleMap)
				}

				ruleSetMap := map[string]interface{}{
					"rule": ruleSlice,
				}
				ruleSetSlice = append(ruleSetSlice, ruleSetMap)
			}
		}

		sourceMap := map[string]interface{}{
			"event_source": source.EventSource,
			"rule_set":     ruleSetSlice,
		}
		resultSlice = append(resultSlice, sourceMap)
	}

	return resultSlice, nil
}

func flattenSecurityCenterAutomationScopes(scopes *[]automations.AutomationScope) ([]string, error) {
	if scopes == nil {
		return []string{}, nil
	}

	resultSlice := make([]string, 0)
	for _, scope := range *scopes {
		if scope.ScopePath == nil {
			return nil, fmt.Errorf("Security Center automation, API returned a scope with an empty ScopePath")
		}

		resultSlice = append(resultSlice, *scope.ScopePath)
	}

	return resultSlice, nil
}

func flattenSecurityCenterAutomationActions(actions *[]automations.AutomationAction, d *pluginsdk.ResourceData) ([]map[string]string, error) {
	if actions == nil {
		return []map[string]string{}, nil
	}

	resultSlice := make([]map[string]string, 0)

	for i, action := range *actions {
		// Use type assertion to discover the underlying action
		// Trying to use action.(automations.AutomationAction).ActionType results in a panic
		actionLogicApp, isLogicApp := action.(automations.AutomationActionLogicApp)
		if isLogicApp {
			if actionLogicApp.LogicAppResourceId == nil {
				return nil, fmt.Errorf("Security Center automation, API returned an action with empty logicAppResourceId")
			}
			actionMap := map[string]string{
				"resource_id": *actionLogicApp.LogicAppResourceId,
				"type":        string(typeLogicApp),
				"trigger_url": "",
			}

			// Need to merge in trigger_url as it's not returned by API Get operation
			// Otherwise don't have consistent state
			if triggerURL, ok := d.GetOk(fmt.Sprintf("action.%d.trigger_url", i)); ok {
				actionMap["trigger_url"] = triggerURL.(string)
			}

			resultSlice = append(resultSlice, actionMap)
		}

		actionEventHub, isEventHub := action.(automations.AutomationActionEventHub)
		if isEventHub {
			if actionEventHub.EventHubResourceId == nil {
				return nil, fmt.Errorf("Security Center automation, API returned an action with empty eventHubResourceId")
			}
			actionMap := map[string]string{
				"resource_id":       *actionEventHub.EventHubResourceId,
				"type":              string(typeEventHub),
				"connection_string": "",
			}

			// Need to merge in connection_string as it's not returned by API Get operation
			// Otherwise don't have consistent state
			if triggerURL, ok := d.GetOk(fmt.Sprintf("action.%d.connection_string", i)); ok {
				actionMap["connection_string"] = triggerURL.(string)
			}

			resultSlice = append(resultSlice, actionMap)
		}

		actionLogAnalytics, isLogAnalytics := action.(automations.AutomationActionWorkspace)
		if isLogAnalytics {
			if actionLogAnalytics.WorkspaceResourceId == nil {
				return nil, fmt.Errorf("Security Center automation, API returned an action with empty workspaceResourceId")
			}
			actionMap := map[string]string{
				"resource_id": *actionLogAnalytics.WorkspaceResourceId,
				"type":        string(typeLogAnalytics),
			}

			resultSlice = append(resultSlice, actionMap)
		}
	}

	return resultSlice, nil
}
