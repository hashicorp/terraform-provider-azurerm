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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	afd "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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
				ValidateFunc: validate.NoEmptyStrings,
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

			"redirect_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"custom_block_response_status_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(100, 530),
			},

			"custom_block_response_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: afd.ValidateCustomBlockResponseBody,
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
						"match_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 100,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_variable": {
										Type:     schema.TypeString,
										Required: true,
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
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.NoEmptyStrings,
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
									"negation_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
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
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
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
	enabled := frontdoor.PolicyEnabledStateDisabled
	if d.Get("enabled").(bool) {
		enabled = frontdoor.PolicyEnabledStateEnabled
	}
	mode := d.Get("mode").(string)
	redirectUrl := d.Get("redirect_url").(string)
	customBlockResponseStatusCode := d.Get("custom_block_response_status_code").(int)
	customBlockResponseBody := d.Get("custom_block_response_body").(string)
	customRules := d.Get("custom_rule").([]interface{})
	managedRules := d.Get("managed_rule").([]interface{})
	tags := d.Get("tags").(map[string]interface{})

	frontdoorWebApplicationFirewallPolicy := frontdoor.WebApplicationFirewallPolicy{
		Name:     utils.String(name),
		Location: utils.String(location),
		WebApplicationFirewallPolicyProperties: &frontdoor.WebApplicationFirewallPolicyProperties{
			PolicySettings: &frontdoor.PolicySettings{
				EnabledState:                  enabled,
				Mode:                          frontdoor.PolicyMode(mode),
				RedirectURL:                   utils.String(redirectUrl),
				CustomBlockResponseStatusCode: utils.Int32(int32(customBlockResponseStatusCode)),
				CustomBlockResponseBody:       utils.String(customBlockResponseBody),
			},
			CustomRules:  expandArmFrontDoorFirewallCustomRules(customRules),
			ManagedRules: expandArmFrontDoorFirewallManagedRules(managedRules),
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
	client := meta.(*ArmClient).frontdoor.FrontDoorsPolicyClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["frontdoorwebapplicationfirewallpolicies"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door Firewall Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Front Door Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if properties := resp.WebApplicationFirewallPolicyProperties; properties != nil {
		if properties.PolicySettings != nil {

			if err := d.Set("enabled", properties.PolicySettings.EnabledState == frontdoor.PolicyEnabledStateEnabled); err != nil {
				return fmt.Errorf("Error setting `enabled`: %+v", err)
			}

			if err := d.Set("mode", string(properties.PolicySettings.Mode)); err != nil {
				return fmt.Errorf("Error setting `mode`: %+v", err)
			}

			if err := d.Set("redirect_url", properties.PolicySettings.RedirectURL); err != nil {
				return fmt.Errorf("Error setting `redirect_url`: %+v", err)
			}

			if err := d.Set("custom_block_response_status_code", properties.PolicySettings.CustomBlockResponseStatusCode); err != nil {
				return fmt.Errorf("Error setting `custom_block_response_status_code`: %+v", err)
			}

			if err := d.Set("custom_block_response_body", properties.PolicySettings.CustomBlockResponseBody); err != nil {
				return fmt.Errorf("Error setting `custom_block_response_body`: %+v", err)
			}
		}

		if err := d.Set("custom_rule", flattenArmFrontDoorFirewallCustomRules(properties.CustomRules)); err != nil {
			return fmt.Errorf("Error flattening `custom_rule`: %+v", err)
		}

		if err := d.Set("managed_rule", flattenArmFrontDoorFirewallManagedRules(properties.ManagedRules)); err != nil {
			return fmt.Errorf("Error flattening `managed_rule`: %+v", err)
		}

		if err := d.Set("frontend_endpoint_ids", afd.FlattenFrontendEndpointLinkSlice(properties.FrontendEndpointLinks)); err != nil {
			return fmt.Errorf("Error flattening `frontend_endpoint_ids`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

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
		custom := cr.(map[string]interface{})

		enabled := frontdoor.CustomRuleEnabledStateDisabled
		if custom["enabled"].(bool) {
			enabled = frontdoor.CustomRuleEnabledStateEnabled
		}
		name := custom["name"].(string)
		priority := int32(custom["priority"].(int))
		ruleType := custom["rule_type"].(string)
		rateLimitDurationInMinutes := int32(custom["rate_limit_duration_in_minutes"].(int))
		rateLimitThreshold := int32(custom["rate_limit_threshold"].(int))
		matchConditions := custom["match_condition"].([]interface{})
		action := custom["action"].(string)

		customRule := frontdoor.CustomRule{
			Name:                       utils.String(name),
			Priority:                   utils.Int32(priority),
			EnabledState:               enabled,
			RuleType:                   frontdoor.RuleType(ruleType),
			RateLimitDurationInMinutes: utils.Int32(rateLimitDurationInMinutes),
			RateLimitThreshold:         utils.Int32(rateLimitThreshold),
			MatchConditions:            expandArmFrontDoorFirewallMatchConditions(matchConditions),
			Action:                     frontdoor.ActionType(action),
		}
		output = append(output, customRule)
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

	result := make([]frontdoor.MatchCondition, 0)

	for _, v := range input {
		match := v.(map[string]interface{})

		matchVariable := match["match_variable"].(string)
		selector := match["selector"].(string)
		operator := match["operator"].(string)
		negateCondition := match["negation_condition"].(bool)
		matchValues := match["match_value"].([]interface{})
		transforms := match["transforms"].([]interface{})

		matchCondition := frontdoor.MatchCondition{
			Operator:        frontdoor.Operator(operator),
			NegateCondition: &negateCondition,
			MatchValue:      utils.ExpandStringSlice(matchValues),
			Transforms:      expandArmFrontDoorFirewallTransforms(transforms),
		}

		if matchVariable != "" {
			matchCondition.MatchVariable = frontdoor.MatchVariable(matchVariable)
		}
		if selector != "" {
			matchCondition.Selector = utils.String(selector)
		}

		result = append(result, matchCondition)
	}

	return &result
}

func expandArmFrontDoorFirewallTransforms(input []interface{}) *[]frontdoor.TransformType {
	if len(input) == 0 {
		return nil
	}

	result := make([]frontdoor.TransformType, 0)

	for _, v := range input {
		result = append(result, frontdoor.TransformType(v.(string)))
	}

	return &result
}

func expandArmFrontDoorFirewallManagedRules(input []interface{}) *frontdoor.ManagedRuleSetList {
	if len(input) == 0 {
		return nil
	}

	managedRules := make([]frontdoor.ManagedRuleSet, 0)

	for _, mr := range input {
		managedRule := mr.(map[string]interface{})

		ruleType := managedRule["type"].(string)
		version := managedRule["version"].(string)
		overrides := managedRule["override"].([]interface{})

		managedRuleSet := frontdoor.ManagedRuleSet{
			RuleSetType:    utils.String(ruleType),
			RuleSetVersion: utils.String(version),
		}

		if ruleGroupOverrides := expandArmFrontDoorFirewallManagedRuleGroupOverride(overrides); ruleGroupOverrides != nil {
			managedRuleSet.RuleGroupOverrides = ruleGroupOverrides
		}

		managedRules = append(managedRules, managedRuleSet)
	}

	result := frontdoor.ManagedRuleSetList{
		ManagedRuleSets: &managedRules,
	}

	return &result
}

func expandArmFrontDoorFirewallManagedRuleGroupOverride(input []interface{}) *[]frontdoor.ManagedRuleGroupOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleGroupOverrides := make([]frontdoor.ManagedRuleGroupOverride, 0)

	for _, v := range input {
		override := v.(map[string]interface{})

		ruleGroupName := override["rule_group_name"].(string)
		rules := override["rule"].([]interface{})

		managedRuleGroupOverride := frontdoor.ManagedRuleGroupOverride{
			RuleGroupName: utils.String(ruleGroupName),
		}

		if managedRuleOverride := expandArmFrontDoorFirewallRuleOverride(rules); managedRuleOverride != nil {
			managedRuleGroupOverride.Rules = managedRuleOverride
		}

		managedRuleGroupOverrides = append(managedRuleGroupOverrides, managedRuleGroupOverride)
	}

	return &managedRuleGroupOverrides
}

func expandArmFrontDoorFirewallRuleOverride(input []interface{}) *[]frontdoor.ManagedRuleOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleOverrides := make([]frontdoor.ManagedRuleOverride, 0)

	for _, v := range input {
		rule := v.(map[string]interface{})

		enabled := frontdoor.ManagedRuleEnabledStateDisabled
		if rule["enabled"].(bool) {
			enabled = frontdoor.ManagedRuleEnabledStateEnabled
		}
		ruleId := rule["rule_id"].(string)
		action := rule["action"].(string)

		managedRuleOverride := frontdoor.ManagedRuleOverride{
			RuleID:       utils.String(ruleId),
			EnabledState: enabled,
			Action:       frontdoor.ActionType(action),
		}

		managedRuleOverrides = append(managedRuleOverrides, managedRuleOverride)
	}

	return &managedRuleOverrides
}

func flattenArmFrontDoorFirewallCustomRules(input *frontdoor.CustomRuleList) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input.Rules {
		output := make(map[string]interface{})

		output["name"] = v.Name
		output["priority"] = int(*v.Priority)
		output["enabled"] = v.EnabledState == frontdoor.CustomRuleEnabledStateEnabled
		output["rule_type"] = string(v.RuleType)
		output["rate_limit_duration_in_minutes"] = int(*v.RateLimitDurationInMinutes)
		output["rate_limit_threshold"] = int(*v.RateLimitThreshold)
		output["match_condition"] = flattenArmFrontDoorFirewallMatchConditions(v.MatchConditions)
		output["action"] = string(v.Action)

		results = append(results, output)
	}

	return results
}

func flattenArmFrontDoorFirewallMatchConditions(input *[]frontdoor.MatchCondition) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		output["match_variable"] = string(v.MatchVariable)
		if selector := v.Selector; selector != nil {
			output["selector"] = string(*v.Selector)
		}
		output["operator"] = string(v.Operator)
		if negateCondition := v.NegateCondition; negateCondition != nil {
			output["negation_condition"] = *v.NegateCondition
		}
		output["match_value"] = utils.FlattenStringSlice(v.MatchValue)
		output["transforms"] = afd.FlattenTransformSlice(v.Transforms)

		results = append(results, output)
	}

	return results
}

func flattenArmFrontDoorFirewallManagedRules(input *frontdoor.ManagedRuleSetList) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input.ManagedRuleSets {
		output := make(map[string]interface{})

		output["type"] = string(*v.RuleSetType)
		output["version"] = string(*v.RuleSetVersion)
		if overrides := v.RuleGroupOverrides; overrides != nil {
			output["override"] = flattenArmFrontDoorFirewallOverrides(overrides)
		}

		results = append(results, output)
	}

	return results
}

func flattenArmFrontDoorFirewallOverrides(input *[]frontdoor.ManagedRuleGroupOverride) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		output["rule_group_name"] = string(*v.RuleGroupName)
		if rules := v.Rules; rules != nil {
			output["rule"] = flattenArmFrontdoorFirewallRules(rules)
		}

		results = append(results, output)
	}

	return results
}

func flattenArmFrontdoorFirewallRules(input *[]frontdoor.ManagedRuleOverride) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		output["rule_id"] = string(*v.RuleID)
		output["enabled"] = v.EnabledState == frontdoor.ManagedRuleEnabledStateEnabled
		output["action"] = string(v.Action)

		results = append(results, output)
	}

	return results
}
