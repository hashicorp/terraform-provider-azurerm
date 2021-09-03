package frontdoor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-05-01/frontdoor"
	"github.com/hashicorp/go-azure-helpers/response"
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

						"rule_action": {
							Type:     pluginsdk.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"request_header_actions": {
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
									"response_header_actions": {
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
									//"route_configuration_override": {
									//	Type:         pluginsdk.TypeString,
									//	Optional:     true,
									//	ValidateFunc: validation.StringIsNotEmpty,
									//},
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
		//Type
		//ID
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

	requestHeaderActions := ruleAction["request_header_actions"].([]interface{})
	responseHeaderActions := ruleAction["response_header_actions"].([]interface{})
	//routeConfigurationOverride := ruleAction["route_configuration_override"].([]interface{})

	frontdoorRulesEngineRuleAction := frontdoor.RulesEngineAction{
		RequestHeaderActions:  expandHeaderAction(requestHeaderActions),
		ResponseHeaderActions: expandHeaderAction(responseHeaderActions),
		//RouteConfigurationOverride: expandRouteConfigOverride(routeConfigurationOverride),
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
			HeaderName: utils.String(headerName),
			Value:      utils.String(value),
		}

		if headerActionType == "Append" {
			frontdoorRulesEngineRuleHeaderAction.HeaderActionType = frontdoor.Append
		}
		if headerActionType == "Delete" {
			frontdoorRulesEngineRuleHeaderAction.HeaderActionType = frontdoor.Delete
		}
		if headerActionType == "Overwrite" {
			frontdoorRulesEngineRuleHeaderAction.HeaderActionType = frontdoor.Overwrite
		}

		output = append(output, frontdoorRulesEngineRuleHeaderAction)
	}

	return &output
}

//func expandRouteConfigOverride(input []interface{}) frontdoor.BasicRouteConfiguration {
//	if len(input) == 0 {
//		return nil
//	}
//	return nil
//}

func expandFrontDoorRulesEngineRules(input []interface{}) *[]frontdoor.RulesEngineRule {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoor.RulesEngineRule, 0)

	for _, r := range input {
		rule := r.(map[string]interface{})

		ruleName := rule["name"].(string)
		priority := int32(rule["priority"].(int))
		actions := rule["rule_action"].([]interface{})

		frontdoorRulesEngineRule := frontdoor.RulesEngineRule{
			Name:     utils.String(ruleName),
			Priority: utils.Int32(priority),
			Action:   expandFrontDoorRulesEngineAction(actions),
			//MatchConditions:
			//MatchProcessingBehavior:
		}

		output = append(output, frontdoorRulesEngineRule)
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
