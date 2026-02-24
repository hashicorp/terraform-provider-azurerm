// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
			_, err := waf.ParseFrontDoorWebApplicationFirewallPolicyID(id)
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

			// NOTE: 'js challenge expiration' is always
			// enabled no matter what and cannot be disabled for Premium_AzureFrontDoor
			// and is not supported in Standard_AzureFrontDoor...
			"js_challenge_cookie_expiration_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(5, 1440),
			},

			// NOTE: 'captcha expiration' is always
			// enabled no matter what and cannot be disabled for Premium_AzureFrontDoor
			// and is not supported in Standard_AzureFrontDoor...

			// NOTE: This field is Optional + Computed because:
			//  * Optional: Users can override the Azure default value (e.g., 30 minutes)
			//  * Computed: Azure automatically enables CAPTCHA policy with a default of 30 minutes on the Premium_AzureFrontDoor SKU,
			//    so the value is defined by Azure even when not explicitly set by the user
			"captcha_cookie_expiration_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(5, 1440),
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
								string(waf.ActionTypeJSChallenge),
								string(waf.ActionTypeCAPTCHA),
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

												// NOTE: 'ActionTypeAnomalyScoring' is only valid with 2.0 and above
												//       'ActionTypeJSChallenge' is only valid with BotManagerRuleSets
												"action": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice(waf.PossibleValuesForActionType(),
														false),
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

			"log_scrubbing": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"scrubbing_rule": {
							Type:     pluginsdk.TypeList,
							MaxItems: 100,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
									},

									"match_variable": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(waf.PossibleValuesForScrubbingRuleEntryMatchVariable(),
											false),
									},

									"operator": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(waf.ScrubbingRuleEntryMatchOperatorEquals),
										ValidateFunc: validation.StringInSlice(waf.PossibleValuesForScrubbingRuleEntryMatchOperator(),
											false),
									},

									"selector": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
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

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
				currentSku := diff.Get("sku_name").(string)
				standardSku := string(waf.SkuNameStandardAzureFrontDoor)
				oldSku, _ := diff.GetChange("sku_name")

				if currentSku == standardSku {
					premiumSku := string(waf.SkuNamePremiumAzureFrontDoor)
					managedRules := diff.Get("managed_rule").([]interface{})
					customRules := expandCdnFrontDoorFirewallCustomRules(diff.Get("custom_rule").([]interface{}))

					// Verify that they are not downgrading the service from Premium SKU -> Standard SKU...
					if oldSku != "" {
						if oldSku.(string) == premiumSku {
							return fmt.Errorf("downgrading from the %q sku to the %q sku is not supported, got %q", premiumSku, standardSku, currentSku)
						}
					}

					// Verify that the Standard SKU is not setting the JSChallenge or Captcha policy...
					if v := diff.Get("js_challenge_cookie_expiration_in_minutes").(int); v > 0 {
						return fmt.Errorf("'js_challenge_cookie_expiration_in_minutes' field is only supported with the %q sku, got %q", premiumSku, currentSku)
					}

					if v := diff.Get("captcha_cookie_expiration_in_minutes").(int); v > 0 {
						return fmt.Errorf("'captcha_cookie_expiration_in_minutes' field is only supported with the %q sku, got %q", premiumSku, currentSku)
					}

					// Verify that the Standard SKU is not using the JSChallenge or CAPTCHA Action type for custom rules...
					if customRules != nil && customRules.Rules != nil {
						for _, v := range *customRules.Rules {
							switch v.Action {
							case waf.ActionTypeJSChallenge:
								return fmt.Errorf("'custom_rule' blocks with the 'action' type of 'JSChallenge' are only supported for the %q sku, got action: %q (custom_rule.name: %q, sku_name: %q)", premiumSku, waf.ActionTypeJSChallenge, *v.Name, currentSku)
							case waf.ActionTypeCAPTCHA:
								return fmt.Errorf("'custom_rule' blocks with the 'action' type of 'CAPTCHA' are only supported for the %q sku, got action: %q (custom_rule.name: %q, sku_name: %q)", premiumSku, waf.ActionTypeCAPTCHA, *v.Name, currentSku)
							}
						}
					}

					// Verify that the Standard SKU is not using managed rules...
					if len(managedRules) > 0 {
						return fmt.Errorf("'managed_rule' code block is only supported with the %q sku, got %q", premiumSku, currentSku)
					}
				}

				return nil
			}),

			// Verify that the scrubbing_rule's are valid...
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
				if v, ok := diff.GetOk("log_scrubbing"); ok {
					_, err := expandCdnFrontDoorFirewallLogScrubbingPolicy(v.([]interface{}))
					if err != nil {
						return err
					}
				}
				return nil
			}),

			// Handle default value reset when field is removed from the configuration
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
				rawConfig := diff.GetRawConfig()

				if diff.Get("sku_name").(string) == string(waf.SkuNamePremiumAzureFrontDoor) {
					// Force the value to default when removed from config
					if rawConfig.IsNull() || rawConfig.GetAttr("js_challenge_cookie_expiration_in_minutes").IsNull() {
						if diff.Get("js_challenge_cookie_expiration_in_minutes").(int) != 30 {
							if err := diff.SetNew("js_challenge_cookie_expiration_in_minutes", 30); err != nil {
								return fmt.Errorf("setting default for `js_challenge_cookie_expiration_in_minutes`: %+v", err)
							}
						}
					}

					// Force the value to default when removed from config
					if rawConfig.IsNull() || rawConfig.GetAttr("captcha_cookie_expiration_in_minutes").IsNull() {
						if diff.Get("captcha_cookie_expiration_in_minutes").(int) != 30 {
							if err := diff.SetNew("captcha_cookie_expiration_in_minutes", 30); err != nil {
								return fmt.Errorf("setting default for `captcha_cookie_expiration_in_minutes`: %+v", err)
							}
						}
					}
				}

				return nil
			}),
		),
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
		return fmt.Errorf("expanding 'managed_rule': %+v", err)
	}

	logScrubbingRules, err := expandCdnFrontDoorFirewallLogScrubbingPolicy(d.Get("log_scrubbing").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding 'log_scrubbing': %+v", err)
	}

	if sku != string(waf.SkuNamePremiumAzureFrontDoor) && managedRules != nil {
		return fmt.Errorf("the 'managed_rule' field is only supported with the 'Premium_AzureFrontDoor' sku, got %q", sku)
	}

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
			},
			CustomRules: expandCdnFrontDoorFirewallCustomRules(customRules),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	// NOTE: CAPTCHA and JS Challenge Expiration policy is enabled by default on Premium SKU's with a default of
	// 30 minutes, if it is not in the config set the default and include it in the policy settings
	// payload block...
	if sku == string(waf.SkuNamePremiumAzureFrontDoor) {
		// Set the Default values...
		jsChallengeExpirationInMinutes := 30
		captchaExpirationInMinutes := 30

		if v, ok := d.GetOk("js_challenge_cookie_expiration_in_minutes"); ok {
			jsChallengeExpirationInMinutes = v.(int)
		}

		if v, ok := d.GetOk("captcha_cookie_expiration_in_minutes"); ok {
			captchaExpirationInMinutes = v.(int)
		}

		payload.Properties.PolicySettings.JavascriptChallengeExpirationInMinutes = pointer.To(int64(jsChallengeExpirationInMinutes))
		payload.Properties.PolicySettings.CaptchaExpirationInMinutes = pointer.To(int64(captchaExpirationInMinutes))
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

	if logScrubbingRules != nil {
		payload.Properties.PolicySettings.LogScrubbing = logScrubbingRules
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

	result, err := client.PoliciesGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := result.Model

	if model == nil {
		return fmt.Errorf("retrieving %s: 'model' was nil", *id)
	}

	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: 'model.Properties' was nil", *id)
	}

	props := *model.Properties

	if d.HasChanges("custom_block_response_body", "custom_block_response_status_code", "enabled", "mode", "redirect_url", "request_body_check_enabled", "js_challenge_cookie_expiration_in_minutes", "captcha_cookie_expiration_in_minutes", "log_scrubbing") {
		enabled := waf.PolicyEnabledStateDisabled
		if d.Get("enabled").(bool) {
			enabled = waf.PolicyEnabledStateEnabled
		}

		requestBodyCheck := waf.PolicyRequestBodyCheckDisabled
		if d.Get("request_body_check_enabled").(bool) {
			requestBodyCheck = waf.PolicyRequestBodyCheckEnabled
		}

		props.PolicySettings = &waf.PolicySettings{
			EnabledState:     pointer.To(enabled),
			Mode:             pointer.To(waf.PolicyMode(d.Get("mode").(string))),
			RequestBodyCheck: pointer.To(requestBodyCheck),
		}

		// NOTE: 'captcha_cookie_expiration_in_minutes' and 'js_challenge_cookie_expiration_in_minutes'
		// is only valid for 'Premium_AzureFrontDoor' skus...
		if model.Sku != nil && model.Sku.Name != nil && *model.Sku.Name == waf.SkuNamePremiumAzureFrontDoor {
			// Set the Default value...
			jsChallengeExpirationInMinutes := 30
			captchaExpirationInMinutes := 30

			if v, ok := d.GetOk("js_challenge_cookie_expiration_in_minutes"); ok {
				jsChallengeExpirationInMinutes = v.(int)
			}

			if v, ok := d.GetOk("captcha_cookie_expiration_in_minutes"); ok {
				captchaExpirationInMinutes = v.(int)
			}

			props.PolicySettings.JavascriptChallengeExpirationInMinutes = pointer.To(int64(jsChallengeExpirationInMinutes))
			props.PolicySettings.CaptchaExpirationInMinutes = pointer.To(int64(captchaExpirationInMinutes))
		}

		if redirectUrl := d.Get("redirect_url").(string); redirectUrl != "" {
			props.PolicySettings.RedirectURL = pointer.To(redirectUrl)
		}

		if body := d.Get("custom_block_response_body").(string); body != "" {
			props.PolicySettings.CustomBlockResponseBody = pointer.To(body)
		}

		if statusCode := int64(d.Get("custom_block_response_status_code").(int)); statusCode > 0 {
			props.PolicySettings.CustomBlockResponseStatusCode = pointer.To(statusCode)
		}
	}

	if d.HasChange("custom_rule") {
		props.CustomRules = expandCdnFrontDoorFirewallCustomRules(d.Get("custom_rule").([]interface{}))
	}

	if d.HasChange("managed_rule") {
		if model.Sku == nil {
			return fmt.Errorf("retrieving %s: 'model.Sku' was nil", *id)
		}

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

	if d.HasChange("log_scrubbing") {
		logScrubbingPolicy, err := expandCdnFrontDoorFirewallLogScrubbingPolicy(d.Get("log_scrubbing").([]interface{}))
		if err != nil {
			return err
		}

		if logScrubbingPolicy != nil {
			props.PolicySettings.LogScrubbing = logScrubbingPolicy
		}
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	model.Properties = pointer.To(props)

	err = client.PoliciesCreateOrUpdateThenPoll(ctx, *id, *model)
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

	resp, err := client.PoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.FrontDoorWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if sku := model.Sku; sku != nil {
			d.Set("sku_name", string(pointer.From(sku.Name)))
		}

		if props := model.Properties; props != nil {
			if err := d.Set("custom_rule", flattenCdnFrontDoorFirewallCustomRules(props.CustomRules)); err != nil {
				return fmt.Errorf("flattening 'custom_rule': %+v", err)
			}

			if err := d.Set("frontend_endpoint_ids", flattenFrontendEndpointLinkSlice(props.FrontendEndpointLinks)); err != nil {
				return fmt.Errorf("flattening 'frontend_endpoint_ids': %+v", err)
			}

			if err := d.Set("managed_rule", flattenCdnFrontDoorFirewallManagedRules(props.ManagedRules)); err != nil {
				return fmt.Errorf("flattening 'managed_rule': %+v", err)
			}

			if policy := props.PolicySettings; policy != nil {
				d.Set("enabled", pointer.From(policy.EnabledState) == waf.PolicyEnabledStateEnabled)
				d.Set("mode", string(pointer.From(policy.Mode)))
				d.Set("request_body_check_enabled", pointer.From(policy.RequestBodyCheck) == waf.PolicyRequestBodyCheckEnabled)
				d.Set("redirect_url", policy.RedirectURL)
				d.Set("custom_block_response_status_code", int(pointer.From(policy.CustomBlockResponseStatusCode)))
				d.Set("custom_block_response_body", policy.CustomBlockResponseBody)

				// NOTE: `js_challenge_cookie_expiration_in_minutes` and
				// `captcha_cookie_expiration_in_minutes` is only returned
				// for Premium_AzureFrontDoor skus, else they will be 'nil'...
				if policy.JavascriptChallengeExpirationInMinutes != nil {
					d.Set("js_challenge_cookie_expiration_in_minutes", int(pointer.From(policy.JavascriptChallengeExpirationInMinutes)))
				}

				if policy.CaptchaExpirationInMinutes != nil {
					d.Set("captcha_cookie_expiration_in_minutes", int(pointer.From(policy.CaptchaExpirationInMinutes)))
				}

				if err := d.Set("log_scrubbing", flattenCdnFrontDoorFirewallLogScrubbingPolicy(policy.LogScrubbing)); err != nil {
					return fmt.Errorf("flattening 'log_scrubbing': %+v", err)
				}
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceCdnFrontDoorFirewallPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorFirewallPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := waf.ParseFrontDoorWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	err = client.PoliciesDeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
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
		matchValues := utils.ExpandStringSlice(match["match_values"].([]interface{}))
		transforms := match["transforms"].([]interface{})

		matchCondition := waf.MatchCondition{
			Operator:        waf.Operator(operator),
			NegateCondition: &negateCondition,
			MatchValue:      *matchValues,
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

		ruleGroupOverrides, err := expandCdnFrontDoorFirewallManagedRuleGroupOverride(overrides, version, fVersion, ruleType)
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

func expandCdnFrontDoorFirewallManagedRuleGroupOverride(input []interface{}, versionRaw string, version float64, ruleType string) (*[]waf.ManagedRuleGroupOverride, error) {
	result := make([]waf.ManagedRuleGroupOverride, 0)
	if len(input) == 0 {
		return nil, nil
	}

	for _, v := range input {
		override := v.(map[string]interface{})

		exclusions := expandCdnFrontDoorFirewallManagedRuleGroupExclusion(override["exclusion"].([]interface{}))
		ruleGroupName := override["rule_group_name"].(string)
		rules, err := expandCdnFrontDoorFirewallRuleOverride(override["rule"].([]interface{}), versionRaw, version, ruleType)
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

func expandCdnFrontDoorFirewallRuleOverride(input []interface{}, versionRaw string, version float64, ruleType string) (*[]waf.ManagedRuleOverride, error) {
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
		action := waf.ActionType(rule["action"].(string))

		// NOTE: Default Rule Sets(DRS) 2.0 and above rules only use action type of 'AnomalyScoring' or 'Log'. Issues 19088 and 19561
		// This will still work for bot rules as well since it will be the default value of 1.0
		switch {
		case version < 2.0 && action == waf.ActionTypeAnomalyScoring:
			return nil, fmt.Errorf("%q is only valid in managed rules where 'type' is DRS and `version` is '2.0' or above, got %q", waf.ActionTypeAnomalyScoring, versionRaw)

		case version >= 2.0 && action != waf.ActionTypeAnomalyScoring && action != waf.ActionTypeLog:
			return nil, fmt.Errorf("the managed rules 'action' field must be set to 'AnomalyScoring' or 'Log' if the managed rule is DRS 2.0 or above, got %q", action)

		case !strings.Contains(strings.ToLower(ruleType), "botmanagerruleset") && action == waf.ActionTypeJSChallenge:
			return nil, fmt.Errorf("%q is only valid if the managed rules 'type' is 'Microsoft_BotManagerRuleSet', got %q", waf.ActionTypeJSChallenge, ruleType)
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

func expandCdnFrontDoorFirewallLogScrubbingPolicy(input []interface{}) (*waf.PolicySettingsLogScrubbing, error) {
	if len(input) == 0 {
		return nil, nil
	}

	inputRaw := input[0].(map[string]interface{})

	policyEnabled := waf.WebApplicationFirewallScrubbingStateDisabled
	if inputRaw["enabled"].(bool) {
		policyEnabled = waf.WebApplicationFirewallScrubbingStateEnabled
	}

	scrubbingRules, err := expandCdnFrontDoorFirewallScrubbingRules(inputRaw["scrubbing_rule"].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &waf.PolicySettingsLogScrubbing{
		State:          pointer.To(policyEnabled),
		ScrubbingRules: scrubbingRules,
	}, nil
}

func expandCdnFrontDoorFirewallScrubbingRules(input []interface{}) (*[]waf.WebApplicationFirewallScrubbingRules, error) {
	if len(input) == 0 {
		return nil, nil
	}

	scrubbingRules := make([]waf.WebApplicationFirewallScrubbingRules, 0)

	for _, rule := range input {
		v := rule.(map[string]interface{})
		var item waf.WebApplicationFirewallScrubbingRules

		enalbed := waf.ScrubbingRuleEntryStateDisabled
		if value := v["enabled"].(bool); value {
			enalbed = waf.ScrubbingRuleEntryStateEnabled
		}

		item.State = pointer.To(enalbed)
		item.MatchVariable = waf.ScrubbingRuleEntryMatchVariable(v["match_variable"].(string))
		item.SelectorMatchOperator = waf.ScrubbingRuleEntryMatchOperator(v["operator"].(string))

		if selector, ok := v["selector"]; ok {
			item.Selector = pointer.To(selector.(string))
		}

		// NOTE: Validate the rules configuration...
		switch {
		case item.MatchVariable == waf.ScrubbingRuleEntryMatchVariableRequestIPAddress || item.MatchVariable == waf.ScrubbingRuleEntryMatchVariableRequestUri:
			// NOTE: 'RequestIPAddress' and 'RequestUri' 'match_variable's can only use the 'EqualsAny' 'operator'...
			if item.SelectorMatchOperator != waf.ScrubbingRuleEntryMatchOperatorEqualsAny {
				return nil, fmt.Errorf("the %q 'match_variable' must use the %q 'operator', got %q", item.MatchVariable, waf.ScrubbingRuleEntryMatchOperatorEqualsAny, item.SelectorMatchOperator)
			}

		case item.SelectorMatchOperator == waf.ScrubbingRuleEntryMatchOperatorEquals:
			// NOTE: If the 'operator' is set to 'Equals' the 'selector' cannot be 'nil'...
			if pointer.From(item.Selector) == "" {
				return nil, fmt.Errorf("the 'selector' field must be set when the %q 'operator' is used, got %q", waf.ScrubbingRuleEntryMatchOperatorEquals, "nil")
			}

		case item.SelectorMatchOperator == waf.ScrubbingRuleEntryMatchOperatorEqualsAny:
			// NOTE: If the 'operator' is set to 'EqualsAny' the 'selector' must be 'nil'...
			if pointer.From(item.Selector) != "" {
				return nil, fmt.Errorf("the 'selector' field cannot be set when the %q 'operator' is used, got %q", waf.ScrubbingRuleEntryMatchOperatorEqualsAny, pointer.From(item.Selector))
			}
		}

		scrubbingRules = append(scrubbingRules, item)
	}

	return pointer.To(scrubbingRules), nil
}

func flattenCdnFrontDoorFirewallCustomRules(input *waf.CustomRuleList) []interface{} {
	if input == nil || input.Rules == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input.Rules {
		action := string(v.Action)
		priority := int(v.Priority)
		ruleType := string(v.RuleType)

		enabled := false
		if v.EnabledState != nil {
			enabled = pointer.From(v.EnabledState) == waf.CustomRuleEnabledStateEnabled
		}

		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		rateLimitDurationInMinutes := 0
		if v.RateLimitDurationInMinutes != nil {
			rateLimitDurationInMinutes = int(*v.RateLimitDurationInMinutes)
		}

		rateLimitThreshold := 0
		if v.RateLimitThreshold != nil {
			rateLimitThreshold = int(*v.RateLimitThreshold)
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

func flattenCdnFrontDoorFirewallMatchConditions(input []waf.MatchCondition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range input {
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
		ruleSetType := r.RuleSetType
		ruleSetVersion := r.RuleSetVersion
		ruleSetAction := string(pointer.From(r.RuleSetAction))

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
		matchVariable := string(v.MatchVariable)
		operator := string(v.SelectorMatchOperator)
		selector := v.Selector

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
		ruleGroupName := v.RuleGroupName

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
		enabled := pointer.From(v.EnabledState) == waf.ManagedRuleEnabledStateEnabled
		ruleId := v.RuleId

		results = append(results, map[string]interface{}{
			"action":    action,
			"enabled":   enabled,
			"exclusion": flattenCdnFrontDoorFirewallExclusions(v.Exclusions),
			"rule_id":   ruleId,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallLogScrubbingPolicy(input *waf.PolicySettingsLogScrubbing) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["enabled"] = pointer.From(input.State) == waf.WebApplicationFirewallScrubbingStateEnabled
	result["scrubbing_rule"] = flattenCdnFrontDoorFirewallLogScrubbingRules(input.ScrubbingRules)

	return []interface{}{result}
}

func flattenCdnFrontDoorFirewallLogScrubbingRules(scrubbingRules *[]waf.WebApplicationFirewallScrubbingRules) interface{} {
	result := make([]interface{}, 0)

	if scrubbingRules == nil || len(*scrubbingRules) == 0 {
		return result
	}

	for _, scrubbingRule := range *scrubbingRules {
		item := map[string]interface{}{}
		item["enabled"] = pointer.From(scrubbingRule.State) == waf.ScrubbingRuleEntryStateEnabled
		item["match_variable"] = scrubbingRule.MatchVariable
		item["operator"] = scrubbingRule.SelectorMatchOperator
		item["selector"] = pointer.From(scrubbingRule.Selector)

		result = append(result, item)
	}

	return result
}
