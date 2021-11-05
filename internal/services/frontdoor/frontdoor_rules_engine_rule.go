package frontdoor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-05-01/frontdoor"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
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
			"location": location.SchemaComputed(),

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
											"RequestURI",
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
										Default:  true, // TODO 3,0 change to false- needs to change https://github.com/hashicorp/terraform-provider-azurerm/pull/13605
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
														string(frontdoor.Append),
														string(frontdoor.Delete),
														string(frontdoor.Overwrite),
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
														string(frontdoor.Append),
														string(frontdoor.Delete),
														string(frontdoor.Overwrite),
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
	client := meta.(*clients.Client).Frontdoor.FrontDoorsRulesEnginesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	frontDoorName := d.Get("frontdoor_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	rulesEngineName := d.Get("name").(string)

	rules := d.Get("rule").([]interface{})

	id := parse.NewRulesEngineID(subscriptionId, resourceGroup, frontDoorName, rulesEngineName).ID()

	frontdoorRulesEngineProperties := frontdoor.RulesEngineProperties{
		Rules: expandFrontDoorRulesEngineRules(rules),
	}

	frontdoorRulesEngine := frontdoor.RulesEngine{
		Name:                  utils.String(rulesEngineName),
		RulesEngineProperties: &frontdoorRulesEngineProperties,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, frontDoorName, rulesEngineName, frontdoorRulesEngine)
	if err != nil {
		return fmt.Errorf("creating Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
	}

	d.SetId(id)
	return resourceFrontDoorRulesEngineRead(d, meta)
}

func expandFrontDoorRulesEngineAction(input []interface{}) *frontdoor.RulesEngineAction {
	if len(input) == 0 {
		return nil
	}

	ruleAction := input[0].(map[string]interface{})

	requestHeaderActions := ruleAction["request_header"].([]interface{})
	responseHeaderActions := ruleAction["response_header"].([]interface{})

	frontdoorRulesEngineRuleAction := frontdoor.RulesEngineAction{
		RequestHeaderActions:  expandHeaderAction(requestHeaderActions),
		ResponseHeaderActions: expandHeaderAction(responseHeaderActions),
	}

	return &frontdoorRulesEngineRuleAction
}

func expandHeaderAction(input []interface{}) *[]frontdoor.HeaderAction {
	if len(input) == 0 {
		return nil
	}
	output := make([]frontdoor.HeaderAction, 0)

	for _, a := range input {
		action := a.(map[string]interface{})

		headerName := action["header_name"].(string)
		value := action["value"].(string)
		headerActionType := action["header_action_type"].(string)

		frontdoorRulesEngineRuleHeaderAction := frontdoor.HeaderAction{
			HeaderName:       utils.String(headerName),
			Value:            utils.String(value),
			HeaderActionType: frontdoor.HeaderActionType(headerActionType),
		}

		output = append(output, frontdoorRulesEngineRuleHeaderAction)
	}

	return &output
}

func expandFrontDoorRulesEngineRules(input []interface{}) *[]frontdoor.RulesEngineRule {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoor.RulesEngineRule, 0)

	for _, r := range input {
		rule := r.(map[string]interface{})

		ruleName := rule["name"].(string)
		priority := int32(rule["priority"].(int))
		actions := rule["action"].([]interface{})
		matchConditions := rule["match_condition"].([]interface{})

		frontdoorRulesEngineRule := frontdoor.RulesEngineRule{
			Name:            utils.String(ruleName),
			Priority:        utils.Int32(priority),
			Action:          expandFrontDoorRulesEngineAction(actions),
			MatchConditions: expandFrontDoorRulesEngineMatchCondition(matchConditions),
		}

		output = append(output, frontdoorRulesEngineRule)
	}
	return &output
}

func expandFrontDoorRulesEngineMatchCondition(input []interface{}) *[]frontdoor.RulesEngineMatchCondition {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoor.RulesEngineMatchCondition, 0)

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

		matchCondition := frontdoor.RulesEngineMatchCondition{
			RulesEngineMatchVariable: frontdoor.RulesEngineMatchVariable(matchVariable),
			Selector:                 utils.String(selector),
			RulesEngineOperator:      frontdoor.RulesEngineOperator(operator),
			NegateCondition:          &negateCondition,
			RulesEngineMatchValue:    &matchValueArray,
			Transforms:               expandFrontDoorRulesEngineMatchConditionTransform(transform),
		}
		output = append(output, matchCondition)
	}
	return &output
}

func expandFrontDoorRulesEngineMatchConditionTransform(input []interface{}) *[]frontdoor.Transform {
	if len(input) == 0 {
		return &[]frontdoor.Transform{}
	}

	output := make([]frontdoor.Transform, 0)

	for _, t := range input {
		result := frontdoor.Transform(t.(string))

		output = append(output, result)
	}
	return &output
}

func resourceFrontDoorRulesEngineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsRulesEnginesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	frontDoorName := d.Get("frontdoor_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	rulesEngineName := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, frontDoorName, rulesEngineName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door Rules Engine %q does not exist - removing from state", rulesEngineName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
	}
	return nil
}

func resourceFrontDoorRulesEngineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsRulesEnginesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	frontDoorName := d.Get("frontdoor_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	rulesEngineName := d.Get("name").(string)

	future, err := client.Delete(ctx, resourceGroup, frontDoorName, rulesEngineName)
	if err != nil {
		return fmt.Errorf("deleting Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
		}
	}
	return nil
}
