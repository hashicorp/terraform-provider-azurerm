// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
					string(frontdoor.SkuNameStandardAzureFrontDoor),
					string(frontdoor.SkuNamePremiumAzureFrontDoor),
				}, false),
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(frontdoor.PolicyModeDetection),
					string(frontdoor.PolicyModePrevention),
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
								string(frontdoor.RuleTypeMatchRule),
								string(frontdoor.RuleTypeRateLimitRule),
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
								string(frontdoor.ActionTypeAllow),
								string(frontdoor.ActionTypeBlock),
								string(frontdoor.ActionTypeLog),
								string(frontdoor.ActionTypeRedirect),
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
											string(frontdoor.MatchVariableCookies),
											string(frontdoor.MatchVariablePostArgs),
											string(frontdoor.MatchVariableQueryString),
											string(frontdoor.MatchVariableRemoteAddr),
											string(frontdoor.MatchVariableRequestBody),
											string(frontdoor.MatchVariableRequestHeader),
											string(frontdoor.MatchVariableRequestMethod),
											string(frontdoor.MatchVariableRequestURI),
											string(frontdoor.MatchVariableSocketAddr),
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
												string(frontdoor.TransformTypeLowercase),
												string(frontdoor.TransformTypeRemoveNulls),
												string(frontdoor.TransformTypeTrim),
												string(frontdoor.TransformTypeUppercase),
												string(frontdoor.TransformTypeURLDecode),
												string(frontdoor.TransformTypeURLEncode),
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
								string(frontdoor.ActionTypeAllow),
								string(frontdoor.ActionTypeLog),
								string(frontdoor.ActionTypeBlock),
								string(frontdoor.ActionTypeRedirect),
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
											string(frontdoor.ManagedRuleExclusionMatchVariableQueryStringArgNames),
											string(frontdoor.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
											string(frontdoor.ManagedRuleExclusionMatchVariableRequestCookieNames),
											string(frontdoor.ManagedRuleExclusionMatchVariableRequestHeaderNames),
											string(frontdoor.ManagedRuleExclusionMatchVariableRequestBodyJSONArgNames),
										}, false),
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorContains),
											string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
											string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEquals),
											string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
											string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
														string(frontdoor.ManagedRuleExclusionMatchVariableQueryStringArgNames),
														string(frontdoor.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
														string(frontdoor.ManagedRuleExclusionMatchVariableRequestCookieNames),
														string(frontdoor.ManagedRuleExclusionMatchVariableRequestHeaderNames),
														string(frontdoor.ManagedRuleExclusionMatchVariableRequestBodyJSONArgNames),
													}, false),
												},
												"operator": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorContains),
														string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
														string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEquals),
														string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
														string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
																	string(frontdoor.ManagedRuleExclusionMatchVariableQueryStringArgNames),
																	string(frontdoor.ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
																	string(frontdoor.ManagedRuleExclusionMatchVariableRequestCookieNames),
																	string(frontdoor.ManagedRuleExclusionMatchVariableRequestHeaderNames),
																	string(frontdoor.ManagedRuleExclusionMatchVariableRequestBodyJSONArgNames),
																}, false),
															},
															"operator": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorContains),
																	string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEndsWith),
																	string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEquals),
																	string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
																	string(frontdoor.ManagedRuleExclusionSelectorMatchOperatorStartsWith),
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
														string(frontdoor.ActionTypeAllow),
														string(frontdoor.ActionTypeLog),
														string(frontdoor.ActionTypeBlock),
														string(frontdoor.ActionTypeRedirect),
														"AnomalyScoring", // Only valid with 2.0 and above
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
	client := meta.(*clients.Client).Cdn.FrontDoorLegacyFirewallPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewFrontDoorFirewallPolicyID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_firewall_policy", id.ID())
	}

	enabled := frontdoor.PolicyEnabledStateDisabled

	if d.Get("enabled").(bool) {
		enabled = frontdoor.PolicyEnabledStateEnabled
	}

	requestBodyCheck := frontdoor.PolicyRequestBodyCheckDisabled

	if d.Get("request_body_check_enabled").(bool) {
		requestBodyCheck = frontdoor.PolicyRequestBodyCheckEnabled
	}

	sku := d.Get("sku_name").(string)
	mode := frontdoor.PolicyMode(d.Get("mode").(string))
	redirectUrl := d.Get("redirect_url").(string)
	customBlockResponseStatusCode := d.Get("custom_block_response_status_code").(int)
	customBlockResponseBody := d.Get("custom_block_response_body").(string)
	customRules := d.Get("custom_rule").([]interface{})
	managedRules, err := expandCdnFrontDoorFirewallManagedRules(d.Get("managed_rule").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding managed_rule: %+v", err)
	}

	if sku != string(frontdoor.SkuNamePremiumAzureFrontDoor) && managedRules != nil {
		return fmt.Errorf("the 'managed_rule' field is only supported with the 'Premium_AzureFrontDoor' sku, got %q", sku)
	}

	t := d.Get("tags").(map[string]interface{})

	payload := frontdoor.WebApplicationFirewallPolicy{
		Location: utils.String(location.Normalize("Global")),
		Sku: &frontdoor.Sku{
			Name: frontdoor.SkuName(sku),
		},
		WebApplicationFirewallPolicyProperties: &frontdoor.WebApplicationFirewallPolicyProperties{
			PolicySettings: &frontdoor.PolicySettings{
				EnabledState:     enabled,
				Mode:             mode,
				RequestBodyCheck: requestBodyCheck,
			},
			CustomRules: expandCdnFrontDoorFirewallCustomRules(customRules),
		},
		Tags: expandFrontDoorTags(tags.Expand(t)),
	}

	if managedRules != nil {
		payload.WebApplicationFirewallPolicyProperties.ManagedRules = managedRules
	}

	if redirectUrl != "" {
		payload.WebApplicationFirewallPolicyProperties.PolicySettings.RedirectURL = utils.String(redirectUrl)
	}

	if customBlockResponseBody != "" {
		payload.WebApplicationFirewallPolicyProperties.PolicySettings.CustomBlockResponseBody = utils.String(customBlockResponseBody)
	}

	if customBlockResponseStatusCode > 0 {
		payload.WebApplicationFirewallPolicyProperties.PolicySettings.CustomBlockResponseStatusCode = utils.Int32(int32(customBlockResponseStatusCode))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName, payload)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorFirewallPolicyRead(d, meta)
}

func resourceCdnFrontDoorFirewallPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorLegacyFirewallPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Sku == nil {
		return fmt.Errorf("retrieving %s: 'sku' was nil", *id)
	}

	if existing.WebApplicationFirewallPolicyProperties == nil {
		return fmt.Errorf("retrieving %s: 'properties' was nil", *id)
	}

	props := *existing.WebApplicationFirewallPolicyProperties

	if d.HasChanges("custom_block_response_body", "custom_block_response_status_code", "enabled", "mode", "redirect_url", "request_body_check_enabled") {
		enabled := frontdoor.PolicyEnabledStateDisabled
		if d.Get("enabled").(bool) {
			enabled = frontdoor.PolicyEnabledStateEnabled
		}
		requestBodyCheck := frontdoor.PolicyRequestBodyCheckDisabled
		if d.Get("request_body_check_enabled").(bool) {
			requestBodyCheck = frontdoor.PolicyRequestBodyCheckEnabled
		}
		props.PolicySettings = &frontdoor.PolicySettings{
			EnabledState:     enabled,
			Mode:             frontdoor.PolicyMode(d.Get("mode").(string)),
			RequestBodyCheck: requestBodyCheck,
		}

		if redirectUrl := d.Get("redirect_url").(string); redirectUrl != "" {
			props.PolicySettings.RedirectURL = utils.String(redirectUrl)
		}

		if body := d.Get("custom_block_response_body").(string); body != "" {
			props.PolicySettings.CustomBlockResponseBody = utils.String(body)
		}

		if statusCode := d.Get("custom_block_response_status_code").(int); statusCode > 0 {
			props.PolicySettings.CustomBlockResponseStatusCode = utils.Int32(int32(statusCode))
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

		if existing.Sku.Name != frontdoor.SkuNamePremiumAzureFrontDoor && managedRules != nil {
			return fmt.Errorf("the 'managed_rule' field is only supported when using the sku 'Premium_AzureFrontDoor', got %q", existing.Sku.Name)
		}

		if managedRules != nil {
			props.ManagedRules = managedRules
		}
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		existing.Tags = expandFrontDoorTags(tags.Expand(t))
	}

	existing.WebApplicationFirewallPolicyProperties = &props
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName, existing)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorFirewallPolicyRead(d, meta)
}

func resourceCdnFrontDoorFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorLegacyFirewallPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorFirewallPolicyID(d.Id())
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

	skuName := ""
	if sku := resp.Sku; sku != nil {
		skuName = string(sku.Name)
	}
	d.Set("sku_name", skuName)

	if properties := resp.WebApplicationFirewallPolicyProperties; properties != nil {
		if policy := properties.PolicySettings; policy != nil {
			d.Set("enabled", policy.EnabledState == frontdoor.PolicyEnabledStateEnabled)
			d.Set("mode", string(policy.Mode))
			d.Set("request_body_check_enabled", policy.RequestBodyCheck == frontdoor.PolicyRequestBodyCheckEnabled)
			d.Set("redirect_url", policy.RedirectURL)
			d.Set("custom_block_response_status_code", policy.CustomBlockResponseStatusCode)
			d.Set("custom_block_response_body", policy.CustomBlockResponseBody)
		}

		if err := d.Set("custom_rule", flattenCdnFrontDoorFirewallCustomRules(properties.CustomRules)); err != nil {
			return fmt.Errorf("flattening 'custom_rule': %+v", err)
		}

		if err := d.Set("frontend_endpoint_ids", flattenFrontendEndpointLinkSlice(properties.FrontendEndpointLinks)); err != nil {
			return fmt.Errorf("flattening 'frontend_endpoint_ids': %+v", err)
		}

		if err := d.Set("managed_rule", flattenCdnFrontDoorFirewallManagedRules(properties.ManagedRules)); err != nil {
			return fmt.Errorf("flattening 'managed_rule': %+v", err)
		}
	}

	if err := tags.FlattenAndSet(d, flattenFrontDoorTags(resp.Tags)); err != nil {
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

func expandCdnFrontDoorFirewallCustomRules(input []interface{}) *frontdoor.CustomRuleList {
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
		matchConditions := expandCdnFrontDoorFirewallMatchConditions(custom["match_condition"].([]interface{}))
		action := custom["action"].(string)

		output = append(output, frontdoor.CustomRule{
			Name:                       utils.String(name),
			Priority:                   &priority,
			EnabledState:               enabled,
			RuleType:                   frontdoor.RuleType(ruleType),
			RateLimitDurationInMinutes: utils.Int32(rateLimitDurationInMinutes),
			RateLimitThreshold:         utils.Int32(rateLimitThreshold),
			MatchConditions:            &matchConditions,
			Action:                     frontdoor.ActionType(action),
		})
	}

	return &frontdoor.CustomRuleList{
		Rules: &output,
	}
}

func expandCdnFrontDoorFirewallMatchConditions(input []interface{}) []frontdoor.MatchCondition {
	result := make([]frontdoor.MatchCondition, 0)
	if len(input) == 0 {
		return nil
	}

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
			Transforms:      expandCdnFrontDoorFirewallTransforms(transforms),
		}

		if matchVariable != "" {
			matchCondition.MatchVariable = frontdoor.MatchVariable(matchVariable)
		}
		if selector != "" {
			matchCondition.Selector = utils.String(selector)
		}

		result = append(result, matchCondition)
	}

	return result
}

func expandCdnFrontDoorFirewallTransforms(input []interface{}) *[]frontdoor.TransformType {
	result := make([]frontdoor.TransformType, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		result = append(result, frontdoor.TransformType(v.(string)))
	}

	return &result
}

func expandCdnFrontDoorFirewallManagedRules(input []interface{}) (*frontdoor.ManagedRuleSetList, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]frontdoor.ManagedRuleSet, 0)
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

		managedRuleSet := frontdoor.ManagedRuleSet{
			Exclusions:         exclusions,
			RuleSetVersion:     &version,
			RuleGroupOverrides: ruleGroupOverrides,
			RuleSetType:        &ruleType,
		}

		if action != "" {
			managedRuleSet.RuleSetAction = frontdoor.ManagedRuleSetActionType(action)
		}

		result = append(result, managedRuleSet)
	}

	return &frontdoor.ManagedRuleSetList{
		ManagedRuleSets: &result,
	}, nil
}

func expandCdnFrontDoorFirewallManagedRuleGroupExclusion(input []interface{}) *[]frontdoor.ManagedRuleExclusion {
	results := make([]frontdoor.ManagedRuleExclusion, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		exclusion := v.(map[string]interface{})

		matchVariable := exclusion["match_variable"].(string)
		operator := exclusion["operator"].(string)
		selector := exclusion["selector"].(string)

		results = append(results, frontdoor.ManagedRuleExclusion{
			MatchVariable:         frontdoor.ManagedRuleExclusionMatchVariable(matchVariable),
			SelectorMatchOperator: frontdoor.ManagedRuleExclusionSelectorMatchOperator(operator),
			Selector:              &selector,
		})
	}

	return &results
}

func expandCdnFrontDoorFirewallManagedRuleGroupOverride(input []interface{}, versionRaw string, version float64) (*[]frontdoor.ManagedRuleGroupOverride, error) {
	result := make([]frontdoor.ManagedRuleGroupOverride, 0)
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

		result = append(result, frontdoor.ManagedRuleGroupOverride{
			Exclusions:    exclusions,
			RuleGroupName: &ruleGroupName,
			Rules:         rules,
		})
	}

	return &result, nil
}

func expandCdnFrontDoorFirewallRuleOverride(input []interface{}, versionRaw string, version float64) (*[]frontdoor.ManagedRuleOverride, error) {
	result := make([]frontdoor.ManagedRuleOverride, 0)
	if len(input) == 0 {
		return nil, nil
	}

	for _, v := range input {
		rule := v.(map[string]interface{})

		enabled := frontdoor.ManagedRuleEnabledStateDisabled
		if rule["enabled"].(bool) {
			enabled = frontdoor.ManagedRuleEnabledStateEnabled
		}

		ruleId := rule["rule_id"].(string)
		actionTypeRaw := rule["action"].(string)
		action := frontdoor.ActionType(actionTypeRaw)

		// NOTE: Default Rule Sets(DRS) 2.0 and above rules only use action type of 'AnomalyScoring' or 'Log'. Issues 19088 and 19561
		// This will still work for bot rules as well since it will be the default value of 1.0
		if version < 2.0 && actionTypeRaw == "AnomalyScoring" {
			return nil, fmt.Errorf("'AnomalyScoring' is only valid in managed rules that are DRS 2.0 and above, got %q", versionRaw)
		} else if version >= 2.0 && actionTypeRaw != "AnomalyScoring" && actionTypeRaw != "Log" {
			return nil, fmt.Errorf("the managed rules 'action' field must be set to 'AnomalyScoring' or 'Log' if the managed rule is DRS 2.0 or above, got %q", action)
		}

		exclusions := expandCdnFrontDoorFirewallManagedRuleGroupExclusion(rule["exclusion"].([]interface{}))

		result = append(result, frontdoor.ManagedRuleOverride{
			RuleID:       &ruleId,
			EnabledState: enabled,
			Action:       action,
			Exclusions:   exclusions,
		})
	}

	return &result, nil
}

func flattenCdnFrontDoorFirewallCustomRules(input *frontdoor.CustomRuleList) []interface{} {
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
		if v.EnabledState != "" {
			enabled = v.EnabledState == frontdoor.CustomRuleEnabledStateEnabled
		}

		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		priority := 0
		if v.Priority != nil {
			priority = int(*v.Priority)
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
			"match_condition":                flattenCdnFrontDoorFirewallMatchConditions(v.MatchConditions),
			"rate_limit_duration_in_minutes": rateLimitDurationInMinutes,
			"rate_limit_threshold":           rateLimitThreshold,
			"priority":                       priority,
			"name":                           name,
			"type":                           ruleType,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallMatchConditions(input *[]frontdoor.MatchCondition) []interface{} {
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

func flattenCdnFrontDoorFirewallManagedRules(input *frontdoor.ManagedRuleSetList) []interface{} {
	if input == nil || input.ManagedRuleSets == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, r := range *input.ManagedRuleSets {
		ruleSetType := ""
		if r.RuleSetType != nil {
			ruleSetType = *r.RuleSetType
		}

		ruleSetVersion := ""
		if r.RuleSetVersion != nil {
			ruleSetVersion = *r.RuleSetVersion
		}

		ruleSetAction := ""
		if r.RuleSetAction != "" {
			ruleSetAction = string(r.RuleSetAction)
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

func flattenCdnFrontDoorFirewallExclusions(input *[]frontdoor.ManagedRuleExclusion) []interface{} {
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
		if v.Selector != nil {
			selector = *v.Selector
		}

		results = append(results, map[string]interface{}{
			"match_variable": matchVariable,
			"operator":       operator,
			"selector":       selector,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallOverrides(input *[]frontdoor.ManagedRuleGroupOverride) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		ruleGroupName := ""
		if v.RuleGroupName != nil {
			ruleGroupName = *v.RuleGroupName
		}

		results = append(results, map[string]interface{}{
			"rule_group_name": ruleGroupName,
			"exclusion":       flattenCdnFrontDoorFirewallExclusions(v.Exclusions),
			"rule":            flattenCdnFrontDoorFirewallRules(v.Rules),
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallRules(input *[]frontdoor.ManagedRuleOverride) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		action := "AnomalyScoring"
		if v.Action != "" {
			action = string(v.Action)
		}

		enabled := false
		if v.EnabledState != "" {
			enabled = v.EnabledState == frontdoor.ManagedRuleEnabledStateEnabled
		}

		ruleId := ""
		if v.RuleID != nil {
			ruleId = *v.RuleID
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
