// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontDoorRulesEngine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorRulesEngineCreateUpdate,
		Read:   resourceFrontDoorRulesEngineRead,
		Update: resourceFrontDoorRulesEngineCreateUpdate,
		Delete: resourceFrontDoorRulesEngineDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.RulesEngineV0ToV1{},
			1: migration.RulesEngineV1ToV2{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RulesEngineID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frontdoor_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.FrontDoorName,
			},
			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"rule": {
				Type:     pluginsdk.TypeList,
				MaxItems: 100,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"priority": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"match_condition": {
							Type:     pluginsdk.TypeList,
							MaxItems: 100,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"variable": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"IsMobile",
											"RemoteAddr",
											"RequestMethod",
											"QueryString",
											"PostArgs",
											"RequestUri",
											"RequestPath",
											"RequestFilename",
											"RequestFilenameExtension",
											"RequestHeader",
											"RequestBody",
											"RequestScheme",
										}, false),
									},

									"selector": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Any",
											"IPMatch",
											"GeoMatch",
											"Equal",
											"Contains",
											"LessThan",
											"GreaterThan",
											"LessThanOrEqual",
											"GreaterThanOrEqual",
											"BeginsWith",
											"EndsWith",
										}, false),
									},

									"transform": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 6,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Lowercase",
												"RemoveNulls",
												"Trim",
												"Uppercase",
												"UrlDecode",
												"UrlEncode",
											}, false),
										},
									},

									"negate_condition": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"value": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 25,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},

						"action": {
							Type:     pluginsdk.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"request_header": {
										Type:     pluginsdk.TypeList,
										MaxItems: 100,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"header_action_type": {
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice([]string{
														string(frontdoors.HeaderActionTypeAppend),
														string(frontdoors.HeaderActionTypeDelete),
														string(frontdoors.HeaderActionTypeOverwrite),
													}, false),
													Optional: true,
												},

												"header_name": {
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
													Optional:     true,
												},

												"value": {
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
													Optional:     true,
												},
											},
										},
									},

									"response_header": {
										Type:     pluginsdk.TypeList,
										MaxItems: 100,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"header_action_type": {
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice([]string{
														string(frontdoors.HeaderActionTypeAppend),
														string(frontdoors.HeaderActionTypeDelete),
														string(frontdoors.HeaderActionTypeOverwrite),
													}, false),
													Optional: true,
												},

												"header_name": {
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
													Optional:     true,
												},

												"value": {
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
													Optional:     true,
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

func resourceFrontDoorRulesEngineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	frontDoorName := d.Get("frontdoor_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	rulesEngineName := d.Get("name").(string)

	rules := d.Get("rule").([]interface{})

	id := frontdoors.NewRulesEngineID(subscriptionId, resourceGroup, frontDoorName, rulesEngineName)

	frontdoorRulesEngineProperties := frontdoors.RulesEngineProperties{
		Rules: expandFrontDoorRulesEngineRules(rules),
	}

	frontdoorRulesEngine := frontdoors.RulesEngine{
		Name:       utils.String(rulesEngineName),
		Properties: &frontdoorRulesEngineProperties,
	}

	if err := client.RulesEnginesCreateOrUpdateThenPoll(ctx, id, frontdoorRulesEngine); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontDoorRulesEngineRead(d, meta)
}

func expandFrontDoorRulesEngineAction(input []interface{}) frontdoors.RulesEngineAction {
	if len(input) == 0 || input[0] == nil {
		return frontdoors.RulesEngineAction{}
	}

	ruleAction := input[0].(map[string]interface{})

	requestHeaderActions := ruleAction["request_header"].([]interface{})
	responseHeaderActions := ruleAction["response_header"].([]interface{})

	frontdoorRulesEngineRuleAction := frontdoors.RulesEngineAction{
		RequestHeaderActions:  expandHeaderAction(requestHeaderActions),
		ResponseHeaderActions: expandHeaderAction(responseHeaderActions),
	}

	return frontdoorRulesEngineRuleAction
}

func expandHeaderAction(input []interface{}) *[]frontdoors.HeaderAction {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	output := make([]frontdoors.HeaderAction, 0)

	for _, a := range input {
		action := a.(map[string]interface{})

		headerName := action["header_name"].(string)
		value := action["value"].(string)
		headerActionType := action["header_action_type"].(string)

		frontdoorRulesEngineRuleHeaderAction := frontdoors.HeaderAction{
			HeaderName:       headerName,
			Value:            utils.String(value),
			HeaderActionType: frontdoors.HeaderActionType(headerActionType),
		}

		output = append(output, frontdoorRulesEngineRuleHeaderAction)
	}

	return &output
}

func expandFrontDoorRulesEngineRules(input []interface{}) *[]frontdoors.RulesEngineRule {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	output := make([]frontdoors.RulesEngineRule, 0)

	for _, r := range input {
		rule := r.(map[string]interface{})

		ruleName := rule["name"].(string)
		priority := int64(rule["priority"].(int))
		actions := rule["action"].([]interface{})
		matchConditions := rule["match_condition"].([]interface{})

		frontdoorRulesEngineRule := frontdoors.RulesEngineRule{
			Name:            ruleName,
			Priority:        priority,
			Action:          expandFrontDoorRulesEngineAction(actions),
			MatchConditions: expandFrontDoorRulesEngineMatchCondition(matchConditions),
		}

		output = append(output, frontdoorRulesEngineRule)
	}
	return &output
}

func expandFrontDoorRulesEngineMatchCondition(input []interface{}) *[]frontdoors.RulesEngineMatchCondition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	output := make([]frontdoors.RulesEngineMatchCondition, 0)

	for _, c := range input {
		condition := c.(map[string]interface{})

		selector := condition["selector"].(string)
		negateCondition := condition["negate_condition"].(bool)
		matchVariable := condition["variable"].(string)
		operator := condition["operator"].(string)
		transform := condition["transform"].([]interface{})
		matchValue := condition["value"].([]interface{})

		matchValueArray := make([]string, 0)
		for _, v := range matchValue {
			matchValueArray = append(matchValueArray, v.(string))
		}

		matchCondition := frontdoors.RulesEngineMatchCondition{
			RulesEngineMatchVariable: frontdoors.RulesEngineMatchVariable(matchVariable),
			Selector:                 utils.String(selector),
			RulesEngineOperator:      frontdoors.RulesEngineOperator(operator),
			NegateCondition:          &negateCondition,
			RulesEngineMatchValue:    matchValueArray,
			Transforms:               expandFrontDoorRulesEngineMatchConditionTransform(transform),
		}
		output = append(output, matchCondition)
	}
	return &output
}

func expandFrontDoorRulesEngineMatchConditionTransform(input []interface{}) *[]frontdoors.Transform {
	if len(input) == 0 || input[0] == nil {
		return &[]frontdoors.Transform{}
	}

	output := make([]frontdoors.Transform, 0)

	for _, t := range input {
		result := frontdoors.Transform(t.(string))

		output = append(output, result)
	}
	return &output
}

func resourceFrontDoorRulesEngineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := frontdoors.ParseRulesEngineIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.RulesEnginesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return nil
}

func resourceFrontDoorRulesEngineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := frontdoors.ParseRulesEngineIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if err := client.RulesEnginesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
