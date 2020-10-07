package securitycenter

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/azuresdkhacks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const typeLogicApp = "logicapp"
const typeEventHub = "eventhub"
const typeLogAnalytics = "loganalytics"

func resourceArmSecurityCenterAutomation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterAutomationCreateUpdate,
		Read:   resourceArmSecurityCenterAutomationRead,
		Update: resourceArmSecurityCenterAutomationCreateUpdate,
		Delete: resourceArmSecurityCenterAutomationDelete,

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azure.NormalizeLocation,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"scopes": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"action": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{typeLogicApp, typeLogAnalytics, typeEventHub}, true),
						},

						"resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"trigger_url": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},

						"connection_string": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"source": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_source": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Alerts", "Assessments", "SubAssessments"}, true),
						},

						"rule_set": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"property_path": {
													Type:     schema.TypeString,
													Required: true,
												},
												"expected_value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Required: true,
												},
												"property_type": {
													Type:     schema.TypeString,
													Required: true,
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
		},
	}
}

func resourceArmSecurityCenterAutomationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutomationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// Core resource props
	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Security Center automation %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_automation", *existing.ID)
		}
	}

	enabled := d.Get("enabled").(bool)
	description := fmt.Sprintf("Created by Terraform")

	// Build automation struct
	automation := security.Automation{
		Location: &location,
		AutomationProperties: &security.AutomationProperties{
			IsEnabled:   &enabled,
			Description: &description,
		},
	}

	automation.AutomationProperties.Scopes = expandScopes(d.Get("scopes").([]interface{}))

	var err error
	automation.AutomationProperties.Actions, err = expandActions(d.Get("action").([]interface{}))
	if err != nil {
		return err
	}

	automation.AutomationProperties.Sources = expandSources(d.Get("source").([]interface{}))
	if err != nil {
		return err
	}

	// Create our patched/hacked struct with real struct embedded
	patchedAutomation := azuresdkhacks.Automation{
		Automation: automation,
	}

	resp, err := client.CreateOrUpdate(ctx, resGroup, name, patchedAutomation)
	if err != nil {
		return fmt.Errorf("Error creating Security Center automation: %+v", err)
	}

	// Important steps
	d.SetId(*resp.ID)
	return resourceArmSecurityCenterAutomationRead(d, meta)
}

func resourceArmSecurityCenterAutomationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutomationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["automations"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Security Center automation %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Security Center automation %s: %v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if properties := resp.AutomationProperties; properties != nil {
		d.Set("enabled", properties.IsEnabled)

		if err := d.Set("scopes", flattenScopes(properties.Scopes)); err != nil {
			return fmt.Errorf("Error reading Security Center automation scopes: %+v", err)
		}

		if err = d.Set("action", flattenActions(properties.Actions, d)); err != nil {
			return fmt.Errorf("Error reading Security Center automation actions: %+v", err)
		}

		if err = d.Set("source", flattenSources(properties.Sources)); err != nil {
			return fmt.Errorf("Error reading Security Center automation sources: %+v", err)
		}
	}

	return nil
}

func resourceArmSecurityCenterAutomationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutomationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["automations"]

	resp, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Security Center automation was not found: %v", err)
			return nil
		}
		return fmt.Errorf("Error deleting Security Center automation: %+v", err)
	}

	return nil
}

func expandSources(sourcesRaw []interface{}) *[]security.AutomationSource {
	f, err := os.OpenFile("provider_debug5.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	if len(sourcesRaw) == 0 {
		return &[]security.AutomationSource{}
	}

	// Output is an array of AutomationSource
	output := make([]security.AutomationSource, 0)

	// Top level loop over sources array
	for _, sourceRaw := range sourcesRaw {
		sourceMap, ok := sourceRaw.(map[string]interface{})
		if !ok {
			continue
		}

		// Build and parse array of RuleSets
		ruleSets := make([]security.AutomationRuleSet, 0)
		ruleSetsRaw := sourceMap["rule_set"].([]interface{})
		for _, ruleSetRaw := range ruleSetsRaw {
			ruleSetMap := ruleSetRaw.(map[string]interface{})
			rulesRaw := ruleSetMap["rule"].([]interface{})

			// Build and parse array of Rules in each RuleSet
			rules := make([]security.AutomationTriggeringRule, 0)
			for _, ruleRaw := range rulesRaw {
				// Parse the rule fields
				ruleMap := ruleRaw.(map[string]interface{})
				rulePath := ruleMap["property_path"].(string)
				ruleType := security.PropertyType(ruleMap["property_type"].(string))
				ruleValue := ruleMap["expected_value"].(string)
				ruleOperator := security.Operator(ruleMap["operator"].(string))

				// Create AutomationTriggeringRule struct and push into array
				rule := security.AutomationTriggeringRule{
					PropertyJPath: &rulePath,
					PropertyType:  ruleType,
					ExpectedValue: &ruleValue,
					Operator:      ruleOperator,
				}
				rules = append(rules, rule)
			}

			// Create AutomationRuleSet struct and push into array
			ruleSet := security.AutomationRuleSet{
				Rules: &rules,
			}
			ruleSets = append(ruleSets, ruleSet)
		}

		// Finally create AutomationSource struct holding our list of RuleSets
		eventSource := security.EventSource(sourceMap["event_source"].(string))
		source := security.AutomationSource{
			EventSource: eventSource,
			RuleSets:    &ruleSets,
		}

		// Finally (no really this time), push AutomationSource into output
		output = append(output, source)
	}

	return &output
}

func expandScopes(scopePathsRaw []interface{}) *[]security.AutomationScope {
	scopes := make([]security.AutomationScope, 0)

	for _, scopePathRaw := range scopePathsRaw {
		if path, ok := scopePathRaw.(string); ok {
			desc := fmt.Sprintf("Scope for %s", path)
			scope := security.AutomationScope{
				ScopePath:   &path,
				Description: &desc,
			}
			scopes = append(scopes, scope)
		}
	}

	return &scopes
}

func expandActions(actionsRaw []interface{}) (*[]security.BasicAutomationAction, error) {
	if len(actionsRaw) == 0 {
		return &[]security.BasicAutomationAction{}, nil
	}

	output := make([]security.BasicAutomationAction, 0)

	for _, actionRaw := range actionsRaw {
		actionMap := actionRaw.(map[string]interface{})

		var autoAction security.BasicAutomationAction
		var resourceID string
		var actionType string
		var ok bool

		// All action types require a resource_id
		if resourceID, ok = actionMap["resource_id"].(string); !ok || resourceID == "" {
			return nil, fmt.Errorf("Security Center automation, resource_id is required for action")
		}

		// Type is important and required
		if actionType, ok = actionMap["type"].(string); !ok || actionType == "" {
			return nil, fmt.Errorf("Security Center automation, type is required for action")
		}

		// Ignore case on type field
		switch strings.ToLower(actionType) {
		// Handle LogicApp action type
		case typeLogicApp:
			var triggerURL string
			if triggerURL, ok = actionMap["trigger_url"].(string); !ok || triggerURL == "" {
				return nil, fmt.Errorf("Security Center automation, trigger_url is required for LogicApp action")
			}
			autoAction = security.AutomationActionLogicApp{
				LogicAppResourceID: &resourceID,
				URI:                &triggerURL,
				ActionType:         security.ActionTypeLogicApp,
			}

		// Handle LogAnalytics action type
		case typeLogAnalytics:
			autoAction = security.AutomationActionWorkspace{
				WorkspaceResourceID: &resourceID,
				ActionType:          security.ActionTypeWorkspace,
			}

		// Handle EventHub action type
		case typeEventHub:
			var connString string
			if connString, ok = actionMap["connection_string"].(string); !ok || connString == "" {
				return nil, fmt.Errorf("Security Center automation, connection_string is required for EventHub action")
			}
			autoAction = security.AutomationActionEventHub{
				EventHubResourceID: &resourceID,
				ConnectionString:   &connString,
				ActionType:         security.ActionTypeEventHub,
			}
		default:
			continue
		}
		output = append(output, autoAction)
	}

	return &output, nil
}

func flattenSources(sources *[]security.AutomationSource) []map[string]interface{} {
	if sources == nil {
		return make([]map[string]interface{}, 0)
	}

	resultSlice := make([]map[string]interface{}, 0)
	for _, source := range *sources {
		ruleSetSlice := make([]interface{}, 0)

		// RuleSets is an optional field need check for nil
		if source.RuleSets != nil {
			for _, ruleSet := range *source.RuleSets {
				ruleSlice := make([]map[string]string, 0)

				for _, rule := range *ruleSet.Rules {
					ruleMap := map[string]string{
						"property_path":  *rule.PropertyJPath,
						"expected_value": *rule.ExpectedValue,
						"operator":       string(rule.Operator),
						"property_type":  string(rule.PropertyType),
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

	return resultSlice
}

func flattenScopes(scopes *[]security.AutomationScope) []string {
	if scopes == nil {
		return []string{}
	}

	resultSlice := make([]string, 0)
	for _, scope := range *scopes {
		if scope.ScopePath == nil {
			continue
		}

		resultSlice = append(resultSlice, *scope.ScopePath)
	}

	return resultSlice
}

func flattenActions(actions *[]security.BasicAutomationAction, d *schema.ResourceData) []map[string]string {
	if actions == nil {
		return []map[string]string{}
	}

	// Get existing schema data for actions
	schemaDataActions, schemaDataActionsOk := d.GetOk("action")

	resultSlice := make([]map[string]string, 0)

	for i, action := range *actions {
		// Use type assertion to discover the underlying action
		actionLogicApp, isLogicApp := action.(security.AutomationActionLogicApp)
		if isLogicApp {
			actionMap := map[string]string{
				"resource_id": *actionLogicApp.LogicAppResourceID,
				"type":        "LogicApp",
				"trigger_url": "",
			}

			// Need to merge in trigger_url as it's not returned by API Get operation
			// Otherwise don't have consistent state
			if schemaDataActionsOk {
				actionsSlice := schemaDataActions.([]interface{})
				dataAction := actionsSlice[i].(map[string]interface{})
				if dataAction["trigger_url"].(string) != "" {
					actionMap["trigger_url"] = dataAction["trigger_url"].(string)
				}
			}

			resultSlice = append(resultSlice, actionMap)
		}

		actionEventHub, isEventHub := action.(security.AutomationActionEventHub)
		if isEventHub {
			actionMap := map[string]string{
				"resource_id":       *actionEventHub.EventHubResourceID,
				"type":              "EventHub",
				"connection_string": "",
			}

			// Need to merge in connection_string as it's not returned by API Get operation
			// Otherwise don't have consistent state
			if schemaDataActionsOk {
				actionsSlice := schemaDataActions.([]interface{})
				dataAction := actionsSlice[i].(map[string]interface{})
				if dataAction["connection_string"].(string) != "" {
					actionMap["connection_string"] = dataAction["connection_string"].(string)
				}
			}

			resultSlice = append(resultSlice, actionMap)
		}

		actionLogAnalytics, isLogAnalytics := action.(security.AutomationActionWorkspace)
		if isLogAnalytics {
			actionMap := map[string]string{
				"resource_id": *actionLogAnalytics.WorkspaceResourceID,
				"type":        "LogAnalytics",
			}
			resultSlice = append(resultSlice, actionMap)
		}
	}

	return resultSlice
}
