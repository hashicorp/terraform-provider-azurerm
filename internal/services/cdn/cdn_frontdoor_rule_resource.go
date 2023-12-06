// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	cdn "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	cdnFrontDoorRuleActions "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorruleactions"
	cdnFrontDoorRuleConditions "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorruleconditions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorRuleCreate,
		Read:   resourceCdnFrontDoorRuleRead,
		Update: resourceCdnFrontDoorRuleUpdate,
		Delete: resourceCdnFrontDoorRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CdnFrontDoorRuleName,
			},

			"cdn_frontdoor_rule_set_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorRuleSetID,
			},

			"behavior_on_match": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(cdn.MatchProcessingBehaviorContinue),
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.MatchProcessingBehaviorContinue),
					string(cdn.MatchProcessingBehaviorStop),
				}, false),
			},

			"order": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"actions": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"url_redirect_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"redirect_type": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.RedirectTypeMoved),
											string(cdn.RedirectTypeFound),
											string(cdn.RedirectTypeTemporaryRedirect),
											string(cdn.RedirectTypePermanentRedirect),
										}, false),
									},

									"redirect_protocol": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(cdn.DestinationProtocolMatchRequest),
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.DestinationProtocolMatchRequest),
											string(cdn.DestinationProtocolHTTP),
											string(cdn.DestinationProtocolHTTPS),
										}, false),
									},

									// NOTE: it is valid for the destination path to be an empty string,
									// Leave blank to preserve the incoming path. Issue #18249
									"destination_path": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: validate.CdnFrontDoorUrlRedirectActionDestinationPath,
									},

									// NOTE: it is valid for the destination hostname to be an empty string.
									// Leave blank to preserve the incoming host. Issue #18249
									"destination_hostname": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(0, 2048),
									},

									// NOTE: it is valid for the query string to be an empty string.
									// Leave blank to preserve the incoming query string. Issue #18249 & #19682
									"query_string": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  "",
										// Update validation logic to match RP. Issue #19097
										ValidateFunc: validate.CdnFrontDoorUrlRedirectActionQueryString,
									},

									// NOTE: it is valid for the destination fragment to be an empty string.
									// Leave blank to preserve the incoming fragment. Issue #18249
									"destination_fragment": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: validation.StringLenBetween(0, 1024),
									},
								},
							},
						},

						"url_rewrite_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"source_pattern": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"destination": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"preserve_unmatched_path": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},

						"request_header_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"header_action": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.HeaderActionAppend),
											string(cdn.HeaderActionOverwrite),
											string(cdn.HeaderActionDelete),
										}, false),
									},

									"header_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"value": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"response_header_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"header_action": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.HeaderActionAppend),
											string(cdn.HeaderActionOverwrite),
											string(cdn.HeaderActionDelete),
										}, false),
									},

									"header_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"value": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"route_configuration_override_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"cdn_frontdoor_origin_group_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validate.FrontDoorOriginGroupID,
									},

									// Removed Default value for issue #18889
									"forwarding_protocol": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.ForwardingProtocolHTTPOnly),
											string(cdn.ForwardingProtocolHTTPSOnly),
											string(cdn.ForwardingProtocolMatchRequest),
										}, false),
									},

									// Removed Default value for issue #19008
									"query_string_caching_behavior": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.RuleQueryStringCachingBehaviorIgnoreQueryString),
											string(cdn.RuleQueryStringCachingBehaviorUseQueryString),
											string(cdn.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
											string(cdn.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
										}, false),
									},

									// NOTE: CSV implemented as a list, code already written for the expanded and flatten to CSV
									// not valid when IncludeAll or ExcludeAll behavior is defined
									"query_string_parameters": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 100,

										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"compression_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},

									// Exposed Disabled for issue #19008
									"cache_behavior": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.RuleCacheBehaviorHonorOrigin),
											string(cdn.RuleCacheBehaviorOverrideAlways),
											string(cdn.RuleCacheBehaviorOverrideIfOriginMissing),
											string(cdn.RuleIsCompressionEnabledDisabled),
										}, false),
									},

									// Made Optional for issue #19008
									"cache_duration": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validate.CdnFrontDoorCacheDuration,
									},
								},
							},
						},
					},
				},
			},

			"conditions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"remote_address_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperatorRemoteAddress(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
								},
							},
						},

						"request_method_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperatorEqualOnly(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorRequestMethodMatchValues(),
								},
							},
						},

						"query_string_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"post_args_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									// In the API this is called selector
									"post_args_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"request_uri_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"request_header_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									// In the API this is called selector
									// match_values are invalid if operator is 'Any'
									"header_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"request_body_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValuesRequired(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"request_scheme_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperatorEqualOnly(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorProtocolMatchValues(),
								},
							},
						},

						"url_path_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorUrlPathConditionMatchValues(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"url_file_extension_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValuesRequired(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"url_filename_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									// making optional for issue #23504
									"match_values": schemaCdnFrontDoorMatchValues(),
									"transforms":   schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"http_version_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperatorEqualOnly(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorHttpVersionMatchValues(),
								},
							},
						},

						"cookies_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									// In the API this is called selector
									"cookie_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"is_device_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperatorEqualOnly(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorIsDeviceMatchValues(),
								},
							},
						},

						"socket_address_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperatorSocketAddress(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
								},
							},
						},

						"client_port_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
								},
							},
						},

						"server_port_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorServerPortMatchValues(),
								},
							},
						},

						"host_name_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperator(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorMatchValues(),
									"transforms":       schemaCdnFrontDoorRuleTransforms(),
								},
							},
						},

						"ssl_protocol_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         schemaCdnFrontDoorOperatorEqualOnly(),
									"negate_condition": schemaCdnFrontDoorNegateCondition(),
									"match_values":     schemaCdnFrontDoorSslProtocolMatchValues(),
								},
							},
						},
					},
				},
			},

			"cdn_frontdoor_rule_set_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCdnFrontDoorRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	ruleSet, err := parse.FrontDoorRuleSetID(d.Get("cdn_frontdoor_rule_set_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorRuleID(ruleSet.SubscriptionId, ruleSet.ResourceGroup, ruleSet.ProfileName, ruleSet.RuleSetName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_rule", id.ID())
	}

	matchProcessingBehaviorValue := cdn.MatchProcessingBehavior(d.Get("behavior_on_match").(string))
	order := d.Get("order").(int)

	actions, err := expandFrontdoorDeliveryRuleActions(d.Get("actions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding 'actions': %+v", err)
	}

	conditions, err := expandFrontdoorDeliveryRuleConditions(d.Get("conditions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding 'conditions': %+v", err)
	}

	props := cdn.Rule{
		RuleProperties: &cdn.RuleProperties{
			Actions:                 &actions,
			Conditions:              &conditions,
			MatchProcessingBehavior: matchProcessingBehaviorValue,
			RuleSetName:             &ruleSet.RuleSetName,
			Order:                   utils.Int32(int32(order)),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCdnFrontDoorRuleRead(d, meta)
}

func resourceCdnFrontDoorRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRuleID(d.Id())
	if err != nil {
		return err
	}

	ruleSet := parse.NewFrontDoorRuleSetID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.RuleSetName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RuleName)
	d.Set("cdn_frontdoor_rule_set_id", ruleSet.ID())

	if props := resp.RuleProperties; props != nil {
		d.Set("behavior_on_match", props.MatchProcessingBehavior)
		d.Set("order", props.Order)

		// BUG: RuleSetName is not being returned by the API
		// Tracking issue opened: https://github.com/Azure/azure-rest-api-specs/issues/20560
		d.Set("cdn_frontdoor_rule_set_name", ruleSet.RuleSetName)

		actions, err := flattenFrontdoorDeliveryRuleActions(props.Actions)
		if err != nil {
			return fmt.Errorf("setting 'actions': %+v", err)
		}
		d.Set("actions", actions)

		conditions, err := flattenFrontdoorDeliveryRuleConditions(props.Conditions)
		if err != nil {
			return fmt.Errorf("setting 'conditions': %+v", err)
		}
		d.Set("conditions", conditions)
	}

	return nil
}

func resourceCdnFrontDoorRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRuleID(d.Id())
	if err != nil {
		return err
	}

	props := cdn.RuleUpdateParameters{
		RuleUpdatePropertiesParameters: &cdn.RuleUpdatePropertiesParameters{},
	}

	if d.HasChange("behavior_on_match") {
		matchProcessingBehaviorValue := cdn.MatchProcessingBehavior(d.Get("behavior_on_match").(string))
		props.RuleUpdatePropertiesParameters.MatchProcessingBehavior = matchProcessingBehaviorValue
	}

	if d.HasChange("order") {
		order := d.Get("order").(int)
		props.RuleUpdatePropertiesParameters.Order = utils.Int32(int32(order))
	}

	if d.HasChange("actions") {
		actions, err := expandFrontdoorDeliveryRuleActions(d.Get("actions").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding 'actions': %+v", err)
		}

		props.RuleUpdatePropertiesParameters.Actions = &actions
	}

	if d.HasChange("conditions") {
		conditions, err := expandFrontdoorDeliveryRuleConditions(d.Get("conditions").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding 'conditions': %+v", err)
		}

		if len(conditions) > 10 {
			return fmt.Errorf("expanding 'conditions': configuration file exceeds the maximum of 10 match conditions, got %d", len(conditions))
		}

		props.RuleUpdatePropertiesParameters.Conditions = &conditions
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorRuleRead(d, meta)
}

func resourceCdnFrontDoorRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandFrontdoorDeliveryRuleActions(input []interface{}) ([]cdn.BasicDeliveryRuleAction, error) {
	results := make([]cdn.BasicDeliveryRuleAction, 0)
	if len(input) == 0 {
		return results, nil
	}

	type expandfunc func(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error)

	m := *cdnFrontDoorRuleActions.InitializeCdnFrontDoorActionMappings()

	actions := map[string]expandfunc{
		m.RouteConfigurationOverride.ConfigName: cdnFrontDoorRuleActions.ExpandCdnFrontDoorRouteConfigurationOverrideAction,
		m.RequestHeader.ConfigName:              cdnFrontDoorRuleActions.ExpandCdnFrontDoorRequestHeaderAction,
		m.ResponseHeader.ConfigName:             cdnFrontDoorRuleActions.ExpandCdnFrontDoorResponseHeaderAction,
		m.URLRedirect.ConfigName:                cdnFrontDoorRuleActions.ExpandCdnFrontDoorUrlRedirectAction,
		m.URLRewrite.ConfigName:                 cdnFrontDoorRuleActions.ExpandCdnFrontDoorUrlRewriteAction,
	}

	basicDeliveryRuleAction := input[0].(map[string]interface{})

	for actionName, expand := range actions {
		raw := basicDeliveryRuleAction[actionName].([]interface{})
		expanded, err := expand(raw)
		if err != nil {
			return nil, err
		}

		if expanded != nil {
			if actionName == m.URLRewrite.ConfigName && len(*expanded) > 1 {
				return nil, fmt.Errorf("the 'url_rewrite_action' is only allowed once in the 'actions' match block, got %d", len(*expanded))
			}

			if actionName == m.URLRedirect.ConfigName && len(*expanded) > 1 {
				return nil, fmt.Errorf("the 'url_redirect_action' is only allowed once in the 'actions' match block, got %d", len(*expanded))
			}

			if actionName == m.RouteConfigurationOverride.ConfigName && len(*expanded) > 1 {
				return nil, fmt.Errorf("the 'route_configuration_override_action' is only allowed once in the 'actions' match block, got %d", len(*expanded))
			}

			results = append(results, *expanded...)
		}
	}

	if len(results) > 5 {
		return nil, fmt.Errorf("the 'actions' match block may only contain up to 5 match actions, got %d", len(results))
	}

	if err := validate.CdnFrontDoorActionsBlock(results); err != nil {
		return nil, err
	}

	return results, nil
}

func expandFrontdoorDeliveryRuleConditions(input []interface{}) ([]cdn.BasicDeliveryRuleCondition, error) {
	results := make([]cdn.BasicDeliveryRuleCondition, 0)
	if len(input) == 0 || input[0] == nil {
		return results, nil
	}

	type expandfunc func(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error)
	m := cdnFrontDoorRuleConditions.InitializeCdnFrontDoorConditionMappings()

	conditions := map[string]expandfunc{
		m.ClientPort.ConfigName:       cdnFrontDoorRuleConditions.ExpandCdnFrontDoorClientPortCondition,
		m.Cookies.ConfigName:          cdnFrontDoorRuleConditions.ExpandCdnFrontDoorCookiesCondition,
		m.HostName.ConfigName:         cdnFrontDoorRuleConditions.ExpandCdnFrontDoorHostNameCondition,
		m.HttpVersion.ConfigName:      cdnFrontDoorRuleConditions.ExpandCdnFrontDoorHttpVersionCondition,
		m.IsDevice.ConfigName:         cdnFrontDoorRuleConditions.ExpandCdnFrontDoorIsDeviceCondition,
		m.PostArgs.ConfigName:         cdnFrontDoorRuleConditions.ExpandCdnFrontDoorPostArgsCondition,
		m.QueryString.ConfigName:      cdnFrontDoorRuleConditions.ExpandCdnFrontDoorQueryStringCondition,
		m.RemoteAddress.ConfigName:    cdnFrontDoorRuleConditions.ExpandCdnFrontDoorRemoteAddressCondition,
		m.RequestBody.ConfigName:      cdnFrontDoorRuleConditions.ExpandCdnFrontDoorRequestBodyCondition,
		m.RequestHeader.ConfigName:    cdnFrontDoorRuleConditions.ExpandCdnFrontDoorRequestHeaderCondition,
		m.RequestMethod.ConfigName:    cdnFrontDoorRuleConditions.ExpandCdnFrontDoorRequestMethodCondition,
		m.RequestScheme.ConfigName:    cdnFrontDoorRuleConditions.ExpandCdnFrontDoorRequestSchemeCondition,
		m.RequestUri.ConfigName:       cdnFrontDoorRuleConditions.ExpandCdnFrontDoorRequestUriCondition,
		m.ServerPort.ConfigName:       cdnFrontDoorRuleConditions.ExpandCdnFrontDoorServerPortCondition,
		m.SocketAddress.ConfigName:    cdnFrontDoorRuleConditions.ExpandCdnFrontDoorSocketAddressCondition,
		m.SslProtocol.ConfigName:      cdnFrontDoorRuleConditions.ExpandCdnFrontDoorSslProtocolCondition,
		m.UrlFileExtension.ConfigName: cdnFrontDoorRuleConditions.ExpandCdnFrontDoorUrlFileExtensionCondition,
		m.UrlFilename.ConfigName:      cdnFrontDoorRuleConditions.ExpandCdnFrontDoorUrlFileNameCondition,
		m.UrlPath.ConfigName:          cdnFrontDoorRuleConditions.ExpandCdnFrontDoorUrlPathCondition,
	}

	basicDeliveryRuleCondition := input[0].(map[string]interface{})

	for conditionName, expand := range conditions {
		raw := basicDeliveryRuleCondition[conditionName].([]interface{})
		if len(raw) > 0 {
			expanded, err := expand(raw)
			if err != nil {
				return nil, err
			}

			if expanded != nil {
				results = append(results, *expanded...)
			}
		}
	}

	if len(results) > 10 {
		return nil, fmt.Errorf("the 'conditions' match block may only contain up to 10 match conditions, got %d", len(results))
	}

	return results, nil
}

func flattenFrontdoorDeliveryRuleConditions(input *[]cdn.BasicDeliveryRuleCondition) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	c := cdnFrontDoorRuleConditions.InitializeCdnFrontDoorConditionMappings()

	clientPortCondition := make([]interface{}, 0)
	cookiesCondition := make([]interface{}, 0)
	hostNameCondition := make([]interface{}, 0)
	httpVersionCondition := make([]interface{}, 0)
	isDeviceCondition := make([]interface{}, 0)
	postArgsCondition := make([]interface{}, 0)
	queryStringCondition := make([]interface{}, 0)
	remoteAddressCondition := make([]interface{}, 0)
	requestBodyCondition := make([]interface{}, 0)
	requestHeaderCondition := make([]interface{}, 0)
	requestMethodCondition := make([]interface{}, 0)
	requestSchemeCondition := make([]interface{}, 0)
	requestURICondition := make([]interface{}, 0)
	serverPortCondition := make([]interface{}, 0)
	socketAddressCondition := make([]interface{}, 0)
	sslProtocolCondition := make([]interface{}, 0)
	urlFileExtensionCondition := make([]interface{}, 0)
	urlFilenameCondition := make([]interface{}, 0)
	urlPathCondition := make([]interface{}, 0)

	for _, BasicDeliveryRuleCondition := range *input {
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleClientPortCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorClientPortCondition(condition)
			if err != nil {
				return nil, err
			}

			clientPortCondition = append(clientPortCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleCookiesCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorCookiesCondition(condition)
			if err != nil {
				return nil, err
			}

			cookiesCondition = append(cookiesCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleHostNameCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorHostNameCondition(condition)
			if err != nil {
				return nil, err
			}

			hostNameCondition = append(hostNameCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleHTTPVersionCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorHttpVersionCondition(condition)
			if err != nil {
				return nil, err
			}
			httpVersionCondition = append(httpVersionCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleIsDeviceCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorIsDeviceCondition(condition)
			if err != nil {
				return nil, err
			}

			isDeviceCondition = append(isDeviceCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRulePostArgsCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorPostArgsCondition(condition)
			if err != nil {
				return nil, err
			}

			postArgsCondition = append(postArgsCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleQueryStringCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorQueryStringCondition(condition)
			if err != nil {
				return nil, err
			}

			queryStringCondition = append(queryStringCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRemoteAddressCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorRemoteAddressCondition(condition)
			if err != nil {
				return nil, err
			}

			remoteAddressCondition = append(remoteAddressCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestBodyCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorRequestBodyCondition(condition)
			if err != nil {
				return nil, err
			}

			requestBodyCondition = append(requestBodyCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestHeaderCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorRequestHeaderCondition(condition)
			if err != nil {
				return nil, err
			}

			requestHeaderCondition = append(requestHeaderCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestMethodCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorRequestMethodCondition(condition)
			if err != nil {
				return nil, err
			}

			requestMethodCondition = append(requestMethodCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestSchemeCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorRequestSchemeCondition(condition)
			if err != nil {
				return nil, err
			}

			requestSchemeCondition = append(requestSchemeCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestURICondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorRequestUriCondition(condition)
			if err != nil {
				return nil, err
			}

			requestURICondition = append(requestURICondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleServerPortCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorServerPortCondition(condition)
			if err != nil {
				return nil, err
			}

			serverPortCondition = append(serverPortCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleSocketAddrCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorSocketAddressCondition(condition)
			if err != nil {
				return nil, err
			}

			socketAddressCondition = append(socketAddressCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleSslProtocolCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorSslProtocolCondition(condition)
			if err != nil {
				return nil, err
			}

			sslProtocolCondition = append(sslProtocolCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLFileExtensionCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorUrlFileExtensionCondition(condition)
			if err != nil {
				return nil, err
			}

			urlFileExtensionCondition = append(urlFileExtensionCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLFileNameCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorUrlFileNameCondition(condition)
			if err != nil {
				return nil, err
			}

			urlFilenameCondition = append(urlFilenameCondition, flattened)
			continue
		}

		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLPathCondition(); ok {
			flattened, err := cdnFrontDoorRuleConditions.FlattenFrontdoorUrlPathCondition(condition)
			if err != nil {
				return nil, err
			}

			urlPathCondition = append(urlPathCondition, flattened)
			continue
		}

		return nil, fmt.Errorf("unknown BasicDeliveryRuleCondition encountered")
	}

	conditions := map[string]interface{}{
		c.ClientPort.ConfigName:       clientPortCondition,
		c.Cookies.ConfigName:          cookiesCondition,
		c.HostName.ConfigName:         hostNameCondition,
		c.HttpVersion.ConfigName:      httpVersionCondition,
		c.IsDevice.ConfigName:         isDeviceCondition,
		c.PostArgs.ConfigName:         postArgsCondition,
		c.QueryString.ConfigName:      queryStringCondition,
		c.RemoteAddress.ConfigName:    remoteAddressCondition,
		c.RequestBody.ConfigName:      requestBodyCondition,
		c.RequestHeader.ConfigName:    requestHeaderCondition,
		c.RequestMethod.ConfigName:    requestMethodCondition,
		c.RequestScheme.ConfigName:    requestSchemeCondition,
		c.RequestUri.ConfigName:       requestURICondition,
		c.ServerPort.ConfigName:       serverPortCondition,
		c.SocketAddress.ConfigName:    socketAddressCondition,
		c.SslProtocol.ConfigName:      sslProtocolCondition,
		c.UrlFileExtension.ConfigName: urlFileExtensionCondition,
		c.UrlFilename.ConfigName:      urlFilenameCondition,
		c.UrlPath.ConfigName:          urlPathCondition,
	}

	// NOTE: Since we are always returning something no matter what this causes
	// a perpetual diff during plan. Only return the conditions map if
	// it actually has a condition defined within it, else return an empty
	// slice
	output := []interface{}{conditions}
	if !ruleHasDeliveryRuleConditions(conditions) {
		output = results
	}

	return output, nil
}

func flattenFrontdoorDeliveryRuleActions(input *[]cdn.BasicDeliveryRuleAction) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	a := cdnFrontDoorRuleActions.InitializeCdnFrontDoorActionMappings()

	requestHeaderActions := make([]interface{}, 0)
	responseHeaderActions := make([]interface{}, 0)
	routeConfigOverrideActions := make([]interface{}, 0)
	urlRedirectActions := make([]interface{}, 0)
	urlRewriteActions := make([]interface{}, 0)

	for _, item := range *input {
		if action, ok := item.AsDeliveryRuleRouteConfigurationOverrideAction(); ok {
			flattened, err := cdnFrontDoorRuleActions.FlattenCdnFrontDoorRouteConfigurationOverrideAction(*action)
			if err != nil {
				return nil, fmt.Errorf("'route_configuration_override_action' unable to parse 'cdn_frontdoor_origin_group_id': %+v", err)
			}

			routeConfigOverrideActions = append(routeConfigOverrideActions, flattened)
			continue
		}

		if action, ok := item.AsDeliveryRuleRequestHeaderAction(); ok {
			if action.Parameters == nil {
				return nil, fmt.Errorf("'parameters' was nil for Delivery Rule Request Header")
			}
			flattened := cdnFrontDoorRuleActions.FlattenHeaderActionParameters(action.Parameters)
			requestHeaderActions = append(requestHeaderActions, flattened)
			continue
		}

		if action, ok := item.AsDeliveryRuleResponseHeaderAction(); ok {
			if action.Parameters == nil {
				return nil, fmt.Errorf("'parameters' was nil for Delivery Rule Response Header")
			}
			flattened := cdnFrontDoorRuleActions.FlattenHeaderActionParameters(action.Parameters)
			responseHeaderActions = append(responseHeaderActions, flattened)
			continue
		}

		if action, ok := item.AsURLRedirectAction(); ok {
			flattened := cdnFrontDoorRuleActions.FlattenCdnFrontDoorUrlRedirectAction(*action)
			urlRedirectActions = append(urlRedirectActions, flattened)
			continue
		}

		if action, ok := item.AsURLRewriteAction(); ok {
			flattened := cdnFrontDoorRuleActions.FlattenCdnFrontDoorUrlRewriteAction(*action)
			urlRewriteActions = append(urlRewriteActions, flattened)
			continue
		}
	}

	if len(requestHeaderActions) == 0 && len(responseHeaderActions) == 0 && len(routeConfigOverrideActions) == 0 && len(urlRedirectActions) == 0 && len(urlRewriteActions) == 0 {
		return []interface{}{}, nil
	}

	return []interface{}{
		map[string]interface{}{
			a.RequestHeader.ConfigName:              requestHeaderActions,
			a.ResponseHeader.ConfigName:             responseHeaderActions,
			a.RouteConfigurationOverride.ConfigName: routeConfigOverrideActions,
			a.URLRedirect.ConfigName:                urlRedirectActions,
			a.URLRewrite.ConfigName:                 urlRewriteActions,
		},
	}, nil
}
