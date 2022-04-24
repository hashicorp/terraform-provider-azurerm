package frontdoor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/sdk/2020-05-01/frontdoors"
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

		SchemaVersion: 1,

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

			"resource_group_name": azure.SchemaResourceGroupName(),

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
										Default:  !features.ThreePointOhBeta(),
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

									"routing_rule_override": {
										Type:     pluginsdk.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{

												"redirect_configuration": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"custom_fragment": {
																Type:     pluginsdk.TypeString,
																Optional: true,
															},
															"custom_host": {
																Type:     pluginsdk.TypeString,
																Optional: true,
															},
															"custom_path": {
																Type:     pluginsdk.TypeString,
																Optional: true,
															},
															"custom_query_string": {
																Type:     pluginsdk.TypeString,
																Optional: true,
															},
															"redirect_protocol": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(frontdoors.FrontDoorRedirectProtocolHttpOnly),
																	string(frontdoors.FrontDoorRedirectProtocolHttpsOnly),
																	string(frontdoors.FrontDoorRedirectProtocolMatchRequest),
																}, false),
															},
															"redirect_type": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(frontdoors.FrontDoorRedirectTypeFound),
																	string(frontdoors.FrontDoorRedirectTypeMoved),
																	string(frontdoors.FrontDoorRedirectTypePermanentRedirect),
																	string(frontdoors.FrontDoorRedirectTypeTemporaryRedirect),
																}, false),
															},
														},
													},
												},
												"forwarding_configuration": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"backend_pool_name": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: azValidate.BackendPoolRoutingRuleName,
															},
															"cache_enabled": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																Default:  false,
															},
															"cache_use_dynamic_compression": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																Default:  false,
															},
															"cache_query_parameter_strip_directive": {
																Type:     pluginsdk.TypeString,
																Optional: true,
																Default:  string(frontdoors.FrontDoorQueryStripAll),
																ValidateFunc: validation.StringInSlice([]string{
																	string(frontdoors.FrontDoorQueryStripAll),
																	string(frontdoors.FrontDoorQueryStripNone),
																	string(frontdoors.FrontDoorQueryStripOnly),
																	string(frontdoors.FrontDoorQueryStripAllExcept),
																}, false),
															},
															"cache_query_parameters": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																MaxItems: 25,
																Elem: &pluginsdk.Schema{
																	Type:         pluginsdk.TypeString,
																	ValidateFunc: validation.StringIsNotEmpty,
																},
															},
															"cache_duration": {
																Type:         pluginsdk.TypeString,
																Optional:     true,
																ValidateFunc: validate.ISO8601DurationBetween("PT1S", "P365D"),
															},
															"custom_forwarding_path": {
																Type:     pluginsdk.TypeString,
																Optional: true,
															},
															"forwarding_protocol": {
																Type:     pluginsdk.TypeString,
																Optional: true,
																Default:  string(frontdoors.FrontDoorForwardingProtocolHttpsOnly),
																ValidateFunc: validation.StringInSlice([]string{
																	string(frontdoors.FrontDoorForwardingProtocolHttpOnly),
																	string(frontdoors.FrontDoorForwardingProtocolHttpsOnly),
																	string(frontdoors.FrontDoorForwardingProtocolMatchRequest),
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
					},
				},
			},

			// Computed values --> Produces this error 'testcase.go:110: Step 1/2 error: After applying this test step, the plan was not empty.'
			"explicit_resource_order": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"engine_rule_ids": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
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
	frontDoorId := frontdoors.NewFrontDoorID(subscriptionId, resourceGroup, frontDoorName)
	rulesEngineName := d.Get("name").(string)

	rules := d.Get("rule").([]interface{})

	id := frontdoors.NewRulesEngineID(subscriptionId, resourceGroup, frontDoorName, rulesEngineName)

	frontdoorRulesEngineProperties := frontdoors.RulesEngineProperties{
		Rules: expandFrontDoorRulesEngineRules(frontDoorId, rules),
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

func expandFrontDoorRulesEngineAction(frontDoorId frontdoors.FrontDoorId, input []interface{}) frontdoors.RulesEngineAction {
	if len(input) == 0 || input[0] == nil {
		return frontdoors.RulesEngineAction{}
	}

	ruleAction := input[0].(map[string]interface{})
	requestHeaderActions := ruleAction["request_header"].([]interface{})
	responseHeaderActions := ruleAction["response_header"].([]interface{})

	var routingConfiguration frontdoors.RouteConfiguration
	if rco := ruleAction["routing_rule_override"].([]interface{}); len(rco) != 0 {
		routeConfigurationOverride := ruleAction["routing_rule_override"].([]interface{})[0].(map[string]interface{})
		if rc := routeConfigurationOverride["redirect_configuration"].([]interface{}); len(rc) != 0 {
			routingConfiguration = expandFrontDoorRedirectConfiguration(rc)
		} else if fc := routeConfigurationOverride["forwarding_configuration"].([]interface{}); len(fc) != 0 {
			routingConfiguration = expandFrontDoorForwardingConfiguration(fc, frontDoorId)
		}
	}

	frontdoorRulesEngineRuleAction := frontdoors.RulesEngineAction{
		RequestHeaderActions:       expandHeaderAction(requestHeaderActions),
		ResponseHeaderActions:      expandHeaderAction(responseHeaderActions),
		RouteConfigurationOverride: &routingConfiguration,
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

func expandFrontDoorRulesEngineRules(frontDoorId frontdoors.FrontDoorId, input []interface{}) *[]frontdoors.RulesEngineRule {
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
			Action:          expandFrontDoorRulesEngineAction(frontDoorId, actions),
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

	d.Set("name", id.RulesEngineName)
	d.Set("frontdoor_name", id.FrontDoorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			explicitResourceOrder := d.Get("explicit_resource_order").([]interface{})

			var flattenedRoutingRules *[]interface{}

			flattenedRoutingRules, err = flattenFrontDoorRulesEngineRule(props.Rules, d.Get("routing_rule_override"), *id, explicitResourceOrder)
			if err != nil {
				return fmt.Errorf("flattening `routing_rule_override`: %+v", err)
			}
			d.Set("rule", flattenedRoutingRules)
		}
	}

	return nil
}

func flattenFrontDoorRulesEngineRule(input *[]frontdoors.RulesEngineRule, oldBlocks interface{}, frontDoorRulesEngineId frontdoors.RulesEngineId, explicitOrder []interface{}) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)

	if len(explicitOrder) > 0 {
		orderedRule := explicitOrder[0].(map[string]interface{})
		orderedRulesEngineRuleIds := orderedRule["engine_rule_ids"].([]interface{})
		combinedRulesEngineRule, err := combineRulesEngineRule(*input, oldBlocks, orderedRulesEngineRuleIds, frontDoorRulesEngineId)
		if err != nil {
			return nil, err
		}
		output = combinedRulesEngineRule
	} else {
		for _, v := range *input {
			rulesEngineRule, err := flattenSingleFrontDoorRulesEngineRule(v, oldBlocks, frontDoorRulesEngineId)
			if err == nil {
				output = append(output, rulesEngineRule)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontDoorRulesEngineRule(input frontdoors.RulesEngineRule, oldBlocks interface{}, frontDoorRulesEngineId frontdoors.RulesEngineId) (map[string]interface{}, error) {

	name := ""
	if input.Name != "" {
		// rewrite the ID to ensure it's consistent
		name = input.Name
	}

	action, err := flattenSingleFrontDooRulesEngineRuleAction(input.Action, oldBlocks)
	if err != nil {
		return nil, fmt.Errorf("flattening `action`: %+v", err)
	}

	output := map[string]interface{}{
		"name":     name,
		"priority": input.Priority,
		"action":   action,
	}

	return output, nil
}

func flattenSingleFrontDooRulesEngineRuleAction(input frontdoors.RulesEngineAction, oldBlocks interface{}) ([]interface{}, error) {
	forwardingConfiguration := make([]interface{}, 0)
	redirectConfiguration := make([]interface{}, 0)

	forwardConfiguration, err := flattenRoutingRuleForwardingConfiguration(input.RouteConfigurationOverride, oldBlocks)
	if err != nil {
		return nil, fmt.Errorf("flattening `forward_configuration`: %+v", err)
	}
	forwardingConfiguration = *forwardConfiguration
	redirectConfiguration = flattenRoutingRuleRedirectConfiguration(input.RouteConfigurationOverride)

	output := make([]interface{}, 0)
	overrides := make([]interface{}, 0)
	override := map[string]interface{}{
		"forwarding_configuration": forwardingConfiguration,
		"redirect_configuration":   redirectConfiguration,
	}
	overrides = append(overrides, override)
	block := map[string]interface{}{
		"routing_rule_override": overrides,
	}
	output = append(output, block)
	return output, nil
}

func combineRulesEngineRule(allRulesEngineRule []frontdoors.RulesEngineRule, oldBlocks interface{}, orderedIds []interface{}, frontDoorRulesEngineId frontdoors.RulesEngineId) ([]interface{}, error) {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, ruleEngineRule := range allRulesEngineRule {

			// TO LOWER Each in

			if strings.Contains(strings.ToLower(v.(string)), strings.ToLower(*&ruleEngineRule.Name)) {
				orderedRoutingRule, err := flattenSingleFrontDoorRulesEngineRule(ruleEngineRule, oldBlocks, frontDoorRulesEngineId)
				if err == nil {
					output = append(output, orderedRoutingRule)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, ruleEngineRule := range allRulesEngineRule {
		found = false
		for _, orderedId := range orderedIds {
			if strings.Contains(strings.ToLower(orderedId.(string)), strings.ToLower(ruleEngineRule.Name)) {
				found = true
				break
			}
		}

		if !found {
			newRuleEngineRule, err := flattenSingleFrontDoorRulesEngineRule(ruleEngineRule, oldBlocks, frontDoorRulesEngineId)
			if err == nil {
				output = append(output, newRuleEngineRule)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
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
