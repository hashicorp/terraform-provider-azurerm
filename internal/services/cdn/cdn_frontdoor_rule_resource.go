package cdn

import (
	"fmt"
	"time"

	cdn "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	cdnfrontdoorruleactions "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorruleactions"
	cdnfrontdoorruleconditions "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorruleconditions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorRuleCreate,
		Read:   resourceCdnFrontdoorRuleRead,
		Update: resourceCdnFrontdoorRuleUpdate,
		Delete: resourceCdnFrontdoorRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateCdnFrontdoorRuleName,
			},

			"cdn_frontdoor_rule_set_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorRuleSetID,
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

						// Name: UrlRedirect
						// DeliveryRuleUrlRedirectActionParameters
						"url_redirect_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,

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

									"destination_path": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: ValidateCdnFrontdoorUrlRedirectActionDestinationPath,
									},

									"destination_hostname": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"query_string": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: ValidateCdnFrontdoorUrlRedirectActionQueryString,
									},

									"destination_fragment": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						// Name: URLRewrite
						// URLRewriteAction
						"url_rewrite_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,

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

						// Name: ModifyRequestHeader
						// DeliveryRuleRequestHeaderAction
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

						// Name: ModifyResponseHeader (NameBasicDeliveryRuleActionNameModifyResponseHeader)
						// DeliveryRuleResponseHeaderAction
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

						// 'Cache_Expiration’, ‘Cache_Key_Query_String’ and ‘Origin_Group_Override' Actions
						// are only supported in the 2020-09-01 API version. All calls for these Actions
						// will fail 90 days after GA of the AFDx service.

						// Name: RouteConfigurationOverride (NameBasicDeliveryRuleActionNameRouteConfigurationOverride)
						// DeliveryRuleRouteConfigurationOverrideAction
						"route_configuration_override_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"cdn_frontdoor_origin_group_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validate.FrontdoorOriginGroupID,
									},

									"forwarding_protocol": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(cdn.ForwardingProtocolMatchRequest),
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.ForwardingProtocolHTTPOnly),
											string(cdn.ForwardingProtocolHTTPSOnly),
											string(cdn.ForwardingProtocolMatchRequest),
										}, false),
									},

									"query_string_caching_behavior": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(cdn.RuleQueryStringCachingBehaviorIgnoreQueryString),
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.RuleQueryStringCachingBehaviorIgnoreQueryString),
											string(cdn.RuleQueryStringCachingBehaviorUseQueryString),
											string(cdn.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
											string(cdn.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
										}, false),
									},

									// CSV implemented as a list, code alread written for the expaned and flatten to CSV
									// not valid when IncludeAll or ExcludeAll behavior is defined
									"query_string_parameters": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 100,

										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									// Content won't be compressed on AzureFrontDoor when requested content is smaller than 1 byte or larger than 1 MB.
									"compression_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"cache_behavior": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(cdn.RuleCacheBehaviorHonorOrigin),
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.RuleCacheBehaviorHonorOrigin),
											string(cdn.RuleCacheBehaviorOverrideAlways),
											string(cdn.RuleCacheBehaviorOverrideIfOriginMissing),
										}, false),
									},

									// Allowed format is d.HH:MM:SS or HH:MM:SS if duration is less than a day
									// maximum duration is 366 days(e.g. 365.23:59:59)
									"cache_duration": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: ValidateCdnFrontdoorCacheDuration,
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
									"operator":         SchemaCdnFrontdoorOperatorRemoteAddress(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
								},
							},
						},

						"request_method_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperatorEqualOnly(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorRequestMethodMatchValues(),
								},
							},
						},

						"query_string_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
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
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
								},
							},
						},

						"request_uri_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
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
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
								},
							},
						},

						"request_body_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValuesRequired(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
								},
							},
						},

						"request_scheme_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperatorEqualOnly(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorProtocolMatchValues(),
								},
							},
						},

						"url_path_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorUrlPathConditionMatchValues(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
								},
							},
						},

						"url_file_extension_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValuesRequired(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
								},
							},
						},

						"url_filename_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValuesRequired(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
								},
							},
						},

						"http_version_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperatorEqualOnly(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorHttpVersionMatchValues(),
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

									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
								},
							},
						},

						"is_device_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperatorEqualOnly(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorIsDeviceMatchValues(),
								},
							},
						},

						// DeliveryRuleSocketAddrCondition
						"socket_address_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperatorSocketAddress(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
								},
							},
						},

						// DeliveryRuleClientPortCondition
						"client_port_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
								},
							},
						},

						// DeliveryRuleServerPortCondition
						"server_port_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorServerPortMatchValues(),
								},
							},
						},

						// DeliveryRuleHostNameCondition
						// NOTE: The match values do not have to adhere to RFC 1123 standards
						// for valid hostnames. Service team delegates that responsibility
						// to the author of the rule and does not validate passed match values.
						"host_name_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperator(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorMatchValues(),
									"transforms":       SchemaCdnFrontdoorRuleTransforms(),
								},
							},
						},

						// DeliveryRuleSslProtocolCondition
						"ssl_protocol_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"operator":         SchemaCdnFrontdoorOperatorEqualOnly(),
									"negate_condition": SchemaCdnFrontdoorNegateCondition(),
									"match_values":     SchemaCdnFrontdoorSslProtocolMatchValues(),
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

func resourceCdnFrontdoorRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	ruleSetId, err := parse.FrontdoorRuleSetID(d.Get("cdn_frontdoor_rule_set_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorRuleID(ruleSetId.SubscriptionId, ruleSetId.ResourceGroup, ruleSetId.ProfileName, ruleSetId.RuleSetName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_rule", id.ID())
		}
	}

	matchProcessingBehaviorValue := cdn.MatchProcessingBehavior(d.Get("behavior_on_match").(string))
	order := d.Get("order").(int)

	actions, err := expandFrontdoorDeliveryRuleActions(d.Get("actions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "actions", err)
	}

	conditions, err := expandFrontdoorDeliveryRuleConditions(d.Get("conditions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "conditions", err)
	}

	props := cdn.Rule{
		RuleProperties: &cdn.RuleProperties{
			Actions:                 &actions,
			Conditions:              &conditions,
			MatchProcessingBehavior: matchProcessingBehaviorValue,
			RuleSetName:             &ruleSetId.RuleSetName,
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

	return resourceCdnFrontdoorRuleRead(d, meta)
}

func resourceCdnFrontdoorRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRuleID(d.Id())
	if err != nil {
		return err
	}

	ruleSetId := parse.NewFrontdoorRuleSetID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.RuleSetName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RuleName)
	d.Set("cdn_frontdoor_rule_set_id", ruleSetId.ID())

	if props := resp.RuleProperties; props != nil {
		d.Set("behavior_on_match", props.MatchProcessingBehavior)
		d.Set("order", props.Order)

		// BUG: RuleSetName is not being returned by the API
		d.Set("cdn_frontdoor_rule_set_name", ruleSetId.RuleSetName)

		actions, err := flattenFrontdoorDeliveryRuleActions(props.Actions)
		if err != nil {
			return fmt.Errorf("setting %q: %+v", "actions", err)
		}
		d.Set("actions", actions)

		conditions, err := flattenFrontdoorDeliveryRuleConditions(props.Conditions)
		if err != nil {
			return fmt.Errorf("setting %q: %+v", "conditions", err)
		}

		d.Set("conditions", conditions)
	}

	return nil
}

func resourceCdnFrontdoorRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRuleID(d.Id())
	if err != nil {
		return err
	}

	matchProcessingBehaviorValue := cdn.MatchProcessingBehavior(d.Get("behavior_on_match").(string))
	order := d.Get("order").(int)

	actions, err := expandFrontdoorDeliveryRuleActions(d.Get("actions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "actions", err)
	}

	conditions, err := expandFrontdoorDeliveryRuleConditions(d.Get("conditions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "conditions", err)
	}

	if len(conditions) > 10 {
		return fmt.Errorf("expanding %q: configuration file exceeds the maximum of 10 match conditions, got %d", "conditions", len(conditions))
	}

	props := cdn.RuleUpdateParameters{
		RuleUpdatePropertiesParameters: &cdn.RuleUpdatePropertiesParameters{
			Actions:                 &actions,
			Conditions:              &conditions,
			MatchProcessingBehavior: matchProcessingBehaviorValue,
			Order:                   utils.Int32(int32(order)),
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontdoorRuleRead(d, meta)
}

func resourceCdnFrontdoorRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRuleID(d.Id())
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

	m := *cdnfrontdoorruleactions.InitializeCdnFrontdoorActionMappings()

	actions := map[string]expandfunc{
		m.RouteConfigurationOverride.ConfigName: cdnfrontdoorruleactions.ExpandCdnFrontdoorRouteConfigurationOverrideAction,
		m.RequestHeader.ConfigName:              cdnfrontdoorruleactions.ExpandCdnFrontdoorRequestHeaderAction,
		m.ResponseHeader.ConfigName:             cdnfrontdoorruleactions.ExpandCdnFrontdoorResponseHeaderAction,
		m.URLRedirect.ConfigName:                cdnfrontdoorruleactions.ExpandCdnFrontdoorUrlRedirectAction,
		m.URLRewrite.ConfigName:                 cdnfrontdoorruleactions.ExpandCdnFrontdoorUrlRewriteAction,
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
				return nil, fmt.Errorf("the %q is only allow once in the %q match block, got %d", m.URLRewrite.ConfigName, "actions", len(*expanded))
			}

			if actionName == m.URLRedirect.ConfigName && len(*expanded) > 1 {
				return nil, fmt.Errorf("the %q is only allow once in the %q match block, got %d", m.URLRedirect.ConfigName, "actions", len(*expanded))
			}

			if actionName == m.RouteConfigurationOverride.ConfigName && len(*expanded) > 1 {
				return nil, fmt.Errorf("the %q is only allow once in the %q match block, got %d", m.RouteConfigurationOverride.ConfigName, "actions", len(*expanded))
			}

			results = append(results, *expanded...)
		}
	}

	if len(results) > 5 {
		return nil, fmt.Errorf("the %q match block may only contain upto 5 match actions, got %d", "actions", len(results))
	}

	// validate action block
	if err := validateActionsBlock(results); err != nil {
		return nil, err
	}

	return results, nil
}

func expandFrontdoorDeliveryRuleConditions(input []interface{}) ([]cdn.BasicDeliveryRuleCondition, error) {
	results := make([]cdn.BasicDeliveryRuleCondition, 0)
	if len(input) == 0 {
		return results, nil
	}

	type expandfunc func(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error)
	m := cdnfrontdoorruleconditions.InitializeCdnFrontdoorConditionMappings()

	// TODO: Add validation for this error once I figure out what that validation should be
	// "BadRequest" Message="Parameters specified for RequestMethod condition are invalid.;
	// Property 'Rule.Conditions[2].Parameters.Selector' is required but it was not set"

	conditions := map[string]expandfunc{
		m.ClientPort.ConfigName:       cdnfrontdoorruleconditions.ExpandCdnFrontdoorClientPortCondition,
		m.Cookies.ConfigName:          cdnfrontdoorruleconditions.ExpandCdnFrontdoorCookiesCondition,
		m.HostName.ConfigName:         cdnfrontdoorruleconditions.ExpandCdnFrontdoorHostNameCondition,
		m.HttpVersion.ConfigName:      cdnfrontdoorruleconditions.ExpandCdnFrontdoorHttpVersionCondition,
		m.IsDevice.ConfigName:         cdnfrontdoorruleconditions.ExpandCdnFrontdoorIsDeviceCondition,
		m.PostArgs.ConfigName:         cdnfrontdoorruleconditions.ExpandCdnFrontdoorPostArgsCondition,
		m.QueryString.ConfigName:      cdnfrontdoorruleconditions.ExpandCdnFrontdoorQueryStringCondition,
		m.RemoteAddress.ConfigName:    cdnfrontdoorruleconditions.ExpandCdnFrontdoorRemoteAddressCondition,
		m.RequestBody.ConfigName:      cdnfrontdoorruleconditions.ExpandCdnFrontdoorRequestBodyCondition,
		m.RequestHeader.ConfigName:    cdnfrontdoorruleconditions.ExpandCdnFrontdoorRequestHeaderCondition,
		m.RequestMethod.ConfigName:    cdnfrontdoorruleconditions.ExpandCdnFrontdoorRequestMethodCondition,
		m.RequestScheme.ConfigName:    cdnfrontdoorruleconditions.ExpandCdnFrontdoorRequestSchemeCondition,
		m.RequestUri.ConfigName:       cdnfrontdoorruleconditions.ExpandCdnFrontdoorRequestUriCondition,
		m.ServerPort.ConfigName:       cdnfrontdoorruleconditions.ExpandCdnFrontdoorServerPortCondition,
		m.SocketAddress.ConfigName:    cdnfrontdoorruleconditions.ExpandCdnFrontdoorSocketAddressCondition,
		m.SslProtocol.ConfigName:      cdnfrontdoorruleconditions.ExpandCdnFrontdoorSslProtocolCondition,
		m.UrlFileExtension.ConfigName: cdnfrontdoorruleconditions.ExpandCdnFrontdoorUrlFileExtensionCondition,
		m.UrlFilename.ConfigName:      cdnfrontdoorruleconditions.ExpandCdnFrontdoorUrlFileNameCondition,
		m.UrlPath.ConfigName:          cdnfrontdoorruleconditions.ExpandCdnFrontdoorUrlPathCondition,
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
		return nil, fmt.Errorf("the %q match block may only contain upto 10 match conditions, got %d", "conditions", len(results))
	}

	return results, nil
}

func flattenFrontdoorDeliveryRuleConditions(input *[]cdn.BasicDeliveryRuleCondition) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	m := cdnfrontdoorruleconditions.InitializeCdnFrontdoorConditionMappings()
	keys := make(map[string]string)

	for _, BasicDeliveryRuleCondition := range *input {
		result := make(map[string]interface{})

		// Client Port
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleClientPortCondition(); ok {
			conditionMapping := m.ClientPort

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorClientPortCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Cookies
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleCookiesCondition(); ok {
			conditionMapping := m.Cookies

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorCookiesCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Host Name
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleHostNameCondition(); ok {
			conditionMapping := m.HostName

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorHostNameCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// HTTP Version
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleHTTPVersionCondition(); ok {
			conditionMapping := m.HttpVersion

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorHttpVersionCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Is Device
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleIsDeviceCondition(); ok {
			conditionMapping := m.IsDevice

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorIsDeviceCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Post Args
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRulePostArgsCondition(); ok {
			conditionMapping := m.PostArgs

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorPostArgsCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Query String
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleQueryStringCondition(); ok {
			conditionMapping := m.QueryString

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorQueryStringCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Remote Address
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRemoteAddressCondition(); ok {
			conditionMapping := m.RemoteAddress

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorRemoteAddressCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Request Body
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestBodyCondition(); ok {
			conditionMapping := m.RequestBody

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorRequestBodyCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Request Header
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestHeaderCondition(); ok {
			conditionMapping := m.RequestHeader

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorRequestHeaderCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Request Method
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestMethodCondition(); ok {
			conditionMapping := m.RequestMethod

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorRequestMethodCondition(condition, m.RequestMethod)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Request Scheme
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestSchemeCondition(); ok {
			conditionMapping := m.RequestScheme

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorRequestSchemeCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Request URI
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestURICondition(); ok {
			conditionMapping := m.RequestUri

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorRequestUriCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Server Port
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleServerPortCondition(); ok {
			conditionMapping := m.ServerPort

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorServerPortCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Socket Address
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleSocketAddrCondition(); ok {
			conditionMapping := m.SocketAddress

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorSocketAddressCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Ssl Protocol
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleSslProtocolCondition(); ok {
			conditionMapping := m.SslProtocol

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorSslProtocolCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// URL File Extension
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLFileExtensionCondition(); ok {
			conditionMapping := m.UrlFileExtension

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorUrlFileExtensionCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// URL Filename
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLFileNameCondition(); ok {
			conditionMapping := m.UrlFilename

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorUrlFileNameCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// URL Path
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLPathCondition(); ok {
			conditionMapping := m.UrlPath

			flattened, err := cdnfrontdoorruleconditions.FlattenFrontdoorUrlPathCondition(condition, conditionMapping)
			if err != nil {
				return nil, err
			}

			if flattened != nil {
				result[conditionMapping.ConfigName] = flattened
				keys[conditionMapping.ConfigName] = conditionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		return nil, fmt.Errorf("unknown BasicDeliveryRuleCondition encountered")
	}

	output := make(map[string][]interface{})

	// at this point our keys map contains the list of all conditions that were flattened
	if len(keys) > 0 {
		// set up a bucket to hold the flattened resource
		for key := range keys {
			output[key] = make([]interface{}, 0)
		}

		// now loop over all of the flattened
		// conditions and add them to the
		// right bucket in output
		for _, conditions := range results {
			condition := conditions.(map[string]interface{})

			for key := range condition {
				output[key] = append(output[key], condition[key])
			}
		}

		// log.Printf("\n\n\n\n\n\n\n\n\n\nXXXX-XX-XXTXX:XX:XX.XXX-0700 [DEBUG] plugin.terraform-provider-azurerm: Flatten Conditions ***************************************************\n")
		// log.Printf("  output  == %+v\n", output)
		// log.Printf(" \n\n\n")
		// log.Printf("  results == %+v\n\n", results)
		// log.Printf("***************************************************\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")

		return []interface{}{output}, nil
	}

	return results, nil
}

func flattenFrontdoorDeliveryRuleActions(input *[]cdn.BasicDeliveryRuleAction) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	m := cdnfrontdoorruleactions.InitializeCdnFrontdoorActionMappings()
	keys := make(map[string]string)

	for _, BasicDeliveryRuleAction := range *input {
		result := make(map[string]interface{})

		// Route Configuration
		if action, ok := BasicDeliveryRuleAction.AsDeliveryRuleRouteConfigurationOverrideAction(); ok {
			flattened, err := cdnfrontdoorruleactions.FlattenCdnFrontdoorRouteConfigurationOverrideAction(action)
			if err != nil {
				return nil, err
			}

			actionMapping := m.RouteConfigurationOverride

			if len(flattened) > 0 {
				result[actionMapping.ConfigName] = flattened
				keys[actionMapping.ConfigName] = actionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Request Header
		if action, ok := BasicDeliveryRuleAction.AsDeliveryRuleRequestHeaderAction(); ok {
			flattened, err := cdnfrontdoorruleactions.FlattenCdnFrontdoorRequestHeaderAction(action)
			if err != nil {
				return nil, err
			}

			actionMapping := m.RequestHeader

			if len(flattened) > 0 {
				result[actionMapping.ConfigName] = flattened
				keys[actionMapping.ConfigName] = actionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// Response Header
		if action, ok := BasicDeliveryRuleAction.AsDeliveryRuleResponseHeaderAction(); ok {
			flattened, err := cdnfrontdoorruleactions.FlattenCdnFrontdoorResponseHeaderAction(action)
			if err != nil {
				return nil, err
			}

			actionMapping := m.ResponseHeader

			if len(flattened) > 0 {
				result[actionMapping.ConfigName] = flattened
				keys[actionMapping.ConfigName] = actionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// URL Redirect
		if action, ok := BasicDeliveryRuleAction.AsURLRedirectAction(); ok {
			flattened, err := cdnfrontdoorruleactions.FlattenCdnFrontdoorUrlRedirectAction(action)
			if err != nil {
				return nil, err
			}

			actionMapping := m.URLRedirect

			if len(flattened) > 0 {
				result[actionMapping.ConfigName] = flattened
				keys[actionMapping.ConfigName] = actionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		// URL Rewrite
		if action, ok := BasicDeliveryRuleAction.AsURLRewriteAction(); ok {
			flattened, err := cdnfrontdoorruleactions.FlattenCdnFrontdoorUrlRewriteAction(action)
			if err != nil {
				return nil, err
			}

			actionMapping := m.URLRewrite

			if len(flattened) > 0 {
				result[actionMapping.ConfigName] = flattened
				keys[actionMapping.ConfigName] = actionMapping.ConfigName
				results = append(results, result)
			}

			continue
		}

		return nil, fmt.Errorf("unknown BasicDeliveryRuleAction encountered")
	}

	output := make(map[string][]interface{})

	// at this point our keys map contains the list of all actions that were flattened
	if len(keys) > 0 {
		// set up a bucket to hold the flattened resource
		for key := range keys {
			output[key] = make([]interface{}, 0)
		}

		// now loop over all of the flattened
		// actions and add them to the
		// right bucket in output
		for _, actions := range results {
			action := actions.(map[string]interface{})

			for key := range action {
				output[key] = append(output[key], action[key])
			}
		}

		// log.Printf("\n\n\n\n\n\n\n\n\n\nXXXX-XX-XXTXX:XX:XX.XXX-0700 [DEBUG] plugin.terraform-provider-azurerm: Flatten Actions ***************************************************\n")
		// log.Printf("  output  == %+v\n", output)
		// log.Printf(" \n\n\n")
		// log.Printf("  results == %+v\n\n", results)
		// log.Printf("***************************************************\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")

		return []interface{}{output}, nil
	}

	// this should be an empty interface at this point
	return results, nil
}
