package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/frontdoor/mgmt/2019-04-01/frontdoor"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	afd "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFrontDoorFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFrontDoorFirewallPolicyCreateUpdate,
		Read:   resourceArmFrontDoorFirewallPolicyRead,
		Update: resourceArmFrontDoorFirewallPolicyCreateUpdate,
		Delete: resourceArmFrontDoorFirewallPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoEmptyStrings,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(frontdoor.Detection),
					string(frontdoor.Prevention),
				}, false),
			},

			"custom_block_response_status_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.CustomBlockResponseBody,
			},

			"redirect_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"custom_block_response_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.CustomBlockResponseBody,
			},

			"custom_rule": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"rule_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.MatchRule),
								string(frontdoor.RateLimitRule),
							}, false),
						},
						"rate_limit_duration_in_minutes": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"rate_limit_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  10,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoor.Allow),
								string(frontdoor.Block),
								string(frontdoor.Log),
								string(frontdoor.Redirect),
							}, false),
						},
						"custom_block_response_body": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"match_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 100,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_variable": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"selector"},
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.Cookies),
											string(frontdoor.PostArgs),
											string(frontdoor.QueryString),
											string(frontdoor.RemoteAddr),
											string(frontdoor.RequestBody),
											string(frontdoor.RequestHeader),
											string(frontdoor.RequestMethod),
											string(frontdoor.RequestURI),
										}, false),
									},
									"selector": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"match_variable"},
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.Cookies),
											string(frontdoor.PostArgs),
											string(frontdoor.QueryString),
											string(frontdoor.RequestHeader),
										}, false),
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.Any),
											string(frontdoor.BeginsWith),
											string(frontdoor.Contains),
											string(frontdoor.EndsWith),
											string(frontdoor.Equal),
											string(frontdoor.GeoMatch),
											string(frontdoor.GreaterThan),
											string(frontdoor.GreaterThanOrEqual),
											string(frontdoor.IPMatch),
											string(frontdoor.LessThan),
											string(frontdoor.LessThanOrEqual),
											string(frontdoor.RegEx),
										}, false),
									},
									"condition": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Is",
											"Is Not",
											"Contains",
											"Not Contains",
										}, false),
										Default: "Is",
									},
									"match_value": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 100,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validate.NoEmptyStrings,
										},
									},
									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 5,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(frontdoor.Lowercase),
												string(frontdoor.RemoveNulls),
												string(frontdoor.Trim),
												string(frontdoor.Uppercase),
												string(frontdoor.URLDecode),
												string(frontdoor.URLEncode),
											}, false),
										},
									},
								},
							},
						},
					},
				},
			},

			"managed_rule": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"override": {
							Type:     schema.TypeList,
							MaxItems: 100,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_group_name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
									"rule": {
										Type:     schema.TypeList,
										MaxItems: 1000,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_id": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validate.NoEmptyStrings,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"action": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(frontdoor.Allow),
														string(frontdoor.Block),
														string(frontdoor.Log),
														string(frontdoor.Redirect),
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

			"frontend_endpoint_ids": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1000,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmFrontDoorFirewallPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontdoor.FrontDoorsPolicyClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing args for Front Door Firewall Policy")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Front Door Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_frontdoor_firewall_policy", *resp.ID)
		}
	}

	location := azure.NormalizeLocation("Global")
	enabled := d.Get("enabled").(bool)
	mode := d.Get("mode").(string)
	redirectUrl := d.Get("redirect_url ").(string)
	customBlockResponseStatusCode := d.Get("custom_block_response_status_code").(int32)
	customBlockResponseBody := d.Get("custom_block_response_body").(string)
	customRules := d.Get("custom_rule").([]interface{})
	managedRules := d.Get("managed_rule").([]interface{})
	frontendEndpoints := d.Get("frontend_endpoint_ids").([]interface{})
	tags := d.Get("tags").(map[string]interface{})

	if err := afd.VerifyCustomRules(customRules); err != nil {
		return fmt.Errorf(`Error validating "custom_rule" for Front Door Firewall Policy %q (Resource Group %q): %+v`, name, resourceGroup, err)
	}

	frontdoorWebApplicationFirewallPolicy := frontdoor.WebApplicationFirewallPolicy{
		Name:     utils.String(name),
		Location: utils.String(location),
		WebApplicationFirewallPolicyProperties: &frontdoor.WebApplicationFirewallPolicyProperties{
			PolicySettings: &frontdoor.PolicySettings{
				EnabledState:                  helper.ConvertToPolicyEnabledStateFromBool(enabled),
				Mode:                          helper.ConvertToPolicyModeFromString(mode),
				RedirectURL:                   utils.String(redirectUrl),
				CustomBlockResponseStatusCode: &customBlockResponseStatusCode,
				CustomBlockResponseBody:       utils.String(customBlockResponseBody),
			},
			customRules:           expandArmFrontDoorFirewallCustomRules(customRules),
			ManagedRules:          expandArmFrontDoorFirewallManagedRules(managedRules),
			FrontendEndpointLinks: expandArmFrontDoorFirewallManagedRules(frontendEndpoints),
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, frontdoorWebApplicationFirewallPolicy)
	if err != nil {
		return fmt.Errorf("Error creating Front Door Firewall policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Front Door Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Front Door Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Front Door Firewall %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmFrontDoorFirewallPolicyRead(d, meta)
}

func resourceArmFrontDoorFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceArmFrontDoorFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontdoor.FrontDoorsPolicyClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["FrontDoorWebApplicationFirewallPolicies"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Front Door Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Front Door Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmFrontDoorFirewallCustomRules(input []interface{}) *frontdoor.CustomRuleList {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoor.CustomRule, 0)

	for _, cr := range input {
		customRule := cr.(map[string]interface{})

		name := customRule["name"].(string)
		priority := int32(customRule["priority"].(int))
		enabled := customRule["enabled"].(bool)
		ruleType := customRule["rule_type"].(string)
		rateLimitDurationInMinutes := int32(customRule["rate_limit_duration_in_minutes"].(int))
		rateLimitThreshold := int32(customRule["rate_limit_threshold"].(int))
		matchConditions := expandArmFrontDoorFirewallMatchConditions(customRule["match_condition"].([]interface{}))
		action := customRule["action_type"].(string)

		afpCustomRule = frontdoor.CustomRule{
			Name:                       utils.String(name),
			Priority:                   utils.Int32(priority),
			EnabledState:               afd.ConvertBoolToCustomRuleEnabledState(enabled),
			RuleType:                   afd.ConvertStringToRuleType(ruleType),
			RateLimitDurationInMinutes: utils.Int32(rateLimitDurationInMinutes),
			RateLimitThreshold:         utils.Int32(rateLimitThreshold),
			MatchConditions:            expandArmFrontDoorFirewallMatchConditions(matchConditions),
			Action:                     afd.ConvertStringToActionType(action),
		}
		output = append(output, afpCustomRule)
	}

	result := frontdoor.CustomRuleList{
		Rules: &output,
	}

	return &result
}

func expandArmFrontDoorFirewallMatchConditions(input []interface{}) *[]frontdoor.MatchCondition {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoor.MatchCondition, 0)

	for _, mc := range input {
		matchCondition := mc.(map[string]interface{})
		matchVariable := matchCondition["match_variable"].(string)
		selector := matchCondition["selector"].(string)
		operator := matchCondition["operator"].(string)
		negateCondition := afd.ConvertConditionToBool(matchCondition["condition"])

		mv := matchCondition["match_value"].([]interface{})
		matchValues := make([]string, 0)
		for _, v := range mv {
			matchValues = append(matchValues, v.(string))
		}
		matchValue := &matchValues

		ts := matchCondition["transforms"].([]interface{})
		transform := make([]string, 0)
		for _, t := range ts {
			transform = append(transform, t.(string))
		}
		transforms := &transform

		fdpMatchCondition := &frontdoor.MatchCondition{
			Operator:        operator,
			NegateCondition: &negateCondition,
			MatchValue:      &matchValue,
			Transforms:      &transforms,
		}

		if matchvariable != "" {
			fdpMatchCondition.MatchVariable = matchVariable
		} else {
			fdpMatchCondition.Selector = utils.String(selector)
		}

		output = append(output, fdpMatchCondition)
	}

	return &output
}

func expandArmFrontDoorFirewallManagedRules(input []interface{}) *frontdoor.ManagedRuleSetList {

}
