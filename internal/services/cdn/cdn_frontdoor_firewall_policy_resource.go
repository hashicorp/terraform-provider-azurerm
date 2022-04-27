package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	legacyfrontdoor "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/legacysdk/2020-11-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorFirewallPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorFirewallPolicyCreateUpdate,
		Read:   resourceCdnFrontdoorFirewallPolicyRead,
		Update: resourceCdnFrontdoorFirewallPolicyCreateUpdate,
		Delete: resourceCdnFrontdoorFirewallPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorPolicyIDInsensitively(id)
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
				ValidateFunc: ValidatedLegacyFrontdoorWAFName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorProfileID,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(legacyfrontdoor.SkuNameStandardAzureFrontDoor),
				ValidateFunc: validation.StringInSlice([]string{
					string(legacyfrontdoor.SkuNameStandardAzureFrontDoor),
					string(legacyfrontdoor.SkuNamePremiumAzureFrontDoor),
				}, false),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(legacyfrontdoor.PolicyModeDetection),
					string(legacyfrontdoor.PolicyModePrevention),
				}, false),
				Default: string(legacyfrontdoor.PolicyModePrevention),
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
				ValidateFunc: ValidateLegacyCustomBlockResponseBody,
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
								string(legacyfrontdoor.RuleTypeMatchRule),
								string(legacyfrontdoor.RuleTypeRateLimitRule),
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
								string(legacyfrontdoor.ActionTypeAllow),
								string(legacyfrontdoor.ActionTypeBlock),
								string(legacyfrontdoor.ActionTypeLog),
								string(legacyfrontdoor.ActionTypeRedirect),
							}, false),
						},

						"match_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 10,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									// TODO - rename to "variable" for consistency
									"match_variable": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(legacyfrontdoor.MatchVariableCookies),
											string(legacyfrontdoor.MatchVariablePostArgs),
											string(legacyfrontdoor.MatchVariableQueryString),
											string(legacyfrontdoor.MatchVariableRemoteAddr),
											string(legacyfrontdoor.MatchVariableRequestBody),
											string(legacyfrontdoor.MatchVariableRequestHeader),
											string(legacyfrontdoor.MatchVariableRequestMethod),
											string(legacyfrontdoor.MatchVariableRequestURI),
											string(legacyfrontdoor.MatchVariableSocketAddr),
										}, false),
									},

									// TODO - rename to "value" for consistency
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
											string(legacyfrontdoor.OperatorAny),
											string(legacyfrontdoor.OperatorBeginsWith),
											string(legacyfrontdoor.OperatorContains),
											string(legacyfrontdoor.OperatorEndsWith),
											string(legacyfrontdoor.OperatorEqual),
											string(legacyfrontdoor.OperatorGeoMatch),
											string(legacyfrontdoor.OperatorGreaterThan),
											string(legacyfrontdoor.OperatorGreaterThanOrEqual),
											string(legacyfrontdoor.OperatorIPMatch),
											string(legacyfrontdoor.OperatorLessThan),
											string(legacyfrontdoor.OperatorLessThanOrEqual),
											string(legacyfrontdoor.OperatorRegEx),
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
												string(legacyfrontdoor.TransformTypeLowercase),
												string(legacyfrontdoor.TransformTypeRemoveNulls),
												string(legacyfrontdoor.TransformTypeTrim),
												string(legacyfrontdoor.TransformTypeUppercase),
												string(legacyfrontdoor.TransformTypeURLDecode),
												string(legacyfrontdoor.TransformTypeURLEncode),
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
											string(legacyfrontdoor.ManagedRuleExclusionMatchVariableQueryStringArgNames),
											string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
											string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestCookieNames),
											string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestHeaderNames),
										}, false),
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorContains),
											string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
											string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEquals),
											string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
											string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
														string(legacyfrontdoor.ManagedRuleExclusionMatchVariableQueryStringArgNames),
														string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
														string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestCookieNames),
														string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestHeaderNames),
													}, false),
												},
												"operator": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorContains),
														string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
														string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEquals),
														string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
														string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
																	string(legacyfrontdoor.ManagedRuleExclusionMatchVariableQueryStringArgNames),
																	string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
																	string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestCookieNames),
																	string(legacyfrontdoor.ManagedRuleExclusionMatchVariableRequestHeaderNames),
																}, false),
															},
															"operator": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorContains),
																	string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
																	string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEquals),
																	string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
																	string(legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
														string(legacyfrontdoor.ActionTypeAllow),
														string(legacyfrontdoor.ActionTypeBlock),
														string(legacyfrontdoor.ActionTypeLog),
														string(legacyfrontdoor.ActionTypeRedirect),
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

func resourceCdnFrontdoorFirewallPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorLegacyPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing args for Cdn Frontdoor Firewall Policy")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := parse.NewFrontdoorPolicyID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Cdn Frontdoor Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_firewall_policy", id.ID())
		}
	}

	location := azure.NormalizeLocation("Global")
	enabled := legacyfrontdoor.PolicyEnabledStateDisabled
	if d.Get("enabled").(bool) {
		enabled = legacyfrontdoor.PolicyEnabledStateEnabled
	}
	sku := d.Get("sku_name").(string)
	profileId := d.Get("cdn_frontdoor_profile_id").(string)
	mode := legacyfrontdoor.PolicyMode(d.Get("mode").(string))
	redirectUrl := d.Get("redirect_url").(string)
	customBlockResponseStatusCode := d.Get("custom_block_response_status_code").(int)
	customBlockResponseBody := d.Get("custom_block_response_body").(string)
	customRules := d.Get("custom_rule").([]interface{})
	managedRules := expandFrontDoorFirewallManagedRules(d.Get("managed_rule").([]interface{}))

	if sku == string(legacyfrontdoor.SkuNameStandardAzureFrontDoor) && managedRules != nil {
		return fmt.Errorf("the %q field is only supported with the %q sku, got %q", "managed_rule", legacyfrontdoor.SkuNamePremiumAzureFrontDoor, sku)
	}

	t := d.Get("tags").(map[string]interface{})

	frontdoorWebApplicationFirewallPolicy := legacyfrontdoor.WebApplicationFirewallPolicy{
		Name:     utils.String(name),
		Location: utils.String(location),
		Sku: &legacyfrontdoor.Sku{
			Name: legacyfrontdoor.SkuName(sku),
		},
		WebApplicationFirewallPolicyProperties: &legacyfrontdoor.WebApplicationFirewallPolicyProperties{
			PolicySettings: &legacyfrontdoor.PolicySettings{
				EnabledState: enabled,
				Mode:         mode,
			},
			CustomRules: expandFrontDoorFirewallCustomRules(customRules),
		},
		Tags: ConvertCdnFrontdoorTags(tags.Expand(t)),
	}

	if managedRules != nil {
		frontdoorWebApplicationFirewallPolicy.WebApplicationFirewallPolicyProperties.ManagedRules = managedRules
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName, frontdoorWebApplicationFirewallPolicy)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_profile_id", profileId)
	return resourceCdnFrontdoorFirewallPolicyRead(d, meta)
}

func resourceCdnFrontdoorFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorLegacyPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorPolicyIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Cdn Frontdoor Firewall Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FrontDoorWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroup)

	if profileId := d.Get("cdn_frontdoor_profile_id").(string); profileId != "" {
		d.Set("cdn_frontdoor_profile_id", profileId)
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", string(sku.Name))
	}

	if properties := resp.WebApplicationFirewallPolicyProperties; properties != nil {
		if policy := properties.PolicySettings; policy != nil {
			if policy.EnabledState != "" {
				d.Set("enabled", policy.EnabledState == legacyfrontdoor.PolicyEnabledStateEnabled)
			}
			if policy.Mode != "" {
				d.Set("mode", string(policy.Mode))
			}
			d.Set("redirect_url", policy.RedirectURL)
			d.Set("custom_block_response_status_code", policy.CustomBlockResponseStatusCode)
			d.Set("custom_block_response_body", policy.CustomBlockResponseBody)
		}

		if err := d.Set("custom_rule", flattenFrontDoorFirewallCustomRules(properties.CustomRules)); err != nil {
			return fmt.Errorf("flattening `custom_rule`: %+v", err)
		}

		if err := d.Set("frontend_endpoint_ids", flattenFrontendEndpointLinkSlice(properties.FrontendEndpointLinks)); err != nil {
			return fmt.Errorf("flattening `frontend_endpoint_ids`: %+v", err)
		}

		if err := d.Set("managed_rule", flattenFrontDoorFirewallManagedRules(properties.ManagedRules)); err != nil {
			return fmt.Errorf("flattening `managed_rule`: %+v", err)
		}
	}

	if err := tags.FlattenAndSet(d, ConvertCdnFrontdoorTagsToTagsFlatten(resp.Tags)); err != nil {
		return err
	}

	return nil
}

func resourceCdnFrontdoorFirewallPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorLegacyPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorPolicyIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandFrontDoorFirewallCustomRules(input []interface{}) *legacyfrontdoor.CustomRuleList {
	if len(input) == 0 {
		return nil
	}

	output := make([]legacyfrontdoor.CustomRule, 0)

	for _, cr := range input {
		custom := cr.(map[string]interface{})

		enabled := legacyfrontdoor.CustomRuleEnabledStateDisabled
		if custom["enabled"].(bool) {
			enabled = legacyfrontdoor.CustomRuleEnabledStateEnabled
		}

		name := custom["name"].(string)
		priority := int32(custom["priority"].(int))
		ruleType := custom["type"].(string)
		rateLimitDurationInMinutes := int32(custom["rate_limit_duration_in_minutes"].(int))
		rateLimitThreshold := int32(custom["rate_limit_threshold"].(int))
		matchConditions := expandFrontDoorFirewallMatchConditions(custom["match_condition"].([]interface{}))
		action := custom["action"].(string)

		customRule := legacyfrontdoor.CustomRule{
			Name:                       utils.String(name),
			Priority:                   &priority,
			EnabledState:               enabled,
			RuleType:                   legacyfrontdoor.RuleType(ruleType),
			RateLimitDurationInMinutes: utils.Int32(rateLimitDurationInMinutes),
			RateLimitThreshold:         utils.Int32(rateLimitThreshold),
			MatchConditions:            &matchConditions,
			Action:                     legacyfrontdoor.ActionType(action),
		}
		output = append(output, customRule)
	}

	return &legacyfrontdoor.CustomRuleList{
		Rules: &output,
	}
}

func expandFrontDoorFirewallMatchConditions(input []interface{}) []legacyfrontdoor.MatchCondition {
	if len(input) == 0 {
		return nil
	}

	result := make([]legacyfrontdoor.MatchCondition, 0)

	for _, v := range input {
		match := v.(map[string]interface{})

		matchVariable := match["match_variable"].(string)
		selector := match["selector"].(string)
		operator := match["operator"].(string)
		negateCondition := match["negation_condition"].(bool)
		matchValues := match["match_values"].([]interface{})
		transforms := match["transforms"].([]interface{})

		matchCondition := legacyfrontdoor.MatchCondition{
			Operator:        legacyfrontdoor.Operator(operator),
			NegateCondition: &negateCondition,
			MatchValue:      utils.ExpandStringSlice(matchValues),
			Transforms:      expandFrontDoorFirewallTransforms(transforms),
		}

		if matchVariable != "" {
			matchCondition.MatchVariable = legacyfrontdoor.MatchVariable(matchVariable)
		}
		if selector != "" {
			matchCondition.Selector = utils.String(selector)
		}

		result = append(result, matchCondition)
	}

	return result
}

func expandFrontDoorFirewallTransforms(input []interface{}) *[]legacyfrontdoor.TransformType {
	if len(input) == 0 {
		return nil
	}

	result := make([]legacyfrontdoor.TransformType, 0)
	for _, v := range input {
		result = append(result, legacyfrontdoor.TransformType(v.(string)))
	}

	return &result
}

func expandFrontDoorFirewallManagedRules(input []interface{}) *legacyfrontdoor.ManagedRuleSetList {
	if len(input) == 0 {
		return nil
	}

	managedRules := make([]legacyfrontdoor.ManagedRuleSet, 0)

	for _, mr := range input {
		managedRule := mr.(map[string]interface{})

		ruleType := managedRule["type"].(string)
		version := managedRule["version"].(string)
		overrides := managedRule["override"].([]interface{})
		exclusions := managedRule["exclusion"].([]interface{})

		managedRuleSet := legacyfrontdoor.ManagedRuleSet{
			RuleSetType:    &ruleType,
			RuleSetVersion: &version,
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(exclusions); exclusions != nil {
			managedRuleSet.Exclusions = exclusions
		}

		if ruleGroupOverrides := expandFrontDoorFirewallManagedRuleGroupOverride(overrides); ruleGroupOverrides != nil {
			managedRuleSet.RuleGroupOverrides = ruleGroupOverrides
		}

		managedRules = append(managedRules, managedRuleSet)
	}

	return &legacyfrontdoor.ManagedRuleSetList{
		ManagedRuleSets: &managedRules,
	}
}

func expandFrontDoorFirewallManagedRuleGroupExclusion(input []interface{}) *[]legacyfrontdoor.ManagedRuleExclusion {
	if len(input) == 0 {
		return nil
	}

	managedRuleExclusions := make([]legacyfrontdoor.ManagedRuleExclusion, 0)
	for _, v := range input {
		exclusion := v.(map[string]interface{})

		matchVariable := exclusion["match_variable"].(string)
		operator := exclusion["operator"].(string)
		selector := exclusion["selector"].(string)

		managedRuleExclusion := legacyfrontdoor.ManagedRuleExclusion{
			MatchVariable:         legacyfrontdoor.ManagedRuleExclusionMatchVariable(matchVariable),
			SelectorMatchOperator: legacyfrontdoor.ManagedRuleExclusionSelectorMatchOperator(operator),
			Selector:              &selector,
		}

		managedRuleExclusions = append(managedRuleExclusions, managedRuleExclusion)
	}

	return &managedRuleExclusions
}

func expandFrontDoorFirewallManagedRuleGroupOverride(input []interface{}) *[]legacyfrontdoor.ManagedRuleGroupOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleGroupOverrides := make([]legacyfrontdoor.ManagedRuleGroupOverride, 0)
	for _, v := range input {
		override := v.(map[string]interface{})

		ruleGroupName := override["rule_group_name"].(string)
		rules := override["rule"].([]interface{})
		exclusions := override["exclusion"].([]interface{})

		managedRuleGroupOverride := legacyfrontdoor.ManagedRuleGroupOverride{
			RuleGroupName: &ruleGroupName,
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

func expandFrontDoorFirewallRuleOverride(input []interface{}) *[]legacyfrontdoor.ManagedRuleOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleOverrides := make([]legacyfrontdoor.ManagedRuleOverride, 0)
	for _, v := range input {
		rule := v.(map[string]interface{})

		enabled := legacyfrontdoor.ManagedRuleEnabledStateDisabled
		if rule["enabled"].(bool) {
			enabled = legacyfrontdoor.ManagedRuleEnabledStateEnabled
		}
		ruleId := rule["rule_id"].(string)
		action := legacyfrontdoor.ActionType(rule["action"].(string))
		exclusions := rule["exclusion"].([]interface{})

		managedRuleOverride := legacyfrontdoor.ManagedRuleOverride{
			RuleID:       &ruleId,
			EnabledState: enabled,
			Action:       action,
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(exclusions); exclusions != nil {
			managedRuleOverride.Exclusions = exclusions
		}

		managedRuleOverrides = append(managedRuleOverrides, managedRuleOverride)
	}

	return &managedRuleOverrides
}

func flattenFrontDoorFirewallCustomRules(input *legacyfrontdoor.CustomRuleList) []interface{} {
	if input == nil || input.Rules == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, r := range *input.Rules {
		output := make(map[string]interface{})

		output["name"] = r.Name
		output["type"] = string(r.RuleType)
		output["action"] = string(r.Action)
		if r.EnabledState != "" {
			output["enabled"] = r.EnabledState == legacyfrontdoor.CustomRuleEnabledStateEnabled
		}
		output["match_condition"] = flattenFrontDoorFirewallMatchConditions(*r.MatchConditions)
		output["priority"] = int(*r.Priority)

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

func flattenFrontDoorFirewallMatchConditions(condition []legacyfrontdoor.MatchCondition) []interface{} {
	if condition == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, c := range condition {
		output := make(map[string]interface{})

		output["match_variable"] = string(c.MatchVariable)
		output["operator"] = string(c.Operator)
		output["match_values"] = c.MatchValue
		output["transforms"] = flattenTransformSlice(c.Transforms)

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

func flattenFrontDoorFirewallManagedRules(input *legacyfrontdoor.ManagedRuleSetList) []interface{} {
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

func flattenFrontDoorFirewallExclusions(managedRuleExclusion *[]legacyfrontdoor.ManagedRuleExclusion) []interface{} {
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

func flattenFrontDoorFirewallOverrides(groupOverride *[]legacyfrontdoor.ManagedRuleGroupOverride) []interface{} {
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

func flattenArmFrontdoorFirewallRules(override *[]legacyfrontdoor.ManagedRuleOverride) []interface{} {
	if override == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0)
	for _, o := range *override {
		output := make(map[string]interface{})

		if o.EnabledState != "" {
			output["enabled"] = o.EnabledState == legacyfrontdoor.ManagedRuleEnabledStateEnabled
		}
		if o.Action != "" {
			output["action"] = string(o.Action)
		}

		output["rule_id"] = o.RuleID

		if v := o.Exclusions; v != nil {
			output["exclusion"] = flattenFrontDoorFirewallExclusions(v)
		}

		results = append(results, output)
	}

	return results
}
