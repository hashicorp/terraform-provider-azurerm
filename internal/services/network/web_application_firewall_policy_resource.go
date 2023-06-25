package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceWebApplicationFirewallPolicy() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceWebApplicationFirewallPolicyCreateUpdate,
		Read:   resourceWebApplicationFirewallPolicyRead,
		Update: resourceWebApplicationFirewallPolicyCreateUpdate,
		Delete: resourceWebApplicationFirewallPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApplicationGatewayWebApplicationFirewallPolicyID(id)
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
								string(network.WebApplicationFirewallActionAllow),
								string(network.WebApplicationFirewallActionBlock),
								string(network.WebApplicationFirewallActionLog),
							}, false),
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
														string(network.WebApplicationFirewallMatchVariableRemoteAddr),
														string(network.WebApplicationFirewallMatchVariableRequestMethod),
														string(network.WebApplicationFirewallMatchVariableQueryString),
														string(network.WebApplicationFirewallMatchVariablePostArgs),
														string(network.WebApplicationFirewallMatchVariableRequestURI),
														string(network.WebApplicationFirewallMatchVariableRequestHeaders),
														string(network.WebApplicationFirewallMatchVariableRequestBody),
														string(network.WebApplicationFirewallMatchVariableRequestCookies),
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
											string(network.WebApplicationFirewallOperatorAny),
											string(network.WebApplicationFirewallOperatorIPMatch),
											string(network.WebApplicationFirewallOperatorGeoMatch),
											string(network.WebApplicationFirewallOperatorEqual),
											string(network.WebApplicationFirewallOperatorContains),
											string(network.WebApplicationFirewallOperatorLessThan),
											string(network.WebApplicationFirewallOperatorGreaterThan),
											string(network.WebApplicationFirewallOperatorLessThanOrEqual),
											string(network.WebApplicationFirewallOperatorGreaterThanOrEqual),
											string(network.WebApplicationFirewallOperatorBeginsWith),
											string(network.WebApplicationFirewallOperatorEndsWith),
											string(network.WebApplicationFirewallOperatorRegex),
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
												string(network.WebApplicationFirewallTransformHTMLEntityDecode),
												string(network.WebApplicationFirewallTransformLowercase),
												string(network.WebApplicationFirewallTransformRemoveNulls),
												string(network.WebApplicationFirewallTransformTrim),
												string(network.WebApplicationFirewallTransformURLDecode),
												string(network.WebApplicationFirewallTransformURLEncode),
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
								string(network.WebApplicationFirewallRuleTypeMatchRule),
								string(network.WebApplicationFirewallRuleTypeInvalid),
							}, false),
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
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
											string(network.OwaspCrsExclusionEntryMatchVariableRequestArgKeys),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestArgNames),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestArgValues),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestCookieKeys),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestCookieNames),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestCookieValues),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestHeaderKeys),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestHeaderNames),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestHeaderValues),
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
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorContains),
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorEndsWith),
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorEquals),
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorEqualsAny),
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorStartsWith),
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
																	string(network.ActionTypeAllow),
																	string(network.ActionTypeAnomalyScoring),
																	string(network.ActionTypeBlock),
																	string(network.ActionTypeLog),
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
								string(network.WebApplicationFirewallModePrevention),
								string(network.WebApplicationFirewallModeDetection),
							}, false),
							Default: string(network.WebApplicationFirewallModePrevention),
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
						"max_request_body_size_in_kb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(8, 2000),
							Default:      128,
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

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			if !features.FourPointOhBeta() {
				// Since ConflictsWith cannot be used on these properties and the properties are optional and computed, diff.GetOK may still return value even the property is not configured. Have to check the configuration with GetRawConfig
				managedRuleSetList := diff.GetRawConfig().AsValueMap()["managed_rules"].AsValueSlice()[0].AsValueMap()["managed_rule_set"].AsValueSlice()
				for _, managedRuleSetVal := range managedRuleSetList {
					ruleGroupOverrideList := managedRuleSetVal.AsValueMap()["rule_group_override"].AsValueSlice()
					for _, ruleGroupOverrideVal := range ruleGroupOverrideList {
						disabledRules := ruleGroupOverrideVal.AsValueMap()["disabled_rules"]
						ruleList := ruleGroupOverrideVal.AsValueMap()["rule"].AsValueSlice()
						if !disabledRules.IsNull() && len(ruleList) != 0 {
							return fmt.Errorf("`disabled_rules` cannot be set when `rule` is set under `rule_group_override`")
						}
					}
				}
			}

			return nil
		}),
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

func resourceWebApplicationFirewallPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WebApplicationFirewallPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewApplicationGatewayWebApplicationFirewallPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_web_application_firewall_policy", *resp.ID)
		}
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

	parameters := network.WebApplicationFirewallPolicy{
		Location: utils.String(location),
		WebApplicationFirewallPolicyPropertiesFormat: &network.WebApplicationFirewallPolicyPropertiesFormat{
			CustomRules:    expandWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(customRules),
			PolicySettings: expandWebApplicationFirewallPolicyPolicySettings(policySettings),
			ManagedRules:   expandedManagedRules,
		},
		Tags: tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceWebApplicationFirewallPolicyRead(d, meta)
}

func resourceWebApplicationFirewallPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WebApplicationFirewallPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationGatewayWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Web Application Firewall Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if webApplicationFirewallPolicyPropertiesFormat := resp.WebApplicationFirewallPolicyPropertiesFormat; webApplicationFirewallPolicyPropertiesFormat != nil {
		if err := d.Set("custom_rules", flattenWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(webApplicationFirewallPolicyPropertiesFormat.CustomRules)); err != nil {
			return fmt.Errorf("setting `custom_rules`: %+v", err)
		}
		if err := d.Set("policy_settings", flattenWebApplicationFirewallPolicyPolicySettings(webApplicationFirewallPolicyPropertiesFormat.PolicySettings)); err != nil {
			return fmt.Errorf("setting `policy_settings`: %+v", err)
		}
		if err := d.Set("managed_rules", flattenWebApplicationFirewallPolicyManagedRulesDefinition(webApplicationFirewallPolicyPropertiesFormat.ManagedRules)); err != nil {
			return fmt.Errorf("setting `managed_rules`: %+v", err)
		}
		if err := d.Set("http_listener_ids", flattenSubResourcesToIDs(webApplicationFirewallPolicyPropertiesFormat.HTTPListeners)); err != nil {
			return fmt.Errorf("setting `http_listeners`: %+v", err)
		}
		if err := d.Set("path_based_rule_ids", flattenSubResourcesToIDs(webApplicationFirewallPolicyPropertiesFormat.PathBasedRules)); err != nil {
			return fmt.Errorf("setting `path_based_rules`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceWebApplicationFirewallPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WebApplicationFirewallPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationGatewayWebApplicationFirewallPolicyID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(input []interface{}) *[]network.WebApplicationFirewallCustomRule {
	results := make([]network.WebApplicationFirewallCustomRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		name := v["name"].(string)
		priority := v["priority"].(int)
		ruleType := v["rule_type"].(string)
		matchConditions := v["match_conditions"].([]interface{})
		action := v["action"].(string)

		result := network.WebApplicationFirewallCustomRule{
			Action:          network.WebApplicationFirewallAction(action),
			MatchConditions: expandWebApplicationFirewallPolicyMatchCondition(matchConditions),
			Name:            utils.String(name),
			Priority:        utils.Int32(int32(priority)),
			RuleType:        network.WebApplicationFirewallRuleType(ruleType),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyPolicySettings(input []interface{}) *network.PolicySettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	enabled := network.WebApplicationFirewallEnabledStateDisabled
	if value, ok := v["enabled"].(bool); ok && value {
		enabled = network.WebApplicationFirewallEnabledStateEnabled
	}
	mode := v["mode"].(string)
	requestBodyCheck := v["request_body_check"].(bool)
	maxRequestBodySizeInKb := v["max_request_body_size_in_kb"].(int)
	fileUploadLimitInMb := v["file_upload_limit_in_mb"].(int)

	result := network.PolicySettings{
		State:                  enabled,
		Mode:                   network.WebApplicationFirewallMode(mode),
		RequestBodyCheck:       utils.Bool(requestBodyCheck),
		MaxRequestBodySizeInKb: utils.Int32(int32(maxRequestBodySizeInKb)),
		FileUploadLimitInMb:    utils.Int32(int32(fileUploadLimitInMb)),
	}
	return &result
}

func expandWebApplicationFirewallPolicyManagedRulesDefinition(input []interface{}, d *pluginsdk.ResourceData) (*network.ManagedRulesDefinition, error) {
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

	return &network.ManagedRulesDefinition{
		Exclusions:      expandWebApplicationFirewallPolicyExclusions(exclusions),
		ManagedRuleSets: expandedManagedRuleSets,
	}, nil
}

func expandWebApplicationFirewallPolicyExclusionManagedRules(input []interface{}) *[]network.ExclusionManagedRule {
	results := make([]network.ExclusionManagedRule, 0)
	for _, item := range input {
		ruleID := item.(string)

		result := network.ExclusionManagedRule{
			RuleID: utils.String(ruleID),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusionManagedRuleGroup(input []interface{}) *[]network.ExclusionManagedRuleGroup {
	results := make([]network.ExclusionManagedRuleGroup, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		ruleGroupName := v["rule_group_name"].(string)

		result := network.ExclusionManagedRuleGroup{
			RuleGroupName: utils.String(ruleGroupName),
		}

		if excludedRules := v["excluded_rules"].([]interface{}); len(excludedRules) > 0 {
			result.Rules = expandWebApplicationFirewallPolicyExclusionManagedRules(excludedRules)
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusionManagedRuleSet(input []interface{}) *[]network.ExclusionManagedRuleSet {
	results := make([]network.ExclusionManagedRuleSet, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		ruleSetType := v["type"].(string)
		ruleSetVersion := v["version"].(string)
		ruleGroups := make([]interface{}, 0)
		if value, exists := v["rule_group"]; exists {
			ruleGroups = value.([]interface{})
		}
		result := network.ExclusionManagedRuleSet{
			RuleSetType:    utils.String(ruleSetType),
			RuleSetVersion: utils.String(ruleSetVersion),
			RuleGroups:     expandWebApplicationFirewallPolicyExclusionManagedRuleGroup(ruleGroups),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusions(input []interface{}) *[]network.OwaspCrsExclusionEntry {
	results := make([]network.OwaspCrsExclusionEntry, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		matchVariable := v["match_variable"].(string)
		selectorMatchOperator := v["selector_match_operator"].(string)
		selector := v["selector"].(string)
		exclusionManagedRuleSets := v["excluded_rule_set"].([]interface{})

		result := network.OwaspCrsExclusionEntry{
			MatchVariable:            network.OwaspCrsExclusionEntryMatchVariable(matchVariable),
			SelectorMatchOperator:    network.OwaspCrsExclusionEntrySelectorMatchOperator(selectorMatchOperator),
			Selector:                 utils.String(selector),
			ExclusionManagedRuleSets: expandWebApplicationFirewallPolicyExclusionManagedRuleSet(exclusionManagedRuleSets),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyManagedRuleSet(input []interface{}, d *pluginsdk.ResourceData) (*[]network.ManagedRuleSet, error) {
	results := make([]network.ManagedRuleSet, 0)

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

		result := network.ManagedRuleSet{
			RuleSetType:        utils.String(ruleSetType),
			RuleSetVersion:     utils.String(ruleSetVersion),
			RuleGroupOverrides: expandedRuleGroupOverrides,
		}

		results = append(results, result)
	}
	return &results, nil
}

func expandWebApplicationFirewallPolicyRuleGroupOverrides(input []interface{}, d *pluginsdk.ResourceData, managedRuleSetIndex int) (*[]network.ManagedRuleGroupOverride, error) {
	results := make([]network.ManagedRuleGroupOverride, 0)
	for i, item := range input {
		v := item.(map[string]interface{})

		ruleGroupName := v["rule_group_name"].(string)

		result := network.ManagedRuleGroupOverride{
			RuleGroupName: utils.String(ruleGroupName),
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

func expandWebApplicationFirewallPolicyRules(input []interface{}) *[]network.ManagedRuleOverride {
	results := make([]network.ManagedRuleOverride, 0)
	for _, item := range input {
		ruleID := item.(string)

		result := network.ManagedRuleOverride{
			RuleID: utils.String(ruleID),
			State:  network.ManagedRuleEnabledStateDisabled,
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyOverrideRules(input []interface{}) *[]network.ManagedRuleOverride {
	results := make([]network.ManagedRuleOverride, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		state := network.ManagedRuleEnabledStateDisabled
		if v["enabled"].(bool) {
			state = network.ManagedRuleEnabledStateEnabled
		}

		result := network.ManagedRuleOverride{
			RuleID: utils.String(v["id"].(string)),
			State:  state,
		}

		action := v["action"].(string)
		if action != "" {
			result.Action = network.ActionType(action)
		}

		results = append(results, result)
	}

	return &results
}

func expandWebApplicationFirewallPolicyMatchCondition(input []interface{}) *[]network.MatchCondition {
	results := make([]network.MatchCondition, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		matchVariables := v["match_variables"].([]interface{})
		operator := v["operator"].(string)
		negationCondition := v["negation_condition"].(bool)
		matchValues := v["match_values"].([]interface{})
		transformsRaw := v["transforms"].(*pluginsdk.Set).List()

		var transforms []network.WebApplicationFirewallTransform
		for _, trans := range transformsRaw {
			transforms = append(transforms, network.WebApplicationFirewallTransform(trans.(string)))
		}
		result := network.MatchCondition{
			MatchValues:      utils.ExpandStringSlice(matchValues),
			MatchVariables:   expandWebApplicationFirewallPolicyMatchVariable(matchVariables),
			NegationConditon: utils.Bool(negationCondition),
			Operator:         network.WebApplicationFirewallOperator(operator),
			Transforms:       &transforms,
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyMatchVariable(input []interface{}) *[]network.MatchVariable {
	results := make([]network.MatchVariable, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		variableName := v["variable_name"].(string)
		selector := v["selector"].(string)

		result := network.MatchVariable{
			Selector:     utils.String(selector),
			VariableName: network.WebApplicationFirewallMatchVariable(variableName),
		}

		results = append(results, result)
	}
	return &results
}

func flattenWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(input *[]network.WebApplicationFirewallCustomRule) []interface{} {
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
		v["match_conditions"] = flattenWebApplicationFirewallPolicyMatchCondition(item.MatchConditions)
		if priority := item.Priority; priority != nil {
			v["priority"] = int(*priority)
		}
		v["rule_type"] = string(item.RuleType)

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPolicyPolicySettings(input *network.PolicySettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["enabled"] = input.State == network.WebApplicationFirewallEnabledStateEnabled
	result["mode"] = string(input.Mode)
	result["request_body_check"] = input.RequestBodyCheck
	result["max_request_body_size_in_kb"] = int(*input.MaxRequestBodySizeInKb)
	result["file_upload_limit_in_mb"] = int(*input.FileUploadLimitInMb)

	return []interface{}{result}
}

func flattenWebApplicationFirewallPolicyManagedRulesDefinition(input *network.ManagedRulesDefinition) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	v := make(map[string]interface{})

	v["exclusion"] = flattenWebApplicationFirewallPolicyExclusions(input.Exclusions)
	v["managed_rule_set"] = flattenWebApplicationFirewallPolicyManagedRuleSets(input.ManagedRuleSets)

	results = append(results, v)

	return results
}

func flattenWebApplicationFirewallPolicyExclusionManagedRules(input *[]network.ExclusionManagedRule) []string {
	results := make([]string, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		if item.RuleID != nil {
			v := *item.RuleID
			results = append(results, v)
		}
	}

	return results
}

func flattenWebApplicationFirewallPolicyExclusionManagedRuleGroups(input *[]network.ExclusionManagedRuleGroup) []interface{} {
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

func flattenWebApplicationFirewallPolicyExclusionManagedRuleSets(input *[]network.ExclusionManagedRuleSet) []interface{} {
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

func flattenWebApplicationFirewallPolicyExclusions(input *[]network.OwaspCrsExclusionEntry) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		selector := item.Selector

		v["match_variable"] = string(item.MatchVariable)
		if selector != nil {
			v["selector"] = *selector
		}

		v["selector_match_operator"] = string(item.SelectorMatchOperator)
		v["excluded_rule_set"] = flattenWebApplicationFirewallPolicyExclusionManagedRuleSets(item.ExclusionManagedRuleSets)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyManagedRuleSets(input *[]network.ManagedRuleSet) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		v["type"] = item.RuleSetType
		v["version"] = item.RuleSetVersion
		v["rule_group_override"] = flattenWebApplicationFirewallPolicyRuleGroupOverrides(item.RuleGroupOverrides)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyRuleGroupOverrides(input *[]network.ManagedRuleGroupOverride) []interface{} {
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

func flattenWebApplicationFirewallPolicyRules(input *[]network.ManagedRuleOverride) []string {
	results := make([]string, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		if (item.State == "" || item.State == network.ManagedRuleEnabledStateDisabled) && item.RuleID != nil {
			v := *item.RuleID

			results = append(results, v)
		}
	}

	return results
}

func flattenWebApplicationFirewallPolicyOverrideRules(input *[]network.ManagedRuleOverride) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})
		if item.RuleID != nil {
			v["id"] = *item.RuleID
		}

		v["enabled"] = item.State == network.ManagedRuleEnabledStateEnabled

		v["action"] = string(item.Action)

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPolicyMatchCondition(input *[]network.MatchCondition) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		var transforms []interface{}
		if item.Transforms != nil {
			for _, trans := range *item.Transforms {
				transforms = append(transforms, string(trans))
			}
		}
		v["match_values"] = utils.FlattenStringSlice(item.MatchValues)
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

func flattenWebApplicationFirewallPolicyMatchVariable(input *[]network.MatchVariable) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if selector := item.Selector; selector != nil {
			v["selector"] = *selector
		}
		v["variable_name"] = string(item.VariableName)

		results = append(results, v)
	}

	return results
}
