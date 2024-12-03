// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2024-02-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorFirewallPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorFirewallPolicyCreate,
		Read:   resourceCdnFrontDoorFirewallPolicyRead,
		Update: resourceCdnFrontDoorFirewallPolicyUpdate,
		Delete: resourceCdnFrontDoorFirewallPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorFirewallPolicyID(id)
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
				ValidateFunc: validate.FrontDoorFirewallPolicyName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(waf.SkuNameStandardAzureFrontDoor),
					string(waf.SkuNamePremiumAzureFrontDoor),
				}, false),
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(waf.PolicyModeDetection),
					string(waf.PolicyModePrevention),
				}, false),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"redirect_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
			},

			"request_body_check_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
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
				ValidateFunc: validation.StringIsBase64,
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
								string(waf.RuleTypeMatchRule),
								string(waf.RuleTypeRateLimitRule),
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
								string(waf.ActionTypeAllow),
								string(waf.ActionTypeBlock),
								string(waf.ActionTypeLog),
								string(waf.ActionTypeRedirect),
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
											string(waf.MatchVariableCookies),
											string(waf.MatchVariablePostArgs),
											string(waf.MatchVariableQueryString),
											string(waf.MatchVariableRemoteAddr),
											string(waf.MatchVariableRequestBody),
											string(waf.MatchVariableRequestHeader),
											string(waf.MatchVariableRequestMethod),
											string(waf.MatchVariableRequestUri),
											string(waf.MatchVariableSocketAddr),
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
											string(waf.OperatorAny),
											string(waf.OperatorBeginsWith),
											string(waf.OperatorContains),
											string(waf.OperatorEndsWith),
											string(waf.OperatorEqual),
											string(waf.OperatorGeoMatch),
											string(waf.OperatorGreaterThan),
											string(waf.OperatorGreaterThanOrEqual),
											string(waf.OperatorIPMatch),
											string(waf.OperatorLessThan),
											string(waf.OperatorLessThanOrEqual),
											string(waf.OperatorRegEx),
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
												string(waf.TransformTypeLowercase),
												string(waf.TransformTypeRemoveNulls),
												string(waf.TransformTypeTrim),
												string(waf.TransformTypeUppercase),
												string(waf.TransformTypeURLDecode),
												string(waf.TransformTypeURLEncode),
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

						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(waf.ActionTypeAllow),
								string(waf.ActionTypeLog),
								string(waf.ActionTypeBlock),
								string(waf.ActionTypeRedirect),
							}, false),
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
											string(waf.ManagedRuleExclusionMatchVariableQueryStringArgNames),
											string(waf.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
											string(waf.ManagedRuleExclusionMatchVariableRequestCookieNames),
											string(waf.ManagedRuleExclusionMatchVariableRequestHeaderNames),
											string(waf.ManagedRuleExclusionMatchVariableRequestBodyJsonArgNames),
										}, false),
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(waf.ManagedRuleExclusionSelectorMatchOperatorContains),
											string(waf.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
											string(waf.ManagedRuleExclusionSelectorMatchOperatorEquals),
											string(waf.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
											string(waf.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
														string(waf.ManagedRuleExclusionMatchVariableQueryStringArgNames),
														string(waf.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
														string(waf.ManagedRuleExclusionMatchVariableRequestCookieNames),
														string(waf.ManagedRuleExclusionMatchVariableRequestHeaderNames),
														string(waf.ManagedRuleExclusionMatchVariableRequestBodyJsonArgNames),
													}, false),
												},
												"operator": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(waf.ManagedRuleExclusionSelectorMatchOperatorContains),
														string(waf.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
														string(waf.ManagedRuleExclusionSelectorMatchOperatorEquals),
														string(waf.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
														string(waf.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
																	string(waf.ManagedRuleExclusionMatchVariableQueryStringArgNames),
																	string(waf.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
																	string(waf.ManagedRuleExclusionMatchVariableRequestCookieNames),
																	string(waf.ManagedRuleExclusionMatchVariableRequestHeaderNames),
																	string(waf.ManagedRuleExclusionMatchVariableRequestBodyJsonArgNames),
																}, false),
															},
															"operator": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(waf.ManagedRuleExclusionSelectorMatchOperatorContains),
																	string(waf.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
																	string(waf.ManagedRuleExclusionSelectorMatchOperatorEquals),
																	string(waf.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
																	string(waf.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
														string(waf.ActionTypeAllow),
														string(waf.ActionTypeLog),
														string(waf.ActionTypeBlock),
														string(waf.ActionTypeRedirect),
														string(waf.ActionTypeAnomalyScoring), // Only valid with 2.0 and above
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

func resourceCdnFrontDoorFirewallPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorFirewallPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := waf.NewFrontDoorWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, name)

	result, err := client.PoliciesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(result.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(result.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_firewall_policy", id.ID())
	}

	enabled := waf.PolicyEnabledStateDisabled

	if d.Get("enabled").(bool) {
		enabled = waf.PolicyEnabledStateEnabled
	}

	requestBodyCheck := waf.PolicyRequestBodyCheckDisabled

	if d.Get("request_body_check_enabled").(bool) {
		requestBodyCheck = waf.PolicyRequestBodyCheckEnabled
	}

	sku := d.Get("sku_name").(string)
	mode := waf.PolicyMode(d.Get("mode").(string))
	redirectUrl := d.Get("redirect_url").(string)
	customBlockResponseStatusCode := d.Get("custom_block_response_status_code").(int)
	customBlockResponseBody := d.Get("custom_block_response_body").(string)
	customRules := d.Get("custom_rule").([]interface{})
	managedRules, err := expandCdnFrontDoorFirewallManagedRules(d.Get("managed_rule").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding managed_rule: %+v", err)
	}

	if sku != string(waf.SkuNamePremiumAzureFrontDoor) && managedRules != nil {
		return fmt.Errorf("the 'managed_rule' field is only supported with the 'Premium_AzureFrontDoor' sku, got %q", sku)
	}

	t := d.Get("tags").(map[string]interface{})

	payload := waf.WebApplicationFirewallPolicy{
		Location: pointer.To(location.Normalize("Global")),
		Sku: &waf.Sku{
			Name: pointer.To(waf.SkuName(sku)),
		},
		Properties: &waf.WebApplicationFirewallPolicyProperties{
			PolicySettings: &waf.PolicySettings{
				EnabledState:     pointer.To(enabled),
				Mode:             pointer.To(mode),
				RequestBodyCheck: pointer.To(requestBodyCheck),
				// NOTE: Add 'javascript_challenge_expiration_in_minutes' here...
			},
			CustomRules: expandCdnFrontDoorFirewallCustomRules(customRules),
		},
		Tags: tags.Expand(t),
	}

	if managedRules != nil {
		payload.Properties.ManagedRules = managedRules
	}

	if redirectUrl != "" {
		payload.Properties.PolicySettings.RedirectURL = pointer.To(redirectUrl)
	}

	if customBlockResponseBody != "" {
		payload.Properties.PolicySettings.CustomBlockResponseBody = pointer.To(customBlockResponseBody)
	}

	if customBlockResponseStatusCode > 0 {
		payload.Properties.PolicySettings.CustomBlockResponseStatusCode = pointer.To(int64(customBlockResponseStatusCode))
	}

	err = client.PoliciesCreateOrUpdateThenPoll(ctx, id, payload)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorFirewallPolicyRead(d, meta)
}

func resourceCdnFrontDoorFirewallPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorFirewallPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := waf.ParseFrontDoorWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	result, err := client.PoliciesGet(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := result.Model

	if model == nil {
		return fmt.Errorf("retrieving %s: 'model' was nil", *id)
	}

	if model.Sku == nil {
		return fmt.Errorf("retrieving %s: 'model.Sku' was nil", *id)
	}

	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: 'model.Properties' was nil", *id)
	}

	props := *model.Properties

	if d.HasChanges("custom_block_response_body", "custom_block_response_status_code", "enabled", "mode", "redirect_url", "request_body_check_enabled") {
		enabled := waf.PolicyEnabledStateDisabled
		if d.Get("enabled").(bool) {
			enabled = waf.PolicyEnabledStateEnabled
		}

		requestBodyCheck := waf.PolicyRequestBodyCheckDisabled
		if d.Get("request_body_check_enabled").(bool) {
			requestBodyCheck = waf.PolicyRequestBodyCheckEnabled
		}
		props.PolicySettings = pointer.To(waf.PolicySettings{
			EnabledState:     pointer.To(enabled),
			Mode:             pointer.To(waf.PolicyMode(d.Get("mode").(string))),
			RequestBodyCheck: pointer.To(requestBodyCheck),
		})

		if redirectUrl := d.Get("redirect_url").(string); redirectUrl != "" {
			props.PolicySettings.RedirectURL = pointer.To(redirectUrl)
		}

		if body := d.Get("custom_block_response_body").(string); body != "" {
			props.PolicySettings.CustomBlockResponseBody = pointer.To(body)
		}

		if statusCode := d.Get("custom_block_response_status_code").(int64); statusCode > 0 {
			props.PolicySettings.CustomBlockResponseStatusCode = pointer.To(statusCode)
		}
	}

	if d.HasChange("custom_rule") {
		props.CustomRules = expandCdnFrontDoorFirewallCustomRules(d.Get("custom_rule").([]interface{}))
	}

	if d.HasChange("managed_rule") {
		managedRules, err := expandCdnFrontDoorFirewallManagedRules(d.Get("managed_rule").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding managed_rule: %+v", err)
		}

		if pointer.From(model.Sku.Name) != waf.SkuNamePremiumAzureFrontDoor && managedRules != nil {
			return fmt.Errorf("the 'managed_rule' field is only supported when using the sku 'Premium_AzureFrontDoor', got %q", pointer.From(model.Sku.Name))
		}

		if managedRules != nil {
			props.ManagedRules = managedRules
		}
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		model.Tags = tags.Expand(t)
	}

	model.Properties = pointer.To(props)

	err = client.PoliciesCreateOrUpdateThenPoll(ctx, pointer.From(id), pointer.From(model))
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorFirewallPolicyRead(d, meta)
}

func resourceCdnFrontDoorFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorFirewallPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := waf.ParseFrontDoorWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	result, err := client.PoliciesGet(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := result.Model

	if model == nil {
		return fmt.Errorf("retrieving %s: 'model' was nil", *id)
	}

	if model.Sku == nil {
		return fmt.Errorf("retrieving %s: 'model.Sku' was nil", *id)
	}

	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: 'model.Properties' was nil", *id)
	}

	d.Set("name", id.FrontDoorWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	skuName := ""
	if sku := model.Sku; sku != nil {
		skuName = string(pointer.From(sku.Name))
	}
	d.Set("sku_name", skuName)

	if props := model.Properties; props != nil {
		if policy := props.PolicySettings; policy != nil {
			d.Set("enabled", pointer.From(policy.EnabledState) == waf.PolicyEnabledStateEnabled)
			d.Set("mode", string(pointer.From(policy.Mode)))
			d.Set("request_body_check_enabled", pointer.From(policy.RequestBodyCheck) == waf.PolicyRequestBodyCheckEnabled)
			d.Set("redirect_url", policy.RedirectURL)
			d.Set("custom_block_response_status_code", policy.CustomBlockResponseStatusCode)
			d.Set("custom_block_response_body", policy.CustomBlockResponseBody)
		}

		if err := d.Set("custom_rule", flattenCdnFrontDoorFirewallCustomRules(props.CustomRules)); err != nil {
			return fmt.Errorf("flattening 'custom_rule': %+v", err)
		}

		if err := d.Set("frontend_endpoint_ids", flattenFrontendEndpointLinkSlice(props.FrontendEndpointLinks)); err != nil {
			return fmt.Errorf("flattening 'frontend_endpoint_ids': %+v", err)
		}

		if err := d.Set("managed_rule", flattenCdnFrontDoorFirewallManagedRules(props.ManagedRules)); err != nil {
			return fmt.Errorf("flattening 'managed_rule': %+v", err)
		}
	}

	if err := tags.FlattenAndSet(d, model.Tags); err != nil {
		return err
	}

	return nil
}

func resourceCdnFrontDoorFirewallPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorLegacyFirewallPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorFirewallPolicyID(d.Id())
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

func expandCdnFrontDoorFirewallCustomRules(input []interface{}) *waf.CustomRuleList {
	if len(input) == 0 {
		return nil
	}

	output := make([]waf.CustomRule, 0)

	for _, cr := range input {
		custom := cr.(map[string]interface{})

		enabled := waf.CustomRuleEnabledStateDisabled
		if custom["enabled"].(bool) {
			enabled = waf.CustomRuleEnabledStateEnabled
		}

		name := custom["name"].(string)
		priority := int64(custom["priority"].(int))
		ruleType := custom["type"].(string)
		rateLimitDurationInMinutes := int64(custom["rate_limit_duration_in_minutes"].(int))
		rateLimitThreshold := int64(custom["rate_limit_threshold"].(int))
		matchConditions := expandCdnFrontDoorFirewallMatchConditions(custom["match_condition"].([]interface{}))
		action := custom["action"].(string)

		output = append(output, waf.CustomRule{
			Name:                       pointer.To(name),
			Priority:                   priority,
			EnabledState:               pointer.To(enabled),
			RuleType:                   waf.RuleType(ruleType),
			RateLimitDurationInMinutes: pointer.To(rateLimitDurationInMinutes),
			RateLimitThreshold:         pointer.To(rateLimitThreshold),
			MatchConditions:            matchConditions,
			Action:                     waf.ActionType(action),
		})
	}

	return &waf.CustomRuleList{
		Rules: &output,
	}
}

func expandCdnFrontDoorFirewallMatchConditions(input []interface{}) []waf.MatchCondition {
	result := make([]waf.MatchCondition, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		match := v.(map[string]interface{})

		matchVariable := match["match_variable"].(string)
		selector := match["selector"].(string)
		operator := match["operator"].(string)
		negateCondition := match["negation_condition"].(bool)
		matchValues := match["match_values"].([]string)
		transforms := match["transforms"].([]interface{})

		matchCondition := waf.MatchCondition{
			Operator:        waf.Operator(operator),
			NegateCondition: &negateCondition,
			MatchValue:      matchValues,
			Transforms:      expandCdnFrontDoorFirewallTransforms(transforms),
		}

		if matchVariable != "" {
			matchCondition.MatchVariable = waf.MatchVariable(matchVariable)
		}
		if selector != "" {
			matchCondition.Selector = pointer.To(selector)
		}

		result = append(result, matchCondition)
	}

	return result
}

func expandCdnFrontDoorFirewallTransforms(input []interface{}) *[]waf.TransformType {
	result := make([]waf.TransformType, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		result = append(result, waf.TransformType(v.(string)))
	}

	return &result
}

func expandCdnFrontDoorFirewallManagedRules(input []interface{}) (*waf.ManagedRuleSetList, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]waf.ManagedRuleSet, 0)
	for _, mr := range input {
		managedRule := mr.(map[string]interface{})

		ruleType := managedRule["type"].(string)
		version := managedRule["version"].(string)
		action := managedRule["action"].(string)
		overrides := managedRule["override"].([]interface{})
		exclusions := expandCdnFrontDoorFirewallManagedRuleGroupExclusion(managedRule["exclusion"].([]interface{}))

		fVersion := 1.0
		if v, err := strconv.ParseFloat(version, 64); err == nil {
			fVersion = v
		}

		// NOTE: The API is deferring the version range from the rule type name
		// 'DefaultRuleSet' is < 1.1 and 'Microsoft_DefaultRuleSet' >= 1.1
		// 'AnomalyScoring' action only valid on 2.0 and above
		if ruleType == "DefaultRuleSet" && fVersion > 1.0 {
			return nil, fmt.Errorf("the managed rule set type %q and version %q is not supported. If you wish to use the 'DefaultRuleSet' type please update your 'version' field to be '1.0' or 'preview-0.1', got %q", ruleType, version, version)
		} else if ruleType == "Microsoft_DefaultRuleSet" && fVersion < 1.1 {
			return nil, fmt.Errorf("the managed rule set type %q and version %q is not supported. If you wish to use the 'Microsoft_DefaultRuleSet' type please update your 'version' field to be '1.1', '2.0' or '2.1', got %q", ruleType, version, version)
		}

		ruleGroupOverrides, err := expandCdnFrontDoorFirewallManagedRuleGroupOverride(overrides, version, fVersion)
		if err != nil {
			return nil, err
		}

		managedRuleSet := waf.ManagedRuleSet{
			Exclusions:         exclusions,
			RuleSetVersion:     version,
			RuleGroupOverrides: ruleGroupOverrides,
			RuleSetType:        ruleType,
		}

		if action != "" {
			managedRuleSet.RuleSetAction = pointer.To(waf.ManagedRuleSetActionType(action))
		}

		result = append(result, managedRuleSet)
	}

	return &waf.ManagedRuleSetList{
		ManagedRuleSets: &result,
	}, nil
}

func expandCdnFrontDoorFirewallManagedRuleGroupExclusion(input []interface{}) *[]waf.ManagedRuleExclusion {
	results := make([]waf.ManagedRuleExclusion, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		exclusion := v.(map[string]interface{})

		matchVariable := exclusion["match_variable"].(string)
		operator := exclusion["operator"].(string)
		selector := exclusion["selector"].(string)

		results = append(results, waf.ManagedRuleExclusion{
			MatchVariable:         waf.ManagedRuleExclusionMatchVariable(matchVariable),
			SelectorMatchOperator: waf.ManagedRuleExclusionSelectorMatchOperator(operator),
			Selector:              selector,
		})
	}

	return &results
}

func expandCdnFrontDoorFirewallManagedRuleGroupOverride(input []interface{}, versionRaw string, version float64) (*[]waf.ManagedRuleGroupOverride, error) {
	result := make([]waf.ManagedRuleGroupOverride, 0)
	if len(input) == 0 {
		return nil, nil
	}

	for _, v := range input {
		override := v.(map[string]interface{})

		exclusions := expandCdnFrontDoorFirewallManagedRuleGroupExclusion(override["exclusion"].([]interface{}))
		ruleGroupName := override["rule_group_name"].(string)
		rules, err := expandCdnFrontDoorFirewallRuleOverride(override["rule"].([]interface{}), versionRaw, version)
		if err != nil {
			return nil, err
		}

		result = append(result, waf.ManagedRuleGroupOverride{
			Exclusions:    exclusions,
			RuleGroupName: ruleGroupName,
			Rules:         rules,
		})
	}

	return &result, nil
}

func expandCdnFrontDoorFirewallRuleOverride(input []interface{}, versionRaw string, version float64) (*[]waf.ManagedRuleOverride, error) {
	result := make([]waf.ManagedRuleOverride, 0)
	if len(input) == 0 {
		return nil, nil
	}

	for _, v := range input {
		rule := v.(map[string]interface{})

		enabled := waf.ManagedRuleEnabledStateDisabled
		if rule["enabled"].(bool) {
			enabled = waf.ManagedRuleEnabledStateEnabled
		}

		ruleId := rule["rule_id"].(string)
		actionTypeRaw := rule["action"].(string)
		action := waf.ActionType(actionTypeRaw)

		// NOTE: Default Rule Sets(DRS) 2.0 and above rules only use action type of 'AnomalyScoring' or 'Log'. Issues 19088 and 19561
		// This will still work for bot rules as well since it will be the default value of 1.0
		if version < 2.0 && actionTypeRaw == string(waf.ActionTypeAnomalyScoring) {
			return nil, fmt.Errorf("'AnomalyScoring' is only valid in managed rules that are DRS 2.0 and above, got %q", versionRaw)
		} else if version >= 2.0 && actionTypeRaw != string(waf.ActionTypeAnomalyScoring) && actionTypeRaw != "Log" {
			return nil, fmt.Errorf("the managed rules 'action' field must be set to 'AnomalyScoring' or 'Log' if the managed rule is DRS 2.0 or above, got %q", action)
		}

		exclusions := expandCdnFrontDoorFirewallManagedRuleGroupExclusion(rule["exclusion"].([]interface{}))

		result = append(result, waf.ManagedRuleOverride{
			RuleId:       ruleId,
			EnabledState: pointer.To(enabled),
			Action:       pointer.To(action),
			Exclusions:   exclusions,
		})
	}

	return &result, nil
}

func flattenCdnFrontDoorFirewallCustomRules(input *waf.CustomRuleList) []interface{} {
	if input == nil || input.Rules == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input.Rules {
		action := ""
		if v.Action != "" {
			action = string(v.Action)
		}

		enabled := false
		if v.EnabledState != nil {
			enabled = pointer.From(v.EnabledState) == waf.CustomRuleEnabledStateEnabled
		}

		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		priority := 0
		if v.Priority != 0 {
			priority = int(v.Priority)
		}

		rateLimitDurationInMinutes := 0
		if v.RateLimitDurationInMinutes != nil {
			rateLimitDurationInMinutes = int(*v.RateLimitDurationInMinutes)
		}

		rateLimitThreshold := 0
		if v.RateLimitThreshold != nil {
			rateLimitThreshold = int(*v.RateLimitThreshold)
		}

		ruleType := ""
		if v.RuleType != "" {
			ruleType = string(v.RuleType)
		}

		results = append(results, map[string]interface{}{
			"action":                         action,
			"enabled":                        enabled,
			"match_condition":                flattenCdnFrontDoorFirewallMatchConditions(pointer.To(v.MatchConditions)),
			"rate_limit_duration_in_minutes": rateLimitDurationInMinutes,
			"rate_limit_threshold":           rateLimitThreshold,
			"priority":                       priority,
			"name":                           name,
			"type":                           ruleType,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallMatchConditions(input *[]waf.MatchCondition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		selector := ""
		if v.Selector != nil {
			selector = *v.Selector
		}

		negateCondition := false
		if v.NegateCondition != nil {
			negateCondition = *v.NegateCondition
		}

		results = append(results, map[string]interface{}{
			"match_variable":     string(v.MatchVariable),
			"match_values":       v.MatchValue,
			"negation_condition": negateCondition,
			"operator":           string(v.Operator),
			"selector":           selector,
			"transforms":         flattenTransformSlice(v.Transforms),
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallManagedRules(input *waf.ManagedRuleSetList) []interface{} {
	if input == nil || input.ManagedRuleSets == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, r := range *input.ManagedRuleSets {
		ruleSetType := ""
		if r.RuleSetType != "" {
			ruleSetType = r.RuleSetType
		}

		ruleSetVersion := ""
		if r.RuleSetVersion != "" {
			ruleSetVersion = r.RuleSetVersion
		}

		ruleSetAction := ""
		if r.RuleSetAction != nil {
			ruleSetAction = string(pointer.From(r.RuleSetAction))
		}

		results = append(results, map[string]interface{}{
			"exclusion": flattenCdnFrontDoorFirewallExclusions(r.Exclusions),
			"override":  flattenCdnFrontDoorFirewallOverrides(r.RuleGroupOverrides),
			"type":      ruleSetType,
			"version":   ruleSetVersion,
			"action":    ruleSetAction,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallExclusions(input *[]waf.ManagedRuleExclusion) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		matchVariable := ""
		if v.MatchVariable != "" {
			matchVariable = string(v.MatchVariable)
		}

		operator := ""
		if v.SelectorMatchOperator != "" {
			operator = string(v.SelectorMatchOperator)
		}

		selector := ""
		if v.Selector != "" {
			selector = v.Selector
		}

		results = append(results, map[string]interface{}{
			"match_variable": matchVariable,
			"operator":       operator,
			"selector":       selector,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallOverrides(input *[]waf.ManagedRuleGroupOverride) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		ruleGroupName := ""
		if v.RuleGroupName != "" {
			ruleGroupName = v.RuleGroupName
		}

		results = append(results, map[string]interface{}{
			"rule_group_name": ruleGroupName,
			"exclusion":       flattenCdnFrontDoorFirewallExclusions(v.Exclusions),
			"rule":            flattenCdnFrontDoorFirewallRules(v.Rules),
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallRules(input *[]waf.ManagedRuleOverride) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		action := waf.ActionTypeAnomalyScoring
		if v.Action != nil {
			action = pointer.From(v.Action)
		}

		enabled := false
		if v.EnabledState != nil {
			enabled = pointer.From(v.EnabledState) == waf.ManagedRuleEnabledStateEnabled
		}

		ruleId := ""
		if v.RuleId != "" {
			ruleId = v.RuleId
		}

		results = append(results, map[string]interface{}{
			"action":    action,
			"enabled":   enabled,
			"exclusion": flattenCdnFrontDoorFirewallExclusions(v.Exclusions),
			"rule_id":   ruleId,
		})
	}

	return results
}
