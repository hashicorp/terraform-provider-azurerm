package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorProfilePolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorProfilePolicyCreateUpdate,
		Read:   resourceFrontdoorProfilePolicyRead,
		Update: resourceFrontdoorProfilePolicyUpdate,
		Delete: resourceFrontdoorProfilePolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webapplicationfirewallpolicies.ParseCdnWebApplicationFirewallPoliciesID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"custom_rules": {
				Type:     pluginsdk.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"rules": {
							Type:     pluginsdk.TypeList,
							ForceNew: true,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"action": {
										Type:     pluginsdk.TypeString,
										ForceNew: true,
										Required: true,
									},

									"enabled_state": {
										Type:     pluginsdk.TypeString,
										ForceNew: true,
										Optional: true,
									},

									"match_conditions": {
										Type:     pluginsdk.TypeList,
										ForceNew: true,
										Required: true,

										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{

												"match_value": {
													Type:     pluginsdk.TypeList,
													ForceNew: true,
													Required: true,

													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},

												"match_variable": {
													Type:     pluginsdk.TypeString,
													ForceNew: true,
													Required: true,
												},

												"negate_condition": {
													Type:     pluginsdk.TypeBool,
													ForceNew: true,
													Required: true,
												},

												"operator": {
													Type:     pluginsdk.TypeString,
													ForceNew: true,
													Required: true,
												},

												"selector": {
													Type:     pluginsdk.TypeString,
													ForceNew: true,
													Required: true,
												},

												"transforms": {
													Type:     pluginsdk.TypeList,
													ForceNew: true,
													Required: true,

													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{},
													},
												},
											},
										},
									},

									"name": {
										Type:     pluginsdk.TypeString,
										ForceNew: true,
										Required: true,
									},

									"priority": {
										Type:     pluginsdk.TypeInt,
										ForceNew: true,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"endpoint_links": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"etag": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"location": azure.SchemaLocation(),

			"managed_rules": {
				Type:     pluginsdk.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"managed_rule_sets": {
							Type:     pluginsdk.TypeList,
							ForceNew: true,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"anomaly_score": {
										Type:     pluginsdk.TypeInt,
										ForceNew: true,
										Optional: true,
									},

									"rule_group_overrides": {
										Type:     pluginsdk.TypeList,
										ForceNew: true,
										Optional: true,

										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{

												"rule_group_name": {
													Type:     pluginsdk.TypeString,
													ForceNew: true,
													Required: true,
												},

												"rules": {
													Type:     pluginsdk.TypeList,
													ForceNew: true,
													Optional: true,

													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{

															"action": {
																Type:     pluginsdk.TypeString,
																ForceNew: true,
																Optional: true,
															},

															"enabled_state": {
																Type:     pluginsdk.TypeString,
																ForceNew: true,
																Optional: true,
															},

															"rule_id": {
																Type:     pluginsdk.TypeString,
																ForceNew: true,
																Required: true,
															},
														},
													},
												},
											},
										},
									},

									"rule_set_type": {
										Type:     pluginsdk.TypeString,
										ForceNew: true,
										Required: true,
									},

									"rule_set_version": {
										Type:     pluginsdk.TypeString,
										ForceNew: true,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"policy_settings": {
				Type:     pluginsdk.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"default_custom_block_response_body": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},

						"default_custom_block_response_status_code": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},

						"default_redirect_url": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},

						"enabled_state": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},

						"mode": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"rate_limit_rules": {
				Type:     pluginsdk.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"rules": {
							Type:     pluginsdk.TypeList,
							ForceNew: true,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"action": {
										Type:     pluginsdk.TypeString,
										ForceNew: true,
										Required: true,
									},

									"enabled_state": {
										Type:     pluginsdk.TypeString,
										ForceNew: true,
										Optional: true,
									},

									"match_conditions": {
										Type:     pluginsdk.TypeList,
										ForceNew: true,
										Required: true,

										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{

												"match_value": {
													Type:     pluginsdk.TypeList,
													ForceNew: true,
													Required: true,

													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},

												"match_variable": {
													Type:     pluginsdk.TypeString,
													ForceNew: true,
													Required: true,
												},

												"negate_condition": {
													Type:     pluginsdk.TypeBool,
													ForceNew: true,
													Required: true,
												},

												"operator": {
													Type:     pluginsdk.TypeString,
													ForceNew: true,
													Required: true,
												},

												"selector": {
													Type:     pluginsdk.TypeString,
													ForceNew: true,
													Required: true,
												},

												"transforms": {
													Type:     pluginsdk.TypeList,
													ForceNew: true,
													Required: true,

													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},
											},
										},
									},

									"name": {
										Type:     pluginsdk.TypeString,
										ForceNew: true,
										Required: true,
									},

									"priority": {
										Type:     pluginsdk.TypeInt,
										ForceNew: true,
										Required: true,
									},

									"rate_limit_duration_in_minutes": {
										Type:     pluginsdk.TypeInt,
										ForceNew: true,
										Required: true,
									},

									"rate_limit_threshold": {
										Type:     pluginsdk.TypeInt,
										ForceNew: true,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"resource_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				ForceNew: true,
				Required: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"name": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Required: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceFrontdoorProfilePolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.WebApplicationFirewallPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webapplicationfirewallpolicies.NewCdnWebApplicationFirewallPoliciesID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.PoliciesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_policy", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location"))
	props := webapplicationfirewallpolicies.CdnWebApplicationFirewallPolicy{
		Etag:     utils.String(d.Get("etag").(string)),
		Location: location,
		Properties: &webapplicationfirewallpolicies.CdnWebApplicationFirewallPolicyProperties{
			CustomRules:    expandCdnWebApplicationFirewallPoliciesCustomRuleList(d.Get("custom_rules").([]interface{})),
			ManagedRules:   expandCdnWebApplicationFirewallPoliciesManagedRuleSetList(d.Get("managed_rules").([]interface{})),
			PolicySettings: expandCdnWebApplicationFirewallPoliciesPolicySettings(d.Get("policy_settings").([]interface{})),
			RateLimitRules: expandCdnWebApplicationFirewallPoliciesRateLimitRuleList(d.Get("rate_limit_rules").([]interface{})),
		},
		Sku:  *expandCdnWebApplicationFirewallPoliciesSku(d.Get("sku").([]interface{})),
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}
	if err := client.PoliciesCreateOrUpdateThenPoll(ctx, id, props); err != nil {

		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorProfilePolicyRead(d, meta)
}

func resourceFrontdoorProfilePolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.WebApplicationFirewallPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseCdnWebApplicationFirewallPoliciesID(d.Id())
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

	d.Set("name", id.PolicyName)

	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("etag", model.Etag)
		d.Set("location", location.Normalize(model.Location))
		if props := model.Properties; props != nil {

			if err := d.Set("custom_rules", flattenCdnWebApplicationFirewallPoliciesCustomRuleList(props.CustomRules)); err != nil {
				return fmt.Errorf("setting `custom_rules`: %+v", err)
			}

			if err := d.Set("endpoint_links", flattenCdnWebApplicationFirewallPoliciesCdnEndpointArray(props.EndpointLinks)); err != nil {
				return fmt.Errorf("setting `endpoint_links`: %+v", err)
			}

			if err := d.Set("managed_rules", flattenCdnWebApplicationFirewallPoliciesManagedRuleSetList(props.ManagedRules)); err != nil {
				return fmt.Errorf("setting `managed_rules`: %+v", err)
			}

			if err := d.Set("policy_settings", flattenCdnWebApplicationFirewallPoliciesPolicySettings(props.PolicySettings)); err != nil {
				return fmt.Errorf("setting `policy_settings`: %+v", err)
			}
			d.Set("provisioning_state", props.ProvisioningState)

			if err := d.Set("rate_limit_rules", flattenCdnWebApplicationFirewallPoliciesRateLimitRuleList(props.RateLimitRules)); err != nil {
				return fmt.Errorf("setting `rate_limit_rules`: %+v", err)
			}
			d.Set("resource_state", props.ResourceState)
		}

		if err := d.Set("sku", flattenCdnWebApplicationFirewallPoliciesSku(&model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, ConvertFrontdoorProfileTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceFrontdoorProfilePolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.WebApplicationFirewallPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseCdnWebApplicationFirewallPoliciesID(d.Id())
	if err != nil {
		return err
	}

	props := webapplicationfirewallpolicies.CdnWebApplicationFirewallPolicyPatchParameters{
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}
	if err := client.PoliciesUpdateThenPoll(ctx, *id, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorProfilePolicyRead(d, meta)
}

func resourceFrontdoorProfilePolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.WebApplicationFirewallPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapplicationfirewallpolicies.ParseCdnWebApplicationFirewallPoliciesID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.PoliciesDelete(ctx, *id); err != nil {

		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandCdnWebApplicationFirewallPoliciesManagedRuleGroupOverrideArray(input []interface{}) *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride {
	results := make([]webapplicationfirewallpolicies.ManagedRuleGroupOverride, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, webapplicationfirewallpolicies.ManagedRuleGroupOverride{
			RuleGroupName: v["rule_group_name"].(string),
			Rules:         expandCdnWebApplicationFirewallPoliciesManagedRuleOverrideArray(v["rules"].([]interface{})),
		})
	}
	return &results
}

func expandCdnWebApplicationFirewallPoliciesManagedRuleOverrideArray(input []interface{}) *[]webapplicationfirewallpolicies.ManagedRuleOverride {
	results := make([]webapplicationfirewallpolicies.ManagedRuleOverride, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		actionValue := webapplicationfirewallpolicies.ActionType(v["action"].(string))
		enabledStateValue := webapplicationfirewallpolicies.ManagedRuleEnabledState(v["enabled_state"].(string))
		results = append(results, webapplicationfirewallpolicies.ManagedRuleOverride{
			Action:       &actionValue,
			EnabledState: &enabledStateValue,
			RuleId:       v["rule_id"].(string),
		})
	}
	return &results
}

func expandCdnWebApplicationFirewallPoliciesRateLimitRuleList(input []interface{}) *webapplicationfirewallpolicies.RateLimitRuleList {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &webapplicationfirewallpolicies.RateLimitRuleList{
		Rules: expandCdnWebApplicationFirewallPoliciesRateLimitRuleArray(v["rules"].([]interface{})),
	}
}

func expandCdnWebApplicationFirewallPoliciesCustomRuleList(input []interface{}) *webapplicationfirewallpolicies.CustomRuleList {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &webapplicationfirewallpolicies.CustomRuleList{
		Rules: expandCdnWebApplicationFirewallPoliciesCustomRuleArray(v["rules"].([]interface{})),
	}
}

func expandCdnWebApplicationFirewallPoliciesManagedRuleSetList(input []interface{}) *webapplicationfirewallpolicies.ManagedRuleSetList {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &webapplicationfirewallpolicies.ManagedRuleSetList{
		ManagedRuleSets: expandCdnWebApplicationFirewallPoliciesManagedRuleSetArray(v["managed_rule_sets"].([]interface{})),
	}
}

func expandCdnWebApplicationFirewallPoliciesMatchValuesArray(input []interface{}) []string {
	results := make([]string, 0)
	if len(input) == 0 {
		return results
	}

	for _, val := range input {
		if val != nil {
			results = append(results, val.(string))
		}
	}

	return results
}

func expandCdnWebApplicationFirewallPoliciesTransformTypeArray(input []interface{}) *[]webapplicationfirewallpolicies.TransformType {
	results := make([]webapplicationfirewallpolicies.TransformType, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		results = append(results, webapplicationfirewallpolicies.TransformType(item.(string)))
	}
	return &results
}

func expandCdnWebApplicationFirewallPoliciesManagedRuleSetArray(input []interface{}) *[]webapplicationfirewallpolicies.ManagedRuleSet {
	results := make([]webapplicationfirewallpolicies.ManagedRuleSet, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, webapplicationfirewallpolicies.ManagedRuleSet{
			AnomalyScore:       utils.Int64(int64(v["anomaly_score"].(int))),
			RuleGroupOverrides: expandCdnWebApplicationFirewallPoliciesManagedRuleGroupOverrideArray(v["rule_group_overrides"].([]interface{})),
			RuleSetType:        v["rule_set_type"].(string),
			RuleSetVersion:     v["rule_set_version"].(string),
		})
	}
	return &results
}

func expandCdnWebApplicationFirewallPoliciesPolicySettings(input []interface{}) *webapplicationfirewallpolicies.PolicySettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	defaultCustomBlockResponseStatusCodeValue := webapplicationfirewallpolicies.DefaultCustomBlockResponseStatusCode(v["default_custom_block_response_status_code"].(string))
	enabledStateValue := webapplicationfirewallpolicies.PolicyEnabledState(v["enabled_state"].(string))
	modeValue := webapplicationfirewallpolicies.PolicyMode(v["mode"].(string))
	return &webapplicationfirewallpolicies.PolicySettings{
		DefaultCustomBlockResponseBody:       utils.String(v["default_custom_block_response_body"].(string)),
		DefaultCustomBlockResponseStatusCode: &defaultCustomBlockResponseStatusCodeValue,
		DefaultRedirectUrl:                   utils.String(v["default_redirect_url"].(string)),
		EnabledState:                         &enabledStateValue,
		Mode:                                 &modeValue,
	}
}

func expandCdnWebApplicationFirewallPoliciesRateLimitRuleArray(input []interface{}) *[]webapplicationfirewallpolicies.RateLimitRule {
	results := make([]webapplicationfirewallpolicies.RateLimitRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		actionValue := webapplicationfirewallpolicies.ActionType(v["action"].(string))
		enabledStateValue := webapplicationfirewallpolicies.CustomRuleEnabledState(v["enabled_state"].(string))
		results = append(results, webapplicationfirewallpolicies.RateLimitRule{
			Action:                     actionValue,
			EnabledState:               &enabledStateValue,
			MatchConditions:            expandCdnWebApplicationFirewallPoliciesMatchConditionArray(v["match_conditions"].([]interface{})),
			Name:                       v["name"].(string),
			Priority:                   int64(v["priority"].(int)),
			RateLimitDurationInMinutes: int64(v["rate_limit_duration_in_minutes"].(int)),
			RateLimitThreshold:         int64(v["rate_limit_threshold"].(int)),
		})
	}
	return &results
}

func expandCdnWebApplicationFirewallPoliciesSku(input []interface{}) *webapplicationfirewallpolicies.Sku {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	nameValue := webapplicationfirewallpolicies.SkuName(v["name"].(string))
	return &webapplicationfirewallpolicies.Sku{
		Name: &nameValue,
	}
}

func expandCdnWebApplicationFirewallPoliciesCustomRuleArray(input []interface{}) *[]webapplicationfirewallpolicies.CustomRule {
	results := make([]webapplicationfirewallpolicies.CustomRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		actionValue := webapplicationfirewallpolicies.ActionType(v["action"].(string))
		enabledStateValue := webapplicationfirewallpolicies.CustomRuleEnabledState(v["enabled_state"].(string))
		results = append(results, webapplicationfirewallpolicies.CustomRule{
			Action:          actionValue,
			EnabledState:    &enabledStateValue,
			MatchConditions: expandCdnWebApplicationFirewallPoliciesMatchConditionArray(v["match_conditions"].([]interface{})),
			Name:            v["name"].(string),
			Priority:        int64(v["priority"].(int)),
		})
	}
	return &results
}

func expandCdnWebApplicationFirewallPoliciesMatchConditionArray(input []interface{}) []webapplicationfirewallpolicies.MatchCondition {
	results := make([]webapplicationfirewallpolicies.MatchCondition, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		matchVariableValue := webapplicationfirewallpolicies.WafMatchVariable(v["match_variable"].(string))
		operatorValue := webapplicationfirewallpolicies.Operator(v["operator"].(string))

		results = append(results, webapplicationfirewallpolicies.MatchCondition{
			MatchValue:      expandCdnWebApplicationFirewallPoliciesMatchValuesArray(v["match_value"].([]interface{})),
			MatchVariable:   matchVariableValue,
			NegateCondition: utils.Bool(v["negate_condition"].(bool)),
			Operator:        operatorValue,
			Selector:        utils.String(v["selector"].(string)),
			Transforms:      expandCdnWebApplicationFirewallPoliciesTransformTypeArray(v["transforms"].([]interface{})),
		})
	}
	return results
}

func flattenCdnWebApplicationFirewallPoliciesCustomRuleList(input *webapplicationfirewallpolicies.CustomRuleList) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})
	result["rules"] = flattenCdnWebApplicationFirewallPoliciesCustomRuleArray(input.Rules)
	return append(results, result)
}

func flattenCdnWebApplicationFirewallPoliciesMatchConditionArray(inputs *[]webapplicationfirewallpolicies.MatchCondition) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})
		result["match_value"] = input.MatchValue
		result["match_variable"] = input.MatchVariable

		if input.NegateCondition != nil {
			result["negate_condition"] = *input.NegateCondition
		}
		result["operator"] = input.Operator

		if input.Selector != nil {
			result["selector"] = *input.Selector
		}
		result["transforms"] = flattenCdnWebApplicationFirewallPoliciesTransformTypeArray(input.Transforms)
		results = append(results, result)
	}

	return results
}

func flattenCdnWebApplicationFirewallPoliciesCdnEndpointArray(inputs *[]webapplicationfirewallpolicies.CdnEndpoint) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})

		if input.Id != nil {
			result["id"] = *input.Id
		}
		results = append(results, result)
	}

	return results
}

func flattenCdnWebApplicationFirewallPoliciesManagedRuleSetArray(inputs *[]webapplicationfirewallpolicies.ManagedRuleSet) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})

		if input.AnomalyScore != nil {
			result["anomaly_score"] = *input.AnomalyScore
		}
		result["rule_group_overrides"] = flattenCdnWebApplicationFirewallPoliciesManagedRuleGroupOverrideArray(input.RuleGroupOverrides)
		result["rule_set_type"] = input.RuleSetType
		result["rule_set_version"] = input.RuleSetVersion
		results = append(results, result)
	}

	return results
}

func flattenCdnWebApplicationFirewallPoliciesPolicySettings(input *webapplicationfirewallpolicies.PolicySettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.DefaultCustomBlockResponseBody != nil {
		result["default_custom_block_response_body"] = *input.DefaultCustomBlockResponseBody
	}

	if input.DefaultCustomBlockResponseStatusCode != nil {
		result["default_custom_block_response_status_code"] = *input.DefaultCustomBlockResponseStatusCode
	}

	if input.DefaultRedirectUrl != nil {
		result["default_redirect_url"] = *input.DefaultRedirectUrl
	}

	if input.EnabledState != nil {
		result["enabled_state"] = *input.EnabledState
	}

	if input.Mode != nil {
		result["mode"] = *input.Mode
	}

	return append(results, result)
}

func flattenCdnWebApplicationFirewallPoliciesRateLimitRuleList(input *webapplicationfirewallpolicies.RateLimitRuleList) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})
	result["rules"] = flattenCdnWebApplicationFirewallPoliciesRateLimitRuleArray(input.Rules)
	return append(results, result)
}

func flattenCdnWebApplicationFirewallPoliciesRateLimitRuleArray(inputs *[]webapplicationfirewallpolicies.RateLimitRule) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})
		result["action"] = input.Action

		if input.EnabledState != nil {
			result["enabled_state"] = *input.EnabledState
		}

		result["match_conditions"] = flattenCdnWebApplicationFirewallPoliciesMatchConditionArray(&input.MatchConditions)
		result["name"] = input.Name
		result["priority"] = input.Priority
		result["rate_limit_duration_in_minutes"] = input.RateLimitDurationInMinutes
		result["rate_limit_threshold"] = input.RateLimitThreshold
		results = append(results, result)
	}

	return results
}

func flattenCdnWebApplicationFirewallPoliciesCustomRuleArray(inputs *[]webapplicationfirewallpolicies.CustomRule) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})
		result["action"] = input.Action

		if input.EnabledState != nil {
			result["enabled_state"] = *input.EnabledState
		}

		result["match_conditions"] = flattenCdnWebApplicationFirewallPoliciesMatchConditionArray(&input.MatchConditions)
		result["name"] = input.Name
		result["priority"] = input.Priority
		results = append(results, result)
	}

	return results
}

func flattenCdnWebApplicationFirewallPoliciesTransformTypeArray(inputs *[]webapplicationfirewallpolicies.TransformType) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, item := range *inputs {
		results = append(results, string(item))
	}

	return results
}

func flattenCdnWebApplicationFirewallPoliciesManagedRuleSetList(input *webapplicationfirewallpolicies.ManagedRuleSetList) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})
	result["managed_rule_sets"] = flattenCdnWebApplicationFirewallPoliciesManagedRuleSetArray(input.ManagedRuleSets)
	return append(results, result)
}

func flattenCdnWebApplicationFirewallPoliciesManagedRuleGroupOverrideArray(inputs *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})
		result["rule_group_name"] = input.RuleGroupName
		result["rules"] = flattenCdnWebApplicationFirewallPoliciesManagedRuleOverrideArray(input.Rules)
		results = append(results, result)
	}

	return results
}

func flattenCdnWebApplicationFirewallPoliciesManagedRuleOverrideArray(inputs *[]webapplicationfirewallpolicies.ManagedRuleOverride) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})

		if input.Action != nil {
			result["action"] = *input.Action
		}

		if input.EnabledState != nil {
			result["enabled_state"] = *input.EnabledState
		}

		result["rule_id"] = input.RuleId
		results = append(results, result)
	}

	return results
}

func flattenCdnWebApplicationFirewallPoliciesSku(input *webapplicationfirewallpolicies.Sku) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.Name != nil {
		result["name"] = *input.Name
	}

	return append(results, result)
}
