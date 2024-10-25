// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceWebApplicationFirewallPolicy() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceWebApplicationFirewallPolicyCreate,
		Read:   resourceWebApplicationFirewallPolicyRead,
		Update: resourceWebApplicationFirewallPolicyUpdate,
		Delete: resourceWebApplicationFirewallPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WebApplicationFirewallPolicyV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"custom_rules": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(webapplicationfirewallpolicies.WebApplicationFirewallActionAllow),
								string(webapplicationfirewallpolicies.WebApplicationFirewallActionBlock),
								string(webapplicationfirewallpolicies.WebApplicationFirewallActionLog),
							}, false),
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"match_conditions": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"match_values": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"match_variables": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"variable_name": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(webapplicationfirewallpolicies.WebApplicationFirewallMatchVariableRemoteAddr),
														string(webapplicationfirewallpolicies.WebApplicationFirewallMatchVariableRequestMethod),
														string(webapplicationfirewallpolicies.WebApplicationFirewallMatchVariableQueryString),
														string(webapplicationfirewallpolicies.WebApplicationFirewallMatchVariablePostArgs),
														string(webapplicationfirewallpolicies.WebApplicationFirewallMatchVariableRequestUri),
														string(webapplicationfirewallpolicies.WebApplicationFirewallMatchVariableRequestHeaders),
														string(webapplicationfirewallpolicies.WebApplicationFirewallMatchVariableRequestBody),
														string(webapplicationfirewallpolicies.WebApplicationFirewallMatchVariableRequestCookies),
													}, false),
												},
												"selector": {
													Type:     pluginsdk.TypeString,
													Optional: true,
												},
											},
										},
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorAny),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorIPMatch),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorGeoMatch),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorEqual),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorContains),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorLessThan),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorGreaterThan),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorLessThanOrEqual),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorGreaterThanOrEqual),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorBeginsWith),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorEndsWith),
											string(webapplicationfirewallpolicies.WebApplicationFirewallOperatorRegex),
										}, false),
									},
									"negation_condition": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},
									"transforms": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(webapplicationfirewallpolicies.WebApplicationFirewallTransformHtmlEntityDecode),
												string(webapplicationfirewallpolicies.WebApplicationFirewallTransformLowercase),
												string(webapplicationfirewallpolicies.WebApplicationFirewallTransformRemoveNulls),
												string(webapplicationfirewallpolicies.WebApplicationFirewallTransformTrim),
												string(webapplicationfirewallpolicies.WebApplicationFirewallTransformURLDecode),
												string(webapplicationfirewallpolicies.WebApplicationFirewallTransformURLEncode),
											}, false),
										},
									},
								},
							},
						},
						"priority": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
						"rule_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(webapplicationfirewallpolicies.WebApplicationFirewallRuleTypeMatchRule),
								string(webapplicationfirewallpolicies.WebApplicationFirewallRuleTypeRateLimitRule),
								string(webapplicationfirewallpolicies.WebApplicationFirewallRuleTypeInvalid),
							}, false),
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"rate_limit_duration": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForApplicationGatewayFirewallRateLimitDuration(), false),
						},
						"rate_limit_threshold": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"group_rate_limit_by": {
							// group variables combination not supported yet, use a single variable name
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForApplicationGatewayFirewallUserSessionVariable(), false),
						},
					},
				},
			},

			"managed_rules": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"exclusion": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"match_variable": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestArgKeys),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestArgNames),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestArgValues),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestCookieKeys),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestCookieNames),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestCookieValues),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestHeaderKeys),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestHeaderNames),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestHeaderValues),
										}, false),
									},
									"selector": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.NoZeroValues,
									},
									"selector_match_operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorContains),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorEndsWith),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorEquals),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorEqualsAny),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorStartsWith),
										}, false),
									},
									"excluded_rule_set": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"type": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													Default:      "OWASP",
													ValidateFunc: validate.ValidateWebApplicationFirewallPolicyExclusionRuleSetType,
												},
												"version": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													Default:      "3.2",
													ValidateFunc: validate.ValidateWebApplicationFirewallPolicyExclusionRuleSetVersion,
												},
												"rule_group": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"rule_group_name": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleGroupName,
															},
															"excluded_rules": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																Elem: &pluginsdk.Schema{
																	Type: pluginsdk.TypeString,
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
						"managed_rule_set": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"type": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "OWASP",
										ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleSetType,
									},
									"version": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleSetVersion,
									},
									"rule_group_override": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"rule_group_name": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleGroupName,
												},
												"rule": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													Computed: !features.FourPointOhBeta(),
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"id": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},

															"enabled": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																Default: func() interface{} {
																	if !features.FourPointOhBeta() {
																		return nil
																	}

																	return false
																}(),
															},

															"action": {
																Type:     pluginsdk.TypeString,
																Optional: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(webapplicationfirewallpolicies.ActionTypeAllow),
																	string(webapplicationfirewallpolicies.ActionTypeAnomalyScoring),
																	string(webapplicationfirewallpolicies.ActionTypeBlock),
																	string(webapplicationfirewallpolicies.ActionTypeJSChallenge),
																	string(webapplicationfirewallpolicies.ActionTypeLog),
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
					},
				},
			},

			"policy_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"mode": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(webapplicationfirewallpolicies.WebApplicationFirewallModePrevention),
								string(webapplicationfirewallpolicies.WebApplicationFirewallModeDetection),
							}, false),
							Default: string(webapplicationfirewallpolicies.WebApplicationFirewallModePrevention),
						},

						"request_body_check": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"file_upload_limit_in_mb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4000),
							Default:      100,
						},

						"request_body_enforcement": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"max_request_body_size_in_kb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(8, 2000),
							Default:      128,
						},

						"request_body_inspect_limit_in_kb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      128,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"js_challenge_cookie_expiration_in_minutes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      30,
							ValidateFunc: validation.IntBetween(5, 1440),
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

									"rule": {
										Type:     pluginsdk.TypeList,
										Optional: true,
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
													ValidateFunc: validation.StringInSlice(
														webapplicationfirewallpolicies.PossibleValuesForScrubbingRuleEntryMatchVariable(),
														false),
												},

												"selector_match_operator": {
													Type:     pluginsdk.TypeString,
													Optional: true,
													Default:  "Equals",
													ValidateFunc: validation.StringInSlice(
														webapplicationfirewallpolicies.PossibleValuesForScrubbingRuleEntryMatchOperator(),
														false),
												},

												"selector": {
													Type:        pluginsdk.TypeString,
													Optional:    true,
													Description: "When matchVariable is a collection, operator used to specify which elements in the collection this rule applies to.",
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

			"http_listener_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"path_based_rule_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"tags": commonschema.Tags(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["managed_rules"].Elem.(*pluginsdk.Resource).Schema["managed_rule_set"].Elem.(*pluginsdk.Resource).Schema["rule_group_override"].Elem.(*pluginsdk.Resource).Schema["disabled_rules"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			Deprecated: "`disabled_rules` will be removed in favour of the `rule` property in version 4.0 of the AzureRM Provider.",
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		}
	}

	return resource
}

func resourceWebApplicationFirewallPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WebApplicationFirewallPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webapplicationfirewallpolicies.NewApplicationGatewayWebApplicationFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_web_application_firewall_policy", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	customRules := d.Get("custom_rules").([]interface{})
	policySettings := d.Get("policy_settings").([]interface{})
	managedRules := d.Get("managed_rules").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	expandedManagedRules, err := expandWebApplicationFirewallPolicyManagedRulesDefinition(managedRules, d)
	if err != nil {
		return err
	}

	parameters := webapplicationfirewallpolicies.WebApplicationFirewallPolicy{
		Location: utils.String(location),
		Properties: &webapplicationfirewallpolicies.WebApplicationFirewallPolicyPropertiesFormat{
			CustomRules:    expandWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(customRules),
			PolicySettings: expandWebApplicationFirewallPolicyPolicySettings(policySettings),
			ManagedRules:   pointer.From(expandedManagedRules),
		},
		Tags: tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceWebApplicationFirewallPolicyRead(d, meta)
}

func resourceWebApplicationFirewallPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WebApplicationFirewallPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webapplicationfirewallpolicies.NewApplicationGatewayWebApplicationFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}

	model := resp.Model

	if d.HasChange("custom_rules") {
		model.Properties.CustomRules = expandWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(d.Get("custom_rules").([]interface{}))
	}

	if d.HasChange("policy_settings") {
		model.Properties.PolicySettings = expandWebApplicationFirewallPolicyPolicySettings(d.Get("policy_settings").([]interface{}))
	}

	if d.HasChange("managed_rules") {
		expandedManagedRules, err := expandWebApplicationFirewallPolicyManagedRulesDefinition(d.Get("managed_rules").([]interface{}), d)
		if err != nil {
			return err
		}
		model.Properties.ManagedRules = pointer.From(expandedManagedRules)
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, id, *model); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	return resourceWebApplicationFirewallPolicyRead(d, meta)
}

func resourceWebApplicationFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WebApplicationFirewallPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Web Application Firewall Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.ApplicationGatewayWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if prop := model.Properties; prop != nil {
			if err := d.Set("custom_rules", flattenWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(prop.CustomRules)); err != nil {
				return fmt.Errorf("setting `custom_rules`: %+v", err)
			}
			if err := d.Set("policy_settings", flattenWebApplicationFirewallPolicyPolicySettings(prop.PolicySettings)); err != nil {
				return fmt.Errorf("setting `policy_settings`: %+v", err)
			}
			if err := d.Set("managed_rules", flattenWebApplicationFirewallPolicyManagedRulesDefinition(prop.ManagedRules)); err != nil {
				return fmt.Errorf("setting `managed_rules`: %+v", err)
			}
			if err := d.Set("http_listener_ids", flattenWebApplicationFirewallPoliciesSubResourcesToIDs(prop.HTTPListeners)); err != nil {
				return fmt.Errorf("setting `http_listeners`: %+v", err)
			}
			if err := d.Set("path_based_rule_ids", flattenWebApplicationFirewallPoliciesSubResourcesToIDs(prop.PathBasedRules)); err != nil {
				return fmt.Errorf("setting `path_based_rules`: %+v", err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceWebApplicationFirewallPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WebApplicationFirewallPolicies
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(input []interface{}) *[]webapplicationfirewallpolicies.WebApplicationFirewallCustomRule {
	results := make([]webapplicationfirewallpolicies.WebApplicationFirewallCustomRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		name := v["name"].(string)
		priority := v["priority"].(int)
		ruleType := v["rule_type"].(string)
		matchConditions := v["match_conditions"].([]interface{})
		action := v["action"].(string)

		enabled := webapplicationfirewallpolicies.WebApplicationFirewallStateEnabled
		if value, ok := v["enabled"].(bool); ok && !value {
			enabled = webapplicationfirewallpolicies.WebApplicationFirewallStateDisabled
		}

		result := webapplicationfirewallpolicies.WebApplicationFirewallCustomRule{
			State:           pointer.To(enabled),
			Action:          webapplicationfirewallpolicies.WebApplicationFirewallAction(action),
			MatchConditions: expandWebApplicationFirewallPolicyMatchCondition(matchConditions),
			Name:            pointer.To(name),
			Priority:        int64(priority),
			RuleType:        webapplicationfirewallpolicies.WebApplicationFirewallRuleType(ruleType),
		}

		if rateLimitDuration, ok := v["rate_limit_duration"]; ok && rateLimitDuration.(string) != "" {
			result.RateLimitDuration = pointer.To(webapplicationfirewallpolicies.ApplicationGatewayFirewallRateLimitDuration(rateLimitDuration.(string)))
		}

		if rateLimitThreshHold, ok := v["rate_limit_threshold"]; ok && rateLimitThreshHold.(int) > 0 {
			result.RateLimitThreshold = pointer.To(int64(rateLimitThreshHold.(int)))
		}

		if groupBy, ok := v["group_rate_limit_by"]; ok && groupBy.(string) != "" {
			groups := []webapplicationfirewallpolicies.GroupByUserSession{
				{
					GroupByVariables: []webapplicationfirewallpolicies.GroupByVariable{
						{
							VariableName: webapplicationfirewallpolicies.ApplicationGatewayFirewallUserSessionVariable(groupBy.(string)),
						},
					},
				},
			}
			result.GroupByUserSession = &groups
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyPolicySettings(input []interface{}) *webapplicationfirewallpolicies.PolicySettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	enabled := webapplicationfirewallpolicies.WebApplicationFirewallEnabledStateDisabled
	if value, ok := v["enabled"].(bool); ok && value {
		enabled = webapplicationfirewallpolicies.WebApplicationFirewallEnabledStateEnabled
	}
	mode := v["mode"].(string)
	requestBodyCheck := v["request_body_check"].(bool)
	requestBodyEnforcement := v["request_body_enforcement"].(bool)
	maxRequestBodySizeInKb := v["max_request_body_size_in_kb"].(int)
	fileUploadLimitInMb := v["file_upload_limit_in_mb"].(int)

	result := webapplicationfirewallpolicies.PolicySettings{
		State:                             pointer.To(enabled),
		Mode:                              pointer.To(webapplicationfirewallpolicies.WebApplicationFirewallMode(mode)),
		RequestBodyCheck:                  pointer.To(requestBodyCheck),
		RequestBodyEnforcement:            pointer.To(requestBodyEnforcement),
		MaxRequestBodySizeInKb:            pointer.To(int64(maxRequestBodySizeInKb)),
		FileUploadLimitInMb:               pointer.To(int64(fileUploadLimitInMb)),
		LogScrubbing:                      expandWebApplicationFirewallPolicyLogScrubbing(v["log_scrubbing"].([]interface{})),
		RequestBodyInspectLimitInKB:       pointer.To(int64(v["request_body_inspect_limit_in_kb"].(int))),
		JsChallengeCookieExpirationInMins: pointer.To(int64(v["js_challenge_cookie_expiration_in_minutes"].(int))),
	}

	return &result
}

func expandWebApplicationFirewallPolicyLogScrubbing(input []interface{}) *webapplicationfirewallpolicies.PolicySettingsLogScrubbing {
	if len(input) == 0 {
		return nil
	}

	var res webapplicationfirewallpolicies.PolicySettingsLogScrubbing
	v := input[0].(map[string]interface{})
	state := webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateDisabled
	if value, ok := v["enabled"].(bool); ok && value {
		state = webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateEnabled
	}
	res.State = &state

	res.ScrubbingRules = expanedWebApplicationPolicyScrubbingRules(v["rule"].([]interface{}))

	return &res
}

func expanedWebApplicationPolicyScrubbingRules(input []interface{}) *[]webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules {
	if len(input) == 0 {
		return nil
	}
	var res []webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules
	for _, rule := range input {
		v := rule.(map[string]interface{})
		var item webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules
		state := webapplicationfirewallpolicies.ScrubbingRuleEntryStateDisabled
		if value, ok := v["enabled"].(bool); ok && value {
			state = webapplicationfirewallpolicies.ScrubbingRuleEntryStateEnabled
		}
		item.State = &state
		item.MatchVariable = webapplicationfirewallpolicies.ScrubbingRuleEntryMatchVariable(v["match_variable"].(string))
		item.SelectorMatchOperator = webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperator(v["selector_match_operator"].(string))
		if val, ok := v["selector"]; ok {
			item.Selector = pointer.To(val.(string))
		}

		res = append(res, item)
	}
	return &res
}

func expandWebApplicationFirewallPolicyManagedRulesDefinition(input []interface{}, d *pluginsdk.ResourceData) (*webapplicationfirewallpolicies.ManagedRulesDefinition, error) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]interface{})

	exclusions := v["exclusion"].([]interface{})
	managedRuleSets := v["managed_rule_set"].([]interface{})

	expandedManagedRuleSets, err := expandWebApplicationFirewallPolicyManagedRuleSet(managedRuleSets, d)
	if err != nil {
		return nil, err
	}

	return &webapplicationfirewallpolicies.ManagedRulesDefinition{
		Exclusions:      expandWebApplicationFirewallPolicyExclusions(exclusions),
		ManagedRuleSets: *expandedManagedRuleSets,
	}, nil
}

func expandWebApplicationFirewallPolicyExclusionManagedRules(input []interface{}) *[]webapplicationfirewallpolicies.ExclusionManagedRule {
	results := make([]webapplicationfirewallpolicies.ExclusionManagedRule, 0)
	for _, item := range input {
		ruleID := item.(string)

		result := webapplicationfirewallpolicies.ExclusionManagedRule{
			RuleId: ruleID,
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusionManagedRuleGroup(input []interface{}) *[]webapplicationfirewallpolicies.ExclusionManagedRuleGroup {
	results := make([]webapplicationfirewallpolicies.ExclusionManagedRuleGroup, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		ruleGroupName := v["rule_group_name"].(string)

		result := webapplicationfirewallpolicies.ExclusionManagedRuleGroup{
			RuleGroupName: ruleGroupName,
		}

		if excludedRules := v["excluded_rules"].([]interface{}); len(excludedRules) > 0 {
			result.Rules = expandWebApplicationFirewallPolicyExclusionManagedRules(excludedRules)
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusionManagedRuleSet(input []interface{}) *[]webapplicationfirewallpolicies.ExclusionManagedRuleSet {
	results := make([]webapplicationfirewallpolicies.ExclusionManagedRuleSet, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		ruleSetType := v["type"].(string)
		ruleSetVersion := v["version"].(string)
		ruleGroups := make([]interface{}, 0)
		if value, exists := v["rule_group"]; exists {
			ruleGroups = value.([]interface{})
		}
		result := webapplicationfirewallpolicies.ExclusionManagedRuleSet{
			RuleSetType:    ruleSetType,
			RuleSetVersion: ruleSetVersion,
			RuleGroups:     expandWebApplicationFirewallPolicyExclusionManagedRuleGroup(ruleGroups),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusions(input []interface{}) *[]webapplicationfirewallpolicies.OwaspCrsExclusionEntry {
	results := make([]webapplicationfirewallpolicies.OwaspCrsExclusionEntry, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		matchVariable := v["match_variable"].(string)
		selectorMatchOperator := v["selector_match_operator"].(string)
		selector := v["selector"].(string)
		exclusionManagedRuleSets := v["excluded_rule_set"].([]interface{})

		result := webapplicationfirewallpolicies.OwaspCrsExclusionEntry{
			MatchVariable:            webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariable(matchVariable),
			SelectorMatchOperator:    webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperator(selectorMatchOperator),
			Selector:                 selector,
			ExclusionManagedRuleSets: expandWebApplicationFirewallPolicyExclusionManagedRuleSet(exclusionManagedRuleSets),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyManagedRuleSet(input []interface{}, d *pluginsdk.ResourceData) (*[]webapplicationfirewallpolicies.ManagedRuleSet, error) {
	results := make([]webapplicationfirewallpolicies.ManagedRuleSet, 0)

	for i, item := range input {
		v := item.(map[string]interface{})
		ruleSetType := v["type"].(string)
		ruleSetVersion := v["version"].(string)
		ruleGroupOverrides := []interface{}{}
		if value, exists := v["rule_group_override"]; exists {
			ruleGroupOverrides = value.([]interface{})
		}

		expandedRuleGroupOverrides, err := expandWebApplicationFirewallPolicyRuleGroupOverrides(ruleGroupOverrides, d, i)
		if err != nil {
			return nil, err
		}

		result := webapplicationfirewallpolicies.ManagedRuleSet{
			RuleSetType:        ruleSetType,
			RuleSetVersion:     ruleSetVersion,
			RuleGroupOverrides: expandedRuleGroupOverrides,
		}

		results = append(results, result)
	}
	return &results, nil
}

func expandWebApplicationFirewallPolicyRuleGroupOverrides(input []interface{}, d *pluginsdk.ResourceData, managedRuleSetIndex int) (*[]webapplicationfirewallpolicies.ManagedRuleGroupOverride, error) {
	results := make([]webapplicationfirewallpolicies.ManagedRuleGroupOverride, 0)
	for i, item := range input {
		v := item.(map[string]interface{})

		ruleGroupName := v["rule_group_name"].(string)

		result := webapplicationfirewallpolicies.ManagedRuleGroupOverride{
			RuleGroupName: ruleGroupName,
		}

		if !features.FourPointOhBeta() {
			// `disabled_rules` will be deprecated from 4.0. In 3.x, `rule` and `disabled_rules` point to the same properties of Azure REST API and conflict with each other in the configuration.
			// Since both properties will be set in the flatten method, need to use `GetRawConfig` to check which one of these two properties is configured in the configuration file.
			managedRuleSetList := d.GetRawConfig().AsValueMap()["managed_rules"].AsValueSlice()[0].AsValueMap()["managed_rule_set"].AsValueSlice()
			if managedRuleSetIndex >= len(managedRuleSetList) {
				return nil, fmt.Errorf("managed rule set index %d exceeds raw config length %d", managedRuleSetIndex, len(managedRuleSetList))
			}

			ruleGroupOverrideList := managedRuleSetList[managedRuleSetIndex].AsValueMap()["rule_group_override"].AsValueSlice()
			if i >= len(ruleGroupOverrideList) {
				return nil, fmt.Errorf("rule group override index %d exceeds raw config length %d", i, len(ruleGroupOverrideList))
			}

			// Since ConflictsWith cannot be used on these properties and the properties are optional and computed, Have to check the configuration with GetRawConfig
			if !ruleGroupOverrideList[i].AsValueMap()["rule"].IsNull() && len(ruleGroupOverrideList[i].AsValueMap()["rule"].AsValueSlice()) > 0 && !ruleGroupOverrideList[i].AsValueMap()["disabled_rules"].IsNull() {
				return nil, fmt.Errorf("`disabled_rules` cannot be set when `rule` is set under `rule_group_override`")
			}

			if disabledRules := v["disabled_rules"].([]interface{}); !ruleGroupOverrideList[i].AsValueMap()["disabled_rules"].IsNull() {
				result.Rules = expandWebApplicationFirewallPolicyRules(disabledRules)
			}

			if rules := v["rule"].([]interface{}); !ruleGroupOverrideList[i].AsValueMap()["rule"].IsNull() && len(ruleGroupOverrideList[i].AsValueMap()["rule"].AsValueSlice()) > 0 {
				result.Rules = expandWebApplicationFirewallPolicyOverrideRules(rules)
			}
		} else {
			if rules := v["rule"].([]interface{}); len(rules) > 0 {
				result.Rules = expandWebApplicationFirewallPolicyOverrideRules(rules)
			}
		}

		results = append(results, result)
	}

	return &results, nil
}

func expandWebApplicationFirewallPolicyRules(input []interface{}) *[]webapplicationfirewallpolicies.ManagedRuleOverride {
	results := make([]webapplicationfirewallpolicies.ManagedRuleOverride, 0)
	for _, item := range input {
		ruleID := item.(string)

		result := webapplicationfirewallpolicies.ManagedRuleOverride{
			RuleId: ruleID,
			State:  pointer.To(webapplicationfirewallpolicies.ManagedRuleEnabledStateDisabled),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyOverrideRules(input []interface{}) *[]webapplicationfirewallpolicies.ManagedRuleOverride {
	results := make([]webapplicationfirewallpolicies.ManagedRuleOverride, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		state := webapplicationfirewallpolicies.ManagedRuleEnabledStateDisabled
		if v["enabled"].(bool) {
			state = webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled
		}

		result := webapplicationfirewallpolicies.ManagedRuleOverride{
			RuleId: v["id"].(string),
			State:  pointer.To(state),
		}

		action := v["action"].(string)
		if action != "" {
			result.Action = pointer.To(webapplicationfirewallpolicies.ActionType(action))
		}

		results = append(results, result)
	}

	return &results
}

func expandWebApplicationFirewallPolicyMatchCondition(input []interface{}) []webapplicationfirewallpolicies.MatchCondition {
	results := make([]webapplicationfirewallpolicies.MatchCondition, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		matchVariables := v["match_variables"].([]interface{})
		operator := v["operator"].(string)
		negationCondition := v["negation_condition"].(bool)
		matchValues := v["match_values"].([]interface{})
		transformsRaw := v["transforms"].(*pluginsdk.Set).List()

		var transforms []webapplicationfirewallpolicies.WebApplicationFirewallTransform
		for _, trans := range transformsRaw {
			transforms = append(transforms, webapplicationfirewallpolicies.WebApplicationFirewallTransform(trans.(string)))
		}
		result := webapplicationfirewallpolicies.MatchCondition{
			MatchValues:      pointer.From(utils.ExpandStringSlice(matchValues)),
			MatchVariables:   expandWebApplicationFirewallPolicyMatchVariable(matchVariables),
			NegationConditon: utils.Bool(negationCondition),
			Operator:         webapplicationfirewallpolicies.WebApplicationFirewallOperator(operator),
			Transforms:       &transforms,
		}

		results = append(results, result)
	}
	return results
}

func expandWebApplicationFirewallPolicyMatchVariable(input []interface{}) []webapplicationfirewallpolicies.MatchVariable {
	results := make([]webapplicationfirewallpolicies.MatchVariable, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		variableName := v["variable_name"].(string)
		selector := v["selector"].(string)

		result := webapplicationfirewallpolicies.MatchVariable{
			Selector:     utils.String(selector),
			VariableName: webapplicationfirewallpolicies.WebApplicationFirewallMatchVariable(variableName),
		}

		results = append(results, result)
	}
	return results
}

func flattenWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(input *[]webapplicationfirewallpolicies.WebApplicationFirewallCustomRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if name := item.Name; name != nil {
			v["name"] = *name
		}
		v["action"] = string(item.Action)
		v["enabled"] = pointer.From(item.State) == webapplicationfirewallpolicies.WebApplicationFirewallStateEnabled
		v["match_conditions"] = flattenWebApplicationFirewallPolicyMatchCondition(item.MatchConditions)
		v["priority"] = int(item.Priority)
		v["rule_type"] = string(item.RuleType)
		v["rate_limit_duration"] = pointer.From(item.RateLimitDuration)
		v["rate_limit_threshold"] = pointer.From(item.RateLimitThreshold)

		if item.GroupByUserSession != nil && len(*item.GroupByUserSession) > 0 {
			if groupVariable := (*item.GroupByUserSession)[0].GroupByVariables; len(groupVariable) > 0 {
				v["group_rate_limit_by"] = groupVariable[0].VariableName
			}
		}

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPolicyPolicySettings(input *webapplicationfirewallpolicies.PolicySettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["enabled"] = pointer.From(input.State) == webapplicationfirewallpolicies.WebApplicationFirewallEnabledStateEnabled
	result["mode"] = string(pointer.From(input.Mode))
	result["request_body_check"] = input.RequestBodyCheck
	result["request_body_enforcement"] = input.RequestBodyEnforcement
	result["max_request_body_size_in_kb"] = int(pointer.From(input.MaxRequestBodySizeInKb))
	result["file_upload_limit_in_mb"] = int(pointer.From(input.FileUploadLimitInMb))
	result["log_scrubbing"] = flattenWebApplicationFirewallPolicyLogScrubbing(input.LogScrubbing)
	result["request_body_inspect_limit_in_kb"] = pointer.From(input.RequestBodyInspectLimitInKB)
	result["js_challenge_cookie_expiration_in_minutes"] = pointer.From(input.JsChallengeCookieExpirationInMins)

	return []interface{}{result}
}

func flattenWebApplicationFirewallPolicyLogScrubbing(input *webapplicationfirewallpolicies.PolicySettingsLogScrubbing) interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	result := make(map[string]interface{})
	result["enabled"] = pointer.From(input.State) == webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateEnabled
	result["rule"] = flattenWebApplicationFirewallPolicyLogScrubbingRules(input.ScrubbingRules)
	return []interface{}{result}
}

func flattenWebApplicationFirewallPolicyLogScrubbingRules(rules *[]webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules) interface{} {
	result := make([]interface{}, 0)
	if rules == nil || len(*rules) == 0 {
		return result
	}
	for _, rule := range *rules {
		item := map[string]interface{}{}
		item["enabled"] = pointer.From(rule.State) == webapplicationfirewallpolicies.ScrubbingRuleEntryStateEnabled
		item["match_variable"] = rule.MatchVariable
		item["selector_match_operator"] = rule.SelectorMatchOperator
		item["selector"] = pointer.From(rule.Selector)
		result = append(result, item)
	}
	return &result

}

func flattenWebApplicationFirewallPolicyManagedRulesDefinition(input webapplicationfirewallpolicies.ManagedRulesDefinition) []interface{} {
	results := make([]interface{}, 0)

	v := make(map[string]interface{})

	v["exclusion"] = flattenWebApplicationFirewallPolicyExclusions(input.Exclusions)
	v["managed_rule_set"] = flattenWebApplicationFirewallPolicyManagedRuleSets(input.ManagedRuleSets)

	results = append(results, v)

	return results
}

func flattenWebApplicationFirewallPolicyExclusionManagedRules(input *[]webapplicationfirewallpolicies.ExclusionManagedRule) []string {
	results := make([]string, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, item.RuleId)
	}

	return results
}

func flattenWebApplicationFirewallPolicyExclusionManagedRuleGroups(input *[]webapplicationfirewallpolicies.ExclusionManagedRuleGroup) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		v["rule_group_name"] = item.RuleGroupName
		v["excluded_rules"] = flattenWebApplicationFirewallPolicyExclusionManagedRules(item.Rules)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyExclusionManagedRuleSets(input *[]webapplicationfirewallpolicies.ExclusionManagedRuleSet) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		v["type"] = item.RuleSetType
		v["version"] = item.RuleSetVersion
		v["rule_group"] = flattenWebApplicationFirewallPolicyExclusionManagedRuleGroups(item.RuleGroups)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyExclusions(input *[]webapplicationfirewallpolicies.OwaspCrsExclusionEntry) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		selector := item.Selector

		v["match_variable"] = string(item.MatchVariable)
		v["selector"] = selector

		v["selector_match_operator"] = string(item.SelectorMatchOperator)
		v["excluded_rule_set"] = flattenWebApplicationFirewallPolicyExclusionManagedRuleSets(item.ExclusionManagedRuleSets)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyManagedRuleSets(input []webapplicationfirewallpolicies.ManagedRuleSet) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		v := make(map[string]interface{})

		v["type"] = item.RuleSetType
		v["version"] = item.RuleSetVersion
		v["rule_group_override"] = flattenWebApplicationFirewallPolicyRuleGroupOverrides(item.RuleGroupOverrides)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyRuleGroupOverrides(input *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		v["rule_group_name"] = item.RuleGroupName

		if !features.FourPointOhBeta() {
			v["disabled_rules"] = flattenWebApplicationFirewallPolicyRules(item.Rules)
		}

		v["rule"] = flattenWebApplicationFirewallPolicyOverrideRules(item.Rules)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyRules(input *[]webapplicationfirewallpolicies.ManagedRuleOverride) []string {
	results := make([]string, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		if item.State == nil || *item.State == webapplicationfirewallpolicies.ManagedRuleEnabledStateDisabled {
			results = append(results, item.RuleId)
		}
	}

	return results
}

func flattenWebApplicationFirewallPolicyOverrideRules(input *[]webapplicationfirewallpolicies.ManagedRuleOverride) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})
		v["id"] = item.RuleId

		v["enabled"] = pointer.From(item.State) == webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled

		v["action"] = string(pointer.From(item.Action))

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPolicyMatchCondition(input []webapplicationfirewallpolicies.MatchCondition) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		v := make(map[string]interface{})

		var transforms []interface{}
		if item.Transforms != nil {
			for _, trans := range *item.Transforms {
				transforms = append(transforms, string(trans))
			}
		}
		v["match_values"] = utils.FlattenStringSlice(pointer.To(item.MatchValues))
		v["match_variables"] = flattenWebApplicationFirewallPolicyMatchVariable(item.MatchVariables)
		if negationCondition := item.NegationConditon; negationCondition != nil {
			v["negation_condition"] = *negationCondition
		}
		v["operator"] = string(item.Operator)
		v["transforms"] = transforms

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPolicyMatchVariable(input []webapplicationfirewallpolicies.MatchVariable) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		v := make(map[string]interface{})

		if selector := item.Selector; selector != nil {
			v["selector"] = *selector
		}
		v["variable_name"] = string(item.VariableName)

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPoliciesSubResourcesToIDs(input *[]webapplicationfirewallpolicies.SubResource) []interface{} {
	ids := make([]interface{}, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}
