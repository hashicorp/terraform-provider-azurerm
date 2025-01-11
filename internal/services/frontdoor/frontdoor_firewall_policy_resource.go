// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontDoorFirewallPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorFirewallPolicyCreateUpdate,
		Read:   resourceFrontDoorFirewallPolicyRead,
		Update: resourceFrontDoorFirewallPolicyCreateUpdate,
		Delete: resourceFrontDoorFirewallPolicyDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WebApplicationFirewallPolicyV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WebApplicationFirewallPolicyIDInsensitively(id)
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
				ValidateFunc: validate.FrontDoorWAFName,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(webapplicationfirewallpolicies.PolicyModeDetection),
					string(webapplicationfirewallpolicies.PolicyModePrevention),
				}, false),
				Default: string(webapplicationfirewallpolicies.PolicyModePrevention),
			},

			"redirect_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
			},

			"custom_block_response_status_code": {
				Type:     pluginsdk.TypeInt,
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.CustomBlockResponseBody,
			},

			"custom_rule": {
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

						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"priority": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  1,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(webapplicationfirewallpolicies.RuleTypeMatchRule),
								string(webapplicationfirewallpolicies.RuleTypeRateLimitRule),
							}, false),
						},

						"rate_limit_duration_in_minutes": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  1,
						},

						"rate_limit_threshold": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  10,
						},

						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(webapplicationfirewallpolicies.ActionTypeAllow),
								string(webapplicationfirewallpolicies.ActionTypeBlock),
								string(webapplicationfirewallpolicies.ActionTypeLog),
								string(webapplicationfirewallpolicies.ActionTypeRedirect),
							}, false),
						},

						"match_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 10,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"match_variable": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.MatchVariableCookies),
											string(webapplicationfirewallpolicies.MatchVariablePostArgs),
											string(webapplicationfirewallpolicies.MatchVariableQueryString),
											string(webapplicationfirewallpolicies.MatchVariableRemoteAddr),
											string(webapplicationfirewallpolicies.MatchVariableRequestBody),
											string(webapplicationfirewallpolicies.MatchVariableRequestHeader),
											string(webapplicationfirewallpolicies.MatchVariableRequestMethod),
											string(webapplicationfirewallpolicies.MatchVariableRequestUri),
											string(webapplicationfirewallpolicies.MatchVariableSocketAddr),
										}, false),
									},

									"match_values": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MaxItems: 600,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringLenBetween(1, 256),
										},
									},

									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.OperatorAny),
											string(webapplicationfirewallpolicies.OperatorBeginsWith),
											string(webapplicationfirewallpolicies.OperatorContains),
											string(webapplicationfirewallpolicies.OperatorEndsWith),
											string(webapplicationfirewallpolicies.OperatorEqual),
											string(webapplicationfirewallpolicies.OperatorGeoMatch),
											string(webapplicationfirewallpolicies.OperatorGreaterThan),
											string(webapplicationfirewallpolicies.OperatorGreaterThanOrEqual),
											string(webapplicationfirewallpolicies.OperatorIPMatch),
											string(webapplicationfirewallpolicies.OperatorLessThan),
											string(webapplicationfirewallpolicies.OperatorLessThanOrEqual),
											string(webapplicationfirewallpolicies.OperatorRegEx),
										}, false),
									},

									"selector": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"negation_condition": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"transforms": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 5,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(webapplicationfirewallpolicies.TransformTypeLowercase),
												string(webapplicationfirewallpolicies.TransformTypeRemoveNulls),
												string(webapplicationfirewallpolicies.TransformTypeTrim),
												string(webapplicationfirewallpolicies.TransformTypeUppercase),
												string(webapplicationfirewallpolicies.TransformTypeURLDecode),
												string(webapplicationfirewallpolicies.TransformTypeURLEncode),
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
				Type:     pluginsdk.TypeList,
				MaxItems: 100,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"exclusion": {
							Type:     pluginsdk.TypeList,
							MaxItems: 100,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"match_variable": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableQueryStringArgNames),
											string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
											string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestCookieNames),
											string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestHeaderNames),
										}, false),
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorContains),
											string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
											string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEquals),
											string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
											string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
										}, false),
									},
									"selector": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"override": {
							Type:     pluginsdk.TypeList,
							MaxItems: 100,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"rule_group_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"exclusion": {
										Type:     pluginsdk.TypeList,
										MaxItems: 100,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"match_variable": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableQueryStringArgNames),
														string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
														string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestCookieNames),
														string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestHeaderNames),
													}, false),
												},
												"operator": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorContains),
														string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
														string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEquals),
														string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
														string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
													}, false),
												},
												"selector": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},

									"rule": {
										Type:     pluginsdk.TypeList,
										MaxItems: 1000,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"rule_id": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"enabled": {
													Type:     pluginsdk.TypeBool,
													Optional: true,
													Default:  false,
												},

												"exclusion": {
													Type:     pluginsdk.TypeList,
													MaxItems: 100,
													Optional: true,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"match_variable": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableQueryStringArgNames),
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestCookieNames),
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariableRequestHeaderNames),
																}, false),
															},
															"operator": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorContains),
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEquals),
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
																	string(webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
																}, false),
															},
															"selector": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},

												"action": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(webapplicationfirewallpolicies.ActionTypeAllow),
														string(webapplicationfirewallpolicies.ActionTypeBlock),
														string(webapplicationfirewallpolicies.ActionTypeLog),
														string(webapplicationfirewallpolicies.ActionTypeRedirect),
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceFrontDoorFirewallPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing args for Front Door Firewall Policy")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := webapplicationfirewallpolicies.NewFrontDoorWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.PoliciesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing Front Door Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_frontdoor_firewall_policy", id.ID())
		}
	}

	location := azure.NormalizeLocation("Global")
	enabled := webapplicationfirewallpolicies.PolicyEnabledStateDisabled
	if d.Get("enabled").(bool) {
		enabled = webapplicationfirewallpolicies.PolicyEnabledStateEnabled
	}
	mode := webapplicationfirewallpolicies.PolicyMode(d.Get("mode").(string))
	redirectUrl := d.Get("redirect_url").(string)
	customBlockResponseStatusCode := d.Get("custom_block_response_status_code").(int)
	customBlockResponseBody := d.Get("custom_block_response_body").(string)
	customRules := d.Get("custom_rule").([]interface{})
	managedRules := d.Get("managed_rule").([]interface{})

	t := d.Get("tags").(map[string]interface{})

	frontdoorWebApplicationFirewallPolicy := webapplicationfirewallpolicies.WebApplicationFirewallPolicy{
		Name:     utils.String(name),
		Location: utils.String(location),
		Properties: &webapplicationfirewallpolicies.WebApplicationFirewallPolicyProperties{
			PolicySettings: &webapplicationfirewallpolicies.PolicySettings{
				EnabledState: &enabled,
				Mode:         &mode,
			},
			CustomRules:  expandFrontDoorFirewallCustomRules(customRules),
			ManagedRules: expandFrontDoorFirewallManagedRules(managedRules),
		},
		Tags: tags.Expand(t),
	}

	if redirectUrl != "" {
		frontdoorWebApplicationFirewallPolicy.Properties.PolicySettings.RedirectURL = utils.String(redirectUrl)
	}
	if customBlockResponseBody != "" {
		frontdoorWebApplicationFirewallPolicy.Properties.PolicySettings.CustomBlockResponseBody = utils.String(customBlockResponseBody)
	}
	if customBlockResponseStatusCode > 0 {
		frontdoorWebApplicationFirewallPolicy.Properties.PolicySettings.CustomBlockResponseStatusCode = utils.Int64(int64(customBlockResponseStatusCode))
	}

	if err := client.PoliciesCreateOrUpdateThenPoll(ctx, id, frontdoorWebApplicationFirewallPolicy); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontDoorFirewallPolicyRead(d, meta)
}

func resourceFrontDoorFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.PoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Front Door Firewall Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FrontDoorWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}
		if properties := model.Properties; properties != nil {
			if policy := properties.PolicySettings; policy != nil {
				if policy.EnabledState != nil {
					d.Set("enabled", *policy.EnabledState == webapplicationfirewallpolicies.PolicyEnabledStateEnabled)
				}
				if policy.Mode != nil {
					d.Set("mode", string(*policy.Mode))
				}
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

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceFrontDoorFirewallPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if err := client.PoliciesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandFrontDoorFirewallCustomRules(input []interface{}) *webapplicationfirewallpolicies.CustomRuleList {
	if len(input) == 0 {
		return nil
	}

	output := make([]webapplicationfirewallpolicies.CustomRule, 0)

	for _, cr := range input {
		custom := cr.(map[string]interface{})

		enabled := webapplicationfirewallpolicies.CustomRuleEnabledStateDisabled
		if custom["enabled"].(bool) {
			enabled = webapplicationfirewallpolicies.CustomRuleEnabledStateEnabled
		}

		name := custom["name"].(string)
		priority := int64(custom["priority"].(int))
		ruleType := custom["type"].(string)
		rateLimitDurationInMinutes := int64(custom["rate_limit_duration_in_minutes"].(int))
		rateLimitThreshold := int64(custom["rate_limit_threshold"].(int))
		matchConditions := expandFrontDoorFirewallMatchConditions(custom["match_condition"].([]interface{}))
		action := custom["action"].(string)

		customRule := webapplicationfirewallpolicies.CustomRule{
			Name:                       utils.String(name),
			Priority:                   priority,
			EnabledState:               &enabled,
			RuleType:                   webapplicationfirewallpolicies.RuleType(ruleType),
			RateLimitDurationInMinutes: utils.Int64(rateLimitDurationInMinutes),
			RateLimitThreshold:         utils.Int64(rateLimitThreshold),
			MatchConditions:            matchConditions,
			Action:                     webapplicationfirewallpolicies.ActionType(action),
		}
		output = append(output, customRule)
	}

	return &webapplicationfirewallpolicies.CustomRuleList{
		Rules: &output,
	}
}

func expandFrontDoorFirewallMatchConditions(input []interface{}) []webapplicationfirewallpolicies.MatchCondition {
	if len(input) == 0 {
		return nil
	}

	result := make([]webapplicationfirewallpolicies.MatchCondition, 0)

	for _, v := range input {
		match := v.(map[string]interface{})

		matchVariable := match["match_variable"].(string)
		selector := match["selector"].(string)
		operator := match["operator"].(string)
		negateCondition := match["negation_condition"].(bool)
		matchValues := match["match_values"].([]interface{})
		transforms := match["transforms"].([]interface{})

		matchCondition := webapplicationfirewallpolicies.MatchCondition{
			Operator:        webapplicationfirewallpolicies.Operator(operator),
			NegateCondition: &negateCondition,
			MatchValue:      *utils.ExpandStringSlice(matchValues),
			Transforms:      expandFrontDoorFirewallTransforms(transforms),
		}

		if matchVariable != "" {
			matchCondition.MatchVariable = webapplicationfirewallpolicies.MatchVariable(matchVariable)
		}
		if selector != "" {
			matchCondition.Selector = utils.String(selector)
		}

		result = append(result, matchCondition)
	}

	return result
}

func expandFrontDoorFirewallTransforms(input []interface{}) *[]webapplicationfirewallpolicies.TransformType {
	if len(input) == 0 {
		return nil
	}

	result := make([]webapplicationfirewallpolicies.TransformType, 0)
	for _, v := range input {
		result = append(result, webapplicationfirewallpolicies.TransformType(v.(string)))
	}

	return &result
}

func expandFrontDoorFirewallManagedRules(input []interface{}) *webapplicationfirewallpolicies.ManagedRuleSetList {
	if len(input) == 0 {
		return nil
	}

	managedRules := make([]webapplicationfirewallpolicies.ManagedRuleSet, 0)

	for _, mr := range input {
		managedRule := mr.(map[string]interface{})

		ruleType := managedRule["type"].(string)
		version := managedRule["version"].(string)
		overrides := managedRule["override"].([]interface{})
		exclusions := managedRule["exclusion"].([]interface{})

		managedRuleSet := webapplicationfirewallpolicies.ManagedRuleSet{
			RuleSetType:    ruleType,
			RuleSetVersion: version,
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(exclusions); exclusions != nil {
			managedRuleSet.Exclusions = exclusions
		}

		if ruleGroupOverrides := expandFrontDoorFirewallManagedRuleGroupOverride(overrides); ruleGroupOverrides != nil {
			managedRuleSet.RuleGroupOverrides = ruleGroupOverrides
		}

		managedRules = append(managedRules, managedRuleSet)
	}

	return &webapplicationfirewallpolicies.ManagedRuleSetList{
		ManagedRuleSets: &managedRules,
	}
}

func expandFrontDoorFirewallManagedRuleGroupExclusion(input []interface{}) *[]webapplicationfirewallpolicies.ManagedRuleExclusion {
	if len(input) == 0 {
		return nil
	}

	managedRuleExclusions := make([]webapplicationfirewallpolicies.ManagedRuleExclusion, 0)
	for _, v := range input {
		exclusion := v.(map[string]interface{})

		matchVariable := exclusion["match_variable"].(string)
		operator := exclusion["operator"].(string)
		selector := exclusion["selector"].(string)

		managedRuleExclusion := webapplicationfirewallpolicies.ManagedRuleExclusion{
			MatchVariable:         webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariable(matchVariable),
			SelectorMatchOperator: webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperator(operator),
			Selector:              selector,
		}

		managedRuleExclusions = append(managedRuleExclusions, managedRuleExclusion)
	}

	return &managedRuleExclusions
}

func expandFrontDoorFirewallManagedRuleGroupOverride(input []interface{}) *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleGroupOverrides := make([]webapplicationfirewallpolicies.ManagedRuleGroupOverride, 0)
	for _, v := range input {
		override := v.(map[string]interface{})

		ruleGroupName := override["rule_group_name"].(string)
		rules := override["rule"].([]interface{})
		exclusions := override["exclusion"].([]interface{})

		managedRuleGroupOverride := webapplicationfirewallpolicies.ManagedRuleGroupOverride{
			RuleGroupName: ruleGroupName,
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

func expandFrontDoorFirewallRuleOverride(input []interface{}) *[]webapplicationfirewallpolicies.ManagedRuleOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleOverrides := make([]webapplicationfirewallpolicies.ManagedRuleOverride, 0)
	for _, v := range input {
		rule := v.(map[string]interface{})

		enabled := webapplicationfirewallpolicies.ManagedRuleEnabledStateDisabled
		if rule["enabled"].(bool) {
			enabled = webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled
		}
		ruleId := rule["rule_id"].(string)
		action := webapplicationfirewallpolicies.ActionType(rule["action"].(string))
		exclusions := rule["exclusion"].([]interface{})

		managedRuleOverride := webapplicationfirewallpolicies.ManagedRuleOverride{
			RuleId:       ruleId,
			EnabledState: &enabled,
			Action:       &action,
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(exclusions); exclusions != nil {
			managedRuleOverride.Exclusions = exclusions
		}

		managedRuleOverrides = append(managedRuleOverrides, managedRuleOverride)
	}

	return &managedRuleOverrides
}

func flattenFrontDoorFirewallCustomRules(input *webapplicationfirewallpolicies.CustomRuleList) []interface{} {
	if input == nil || input.Rules == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, r := range *input.Rules {
		output := make(map[string]interface{})

		output["name"] = r.Name
		output["type"] = string(r.RuleType)
		output["action"] = string(r.Action)
		if r.EnabledState != nil {
			output["enabled"] = *r.EnabledState == webapplicationfirewallpolicies.CustomRuleEnabledStateEnabled
		}
		output["match_condition"] = flattenFrontDoorFirewallMatchConditions(r.MatchConditions)
		output["priority"] = int(r.Priority)

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

func flattenFrontDoorFirewallMatchConditions(condition []webapplicationfirewallpolicies.MatchCondition) []interface{} {
	if condition == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, c := range condition {
		output := make(map[string]interface{})

		output["match_variable"] = string(c.MatchVariable)
		output["operator"] = string(c.Operator)
		output["match_values"] = c.MatchValue
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

func flattenFrontDoorFirewallManagedRules(input *webapplicationfirewallpolicies.ManagedRuleSetList) []interface{} {
	if input == nil || input.ManagedRuleSets == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, r := range *input.ManagedRuleSets {
		output := make(map[string]interface{})

		output["type"] = r.RuleSetType

		output["version"] = r.RuleSetVersion

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

func flattenFrontDoorFirewallExclusions(managedRuleExclusion *[]webapplicationfirewallpolicies.ManagedRuleExclusion) []interface{} {
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

func flattenFrontDoorFirewallOverrides(groupOverride *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride) []interface{} {
	if groupOverride == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, o := range *groupOverride {
		output := make(map[string]interface{})

		output["rule_group_name"] = o.RuleGroupName

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

func flattenArmFrontdoorFirewallRules(override *[]webapplicationfirewallpolicies.ManagedRuleOverride) []interface{} {
	if override == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, o := range *override {
		output := make(map[string]interface{})

		if o.EnabledState != nil {
			output["enabled"] = *o.EnabledState == webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled
		}
		if o.Action != nil {
			output["action"] = string(*o.Action)
		}

		output["rule_id"] = o.RuleId

		if v := o.Exclusions; v != nil {
			output["exclusion"] = flattenFrontDoorFirewallExclusions(v)
		}

		results = append(results, output)
	}

	return results
}
