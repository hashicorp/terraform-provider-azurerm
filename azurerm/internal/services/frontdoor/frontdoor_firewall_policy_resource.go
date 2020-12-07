package frontdoor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceFrontDoorFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceFrontDoorFirewallPolicyCreateUpdate,
		Read:   resourceFrontDoorFirewallPolicyRead,
		Update: resourceFrontDoorFirewallPolicyCreateUpdate,
		Delete: resourceFrontDoorFirewallPolicyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.WebApplicationFirewallPolicyID(id)
			return err
		}),

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
				ValidateFunc: validate.FrontDoorWAFName,
			},

			"location": location.SchemaComputed(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(frontdoor.Detection),
					string(frontdoor.Prevention),
				}, false),
				Default: string(frontdoor.Prevention),
			},

			"redirect_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
			},

			"custom_block_response_status_code": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					200,
					403,
					405,
					406,
					429,
				}),
			},

			"custom_block_response_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.CustomBlockResponseBody,
			},

			"custom_rule": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},

						"type": {
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
											string(frontdoor.SocketAddr),
										}, false),
									},

									"match_values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 100,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},

									"operator": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.OperatorAny),
											string(frontdoor.OperatorBeginsWith),
											string(frontdoor.OperatorContains),
											string(frontdoor.OperatorEndsWith),
											string(frontdoor.OperatorEqual),
											string(frontdoor.OperatorGeoMatch),
											string(frontdoor.OperatorGreaterThan),
											string(frontdoor.OperatorGreaterThanOrEqual),
											string(frontdoor.OperatorIPMatch),
											string(frontdoor.OperatorLessThan),
											string(frontdoor.OperatorLessThanOrEqual),
											string(frontdoor.OperatorRegEx),
										}, false),
									},

									"selector": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"negation_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"version": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"exclusion": {
							Type:     schema.TypeList,
							MaxItems: 100,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_variable": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.QueryStringArgNames),
											string(frontdoor.RequestBodyPostArgNames),
											string(frontdoor.RequestCookieNames),
											string(frontdoor.RequestHeaderNames),
										}, false),
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.Contains),
											string(frontdoor.EndsWith),
											string(frontdoor.Equals),
											string(frontdoor.EqualsAny),
											string(frontdoor.StartsWith),
										}, false),
									},
									"selector": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"override": {
							Type:     schema.TypeList,
							MaxItems: 100,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_group_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"exclusion": {
										Type:     schema.TypeList,
										MaxItems: 100,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_variable": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(frontdoor.QueryStringArgNames),
														string(frontdoor.RequestBodyPostArgNames),
														string(frontdoor.RequestCookieNames),
														string(frontdoor.RequestHeaderNames),
													}, false),
												},
												"operator": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(frontdoor.Contains),
														string(frontdoor.EndsWith),
														string(frontdoor.Equals),
														string(frontdoor.EqualsAny),
														string(frontdoor.StartsWith),
													}, false),
												},
												"selector": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
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
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},

												"exclusion": {
													Type:     schema.TypeList,
													MaxItems: 100,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_variable": {
																Type:     schema.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(frontdoor.QueryStringArgNames),
																	string(frontdoor.RequestBodyPostArgNames),
																	string(frontdoor.RequestCookieNames),
																	string(frontdoor.RequestHeaderNames),
																}, false),
															},
															"operator": {
																Type:     schema.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(frontdoor.Contains),
																	string(frontdoor.EndsWith),
																	string(frontdoor.Equals),
																	string(frontdoor.EqualsAny),
																	string(frontdoor.StartsWith),
																}, false),
															},
															"selector": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
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

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.WebApplicationFirewallPolicyV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.WebApplicationFirewallPolicyV0ToV1,
				Version: 0,
			},
		},
	}
}

func resourceFrontDoorFirewallPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing args for Front Door Firewall Policy")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := parse.NewWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, name).ID("")

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Front Door Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_frontdoor_firewall_policy", id)
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

	t := d.Get("tags").(map[string]interface{})

	frontdoorWebApplicationFirewallPolicy := frontdoor.WebApplicationFirewallPolicy{
		Name:     utils.String(name),
		Location: utils.String(location),
		WebApplicationFirewallPolicyProperties: &frontdoor.WebApplicationFirewallPolicyProperties{
			PolicySettings: &frontdoor.PolicySettings{
				EnabledState: enabled,
				Mode:         frontdoor.PolicyMode(mode),
			},
			CustomRules:  expandFrontDoorFirewallCustomRules(customRules),
			ManagedRules: expandFrontDoorFirewallManagedRules(managedRules),
		},
		Tags: tags.Expand(t),
	}

	if redirectUrl != "" {
		frontdoorWebApplicationFirewallPolicy.WebApplicationFirewallPolicyProperties.PolicySettings.RedirectURL = utils.String(redirectUrl)
	}
	if customBlockResponseBody != "" {
		frontdoorWebApplicationFirewallPolicy.WebApplicationFirewallPolicyProperties.PolicySettings.CustomBlockResponseBody = utils.String(customBlockResponseBody)
	}
	if customBlockResponseStatusCode > 0 {
		frontdoorWebApplicationFirewallPolicy.WebApplicationFirewallPolicyProperties.PolicySettings.CustomBlockResponseStatusCode = utils.Int32(int32(customBlockResponseStatusCode))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, frontdoorWebApplicationFirewallPolicy)
	if err != nil {
		return fmt.Errorf("creating Front Door Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Front Door Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id)
	return resourceFrontDoorFirewallPolicyRead(d, meta)
}

func resourceFrontDoorFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door Firewall Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Front Door Firewall Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if properties := resp.WebApplicationFirewallPolicyProperties; properties != nil {
		if policy := properties.PolicySettings; policy != nil {
			d.Set("enabled", policy.EnabledState == frontdoor.PolicyEnabledStateEnabled)
			d.Set("mode", string(policy.Mode))
			d.Set("redirect_url", policy.RedirectURL)
			d.Set("custom_block_response_status_code", policy.CustomBlockResponseStatusCode)
			d.Set("custom_block_response_body", policy.CustomBlockResponseBody)
		}

		if err := d.Set("custom_rule", flattenFrontDoorFirewallCustomRules(properties.CustomRules)); err != nil {
			return fmt.Errorf("flattening `custom_rule`: %+v", err)
		}

		if err := d.Set("frontend_endpoint_ids", FlattenFrontendEndpointLinkSlice(properties.FrontendEndpointLinks)); err != nil {
			return fmt.Errorf("flattening `frontend_endpoint_ids`: %+v", err)
		}

		if err := d.Set("managed_rule", flattenFrontDoorFirewallManagedRules(properties.ManagedRules)); err != nil {
			return fmt.Errorf("flattening `managed_rule`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceFrontDoorFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Front Door Firewall %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Front Door Firewall %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandFrontDoorFirewallCustomRules(input []interface{}) *frontdoor.CustomRuleList {
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
		ruleType := custom["type"].(string)
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
			MatchConditions:            expandFrontDoorFirewallMatchConditions(matchConditions),
			Action:                     frontdoor.ActionType(action),
		}
		output = append(output, customRule)
	}

	return &frontdoor.CustomRuleList{
		Rules: &output,
	}
}

func expandFrontDoorFirewallMatchConditions(input []interface{}) *[]frontdoor.MatchCondition {
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
		matchValues := match["match_values"].([]interface{})
		transforms := match["transforms"].([]interface{})

		matchCondition := frontdoor.MatchCondition{
			Operator:        frontdoor.Operator(operator),
			NegateCondition: &negateCondition,
			MatchValue:      utils.ExpandStringSlice(matchValues),
			Transforms:      expandFrontDoorFirewallTransforms(transforms),
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

func expandFrontDoorFirewallTransforms(input []interface{}) *[]frontdoor.TransformType {
	if len(input) == 0 {
		return nil
	}

	result := make([]frontdoor.TransformType, 0)
	for _, v := range input {
		result = append(result, frontdoor.TransformType(v.(string)))
	}

	return &result
}

func expandFrontDoorFirewallManagedRules(input []interface{}) *frontdoor.ManagedRuleSetList {
	if len(input) == 0 {
		return nil
	}

	managedRules := make([]frontdoor.ManagedRuleSet, 0)

	for _, mr := range input {
		managedRule := mr.(map[string]interface{})

		ruleType := managedRule["type"].(string)
		version := managedRule["version"].(string)
		overrides := managedRule["override"].([]interface{})
		exclusions := managedRule["exclusion"].([]interface{})

		managedRuleSet := frontdoor.ManagedRuleSet{
			RuleSetType:    utils.String(ruleType),
			RuleSetVersion: utils.String(version),
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(exclusions); exclusions != nil {
			managedRuleSet.Exclusions = exclusions
		}

		if ruleGroupOverrides := expandFrontDoorFirewallManagedRuleGroupOverride(overrides); ruleGroupOverrides != nil {
			managedRuleSet.RuleGroupOverrides = ruleGroupOverrides
		}

		managedRules = append(managedRules, managedRuleSet)
	}

	return &frontdoor.ManagedRuleSetList{
		ManagedRuleSets: &managedRules,
	}
}

func expandFrontDoorFirewallManagedRuleGroupExclusion(input []interface{}) *[]frontdoor.ManagedRuleExclusion {
	if len(input) == 0 {
		return nil
	}

	managedRuleExclusions := make([]frontdoor.ManagedRuleExclusion, 0)
	for _, v := range input {
		exclusion := v.(map[string]interface{})

		matchVariable := exclusion["match_variable"].(string)
		operator := exclusion["operator"].(string)
		selector := exclusion["selector"].(string)

		managedRuleExclusion := frontdoor.ManagedRuleExclusion{
			MatchVariable:         frontdoor.ManagedRuleExclusionMatchVariable(matchVariable),
			SelectorMatchOperator: frontdoor.ManagedRuleExclusionSelectorMatchOperator(operator),
			Selector:              utils.String(selector),
		}

		managedRuleExclusions = append(managedRuleExclusions, managedRuleExclusion)
	}

	return &managedRuleExclusions
}

func expandFrontDoorFirewallManagedRuleGroupOverride(input []interface{}) *[]frontdoor.ManagedRuleGroupOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleGroupOverrides := make([]frontdoor.ManagedRuleGroupOverride, 0)
	for _, v := range input {
		override := v.(map[string]interface{})

		ruleGroupName := override["rule_group_name"].(string)
		rules := override["rule"].([]interface{})
		exclusions := override["exclusion"].([]interface{})

		managedRuleGroupOverride := frontdoor.ManagedRuleGroupOverride{
			RuleGroupName: utils.String(ruleGroupName),
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(exclusions); exclusions != nil {
			managedRuleGroupOverride.Exclusions = exclusions
		}

		if managedRuleOverride := expandFrontDoorFirewallRuleOverride(rules); managedRuleOverride != nil {
			managedRuleGroupOverride.Rules = managedRuleOverride
		}

		managedRuleGroupOverrides = append(managedRuleGroupOverrides, managedRuleGroupOverride)
	}

	return &managedRuleGroupOverrides
}

func expandFrontDoorFirewallRuleOverride(input []interface{}) *[]frontdoor.ManagedRuleOverride {
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
		exclusions := rule["exclusion"].([]interface{})

		managedRuleOverride := frontdoor.ManagedRuleOverride{
			RuleID:       utils.String(ruleId),
			EnabledState: enabled,
			Action:       frontdoor.ActionType(action),
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(exclusions); exclusions != nil {
			managedRuleOverride.Exclusions = exclusions
		}

		managedRuleOverrides = append(managedRuleOverrides, managedRuleOverride)
	}

	return &managedRuleOverrides
}

func flattenFrontDoorFirewallCustomRules(input *frontdoor.CustomRuleList) []interface{} {
	if input == nil || input.Rules == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, r := range *input.Rules {
		output := make(map[string]interface{})

		output["name"] = r.Name
		output["type"] = string(r.RuleType)
		output["action"] = string(r.Action)
		output["enabled"] = r.EnabledState == frontdoor.CustomRuleEnabledStateEnabled
		output["match_condition"] = flattenFrontDoorFirewallMatchConditions(r.MatchConditions)

		if v := r.Priority; v != nil {
			output["priority"] = int(*v)
		}

		if v := r.RateLimitDurationInMinutes; v != nil {
			output["rate_limit_duration_in_minutes"] = int(*v)
		}

		if v := r.RateLimitThreshold; v != nil {
			output["rate_limit_threshold"] = int(*v)
		}

		results = append(results, output)
	}

	return results
}

func flattenFrontDoorFirewallMatchConditions(condition *[]frontdoor.MatchCondition) []interface{} {
	if condition == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, c := range *condition {
		output := make(map[string]interface{})

		output["match_variable"] = string(c.MatchVariable)
		output["operator"] = string(c.Operator)
		output["match_values"] = utils.FlattenStringSlice(c.MatchValue)
		output["transforms"] = FlattenTransformSlice(c.Transforms)

		if v := c.Selector; v != nil {
			output["selector"] = *v
		}

		if v := c.NegateCondition; v != nil {
			output["negation_condition"] = *v
		}

		results = append(results, output)
	}

	return results
}

func flattenFrontDoorFirewallManagedRules(input *frontdoor.ManagedRuleSetList) []interface{} {
	if input == nil || input.ManagedRuleSets == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, r := range *input.ManagedRuleSets {
		output := make(map[string]interface{})

		if v := r.RuleSetType; v != nil {
			output["type"] = *v
		}

		if v := r.RuleSetVersion; v != nil {
			output["version"] = *v
		}

		if v := r.RuleGroupOverrides; v != nil {
			output["override"] = flattenFrontDoorFirewallOverrides(v)
		}

		if v := r.Exclusions; v != nil {
			output["exclusion"] = flattenFrontDoorFirewallExclusions(v)
		}

		results = append(results, output)
	}

	return results
}

func flattenFrontDoorFirewallExclusions(managedRuleExclusion *[]frontdoor.ManagedRuleExclusion) []interface{} {
	if managedRuleExclusion == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, o := range *managedRuleExclusion {
		output := make(map[string]interface{})

		output["match_variable"] = o.MatchVariable
		output["operator"] = o.SelectorMatchOperator
		output["selector"] = o.Selector

		results = append(results, output)
	}

	return results
}

func flattenFrontDoorFirewallOverrides(groupOverride *[]frontdoor.ManagedRuleGroupOverride) []interface{} {
	if groupOverride == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, o := range *groupOverride {
		output := make(map[string]interface{})

		if v := o.RuleGroupName; v != nil {
			output["rule_group_name"] = *v
		}

		if v := o.Exclusions; v != nil {
			output["exclusion"] = flattenFrontDoorFirewallExclusions(v)
		}

		if rules := o.Rules; rules != nil {
			output["rule"] = flattenArmFrontdoorFirewallRules(rules)
		}

		results = append(results, output)
	}

	return results
}

func flattenArmFrontdoorFirewallRules(override *[]frontdoor.ManagedRuleOverride) []interface{} {
	if override == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, o := range *override {
		output := make(map[string]interface{})

		output["enabled"] = o.EnabledState == frontdoor.ManagedRuleEnabledStateEnabled
		output["action"] = string(o.Action)

		if v := o.RuleID; v != nil {
			output["rule_id"] = *v
		}

		if v := o.Exclusions; v != nil {
			output["exclusion"] = flattenFrontDoorFirewallExclusions(v)
		}

		results = append(results, output)
	}

	return results
}
