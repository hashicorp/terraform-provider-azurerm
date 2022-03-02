package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorRuleCreate,
		Read:   resourceFrontdoorRuleRead,
		Update: resourceFrontdoorRuleUpdate,
		Delete: resourceFrontdoorRuleDelete,

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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontdoor_rule_set_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorRuleSetID,
			},

			"match_processing_behavior": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(track1.MatchProcessingBehaviorContinue),
				ValidateFunc: validation.StringInSlice([]string{
					string(track1.MatchProcessingBehaviorContinue),
					string(track1.MatchProcessingBehaviorStop),
				}, false),
			},

			"order": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			// I don't need name as a field I can derive the correct value based on what
			// type of parameter you define in the config

			// string(track1.NameBasicDeliveryRuleActionNameCacheExpiration),
			// string(track1.NameBasicDeliveryRuleActionNameCacheKeyQueryString),
			// string(track1.NameBasicDeliveryRuleActionNameModifyRequestHeader),
			// string(track1.NameBasicDeliveryRuleActionNameModifyResponseHeader),
			// string(track1.NameBasicDeliveryRuleActionNameOriginGroupOverride),
			// string(track1.NameBasicDeliveryRuleActionNameRouteConfigurationOverride),
			// string(track1.NameBasicDeliveryRuleActionNameURLRedirect),
			// string(track1.NameBasicDeliveryRuleActionNameURLRewrite),
			// string(track1.NameBasicDeliveryRuleActionNameURLSigning),

			// type BasicDeliveryRuleAction interface {
			"actions": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 5,

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
											string(track1.RedirectTypeMoved),
											string(track1.RedirectTypeFound),
											string(track1.RedirectTypeTemporaryRedirect),
											string(track1.RedirectTypePermanentRedirect),
										}, false),
									},

									"redirect_protocol": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(track1.DestinationProtocolMatchRequest),
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.DestinationProtocolMatchRequest),
											string(track1.DestinationProtocolHTTP),
											string(track1.DestinationProtocolHTTPS),
										}, false),
									},

									// TODO: Write validation function for this
									// Path cannot be empty and must start with /. Leave empty to use the incoming path as destination path.
									"destination_path": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"destination_hostname": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									// TODO: Write validation function for this
									// Query string must be in <key>=<value> format. ? and & will be added automatically so do not include them.
									"query_string": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"destination_fragment": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"url_signing_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"algorithm": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(track1.AlgorithmSHA256),
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.AlgorithmSHA256),
										}, false),
									},

									"parameter_name_override": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MaxItems: 100,

										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{

												"param_name": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"param_indicator": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(track1.ParamIndicatorExpires),
														string(track1.ParamIndicatorKeyID),
														string(track1.ParamIndicatorSignature),
													}, false),
												},
											},
										},
									},
								},
							},
						},

						"origin_group_override_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"origin_group_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validate.FrontdoorOriginGroupID,
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
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"header_action": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.HeaderActionAppend),
											string(track1.HeaderActionOverwrite),
											string(track1.HeaderActionDelete),
										}, false),
									},

									"header_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"value": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"response_header_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"header_action": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.HeaderActionAppend),
											string(track1.HeaderActionOverwrite),
											string(track1.HeaderActionDelete),
										}, false),
									},

									"header_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"value": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"cache_expiration_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"cache_behavior": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(track1.CacheBehaviorSetIfMissing),
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.CacheBehaviorBypassCache),
											string(track1.CacheBehaviorOverride),
											string(track1.CacheBehaviorSetIfMissing),
										}, false),
									},

									"cache_type": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  "All",
										ValidateFunc: validation.StringInSlice([]string{
											"All",
										}, false),
									},

									// Allowed format is d.hh:mm:ss, maximum duration is 366 days
									"cache_duration": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: ValidateFrontdoorCacheDuration,
									},
								},
							},
						},

						"cache_key_query_string_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"query_string_behavior": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(track1.QueryStringBehaviorExcludeAll),
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.QueryStringBehaviorInclude),
											string(track1.QueryStringBehaviorIncludeAll),
											string(track1.QueryStringBehaviorExclude),
											string(track1.QueryStringBehaviorExcludeAll),
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
								},
							},
						},
					},
				},
			},

			// type BasicDeliveryRuleCondition interface {
			"conditions": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 10,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"remote_addr_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.RemoteAddressOperatorAny),
											string(track1.RemoteAddressOperatorIPMatch),
											string(track1.RemoteAddressOperatorGeoMatch),
										}, false),
									},

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),
								},
							},
						},

						"request_method_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperatorEqualOnly(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorRequestMethodMatchValues(),
								},
							},
						},

						"query_string_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"postargs_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									// In the API this is called selector
									"postargs_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"request_uri_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"request_header_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									// In the API this is called selector
									// match_values are optional if operator is Any
									"header_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"request_body_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValuesRequired(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"request_scheme_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperatorEqualOnly(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorProtocolMatchValues(),
								},
							},
						},

						"url_path_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorUrlPathConditionMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"url_extension_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValuesRequired(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"url_filename_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValuesRequired(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"http_version_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperatorEqualOnly(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorHttpVersionMatchValues(),
								},
							},
						},

						"cookies_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									// In the API this is called selector
									"cookie_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						"is_device_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperatorEqualOnly(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorIsDeviceMatchValues(),
								},
							},
						},

						// DeliveryRuleSocketAddrCondition
						"socket_addr_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.SocketAddrOperatorAny),
											string(track1.SocketAddrOperatorIPMatch),
										}, false),
									},

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						// DeliveryRuleClientPortCondition
						"client_port_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						// DeliveryRuleServerPortCondition
						"server_port_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						// DeliveryRuleHostNameCondition
						"host_name_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": SchemaFrontdoorOperator(),

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": SchemaFrontdoorMatchValues(),

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},

						// DeliveryRuleSslProtocolCondition
						"ssl_protocol_condition": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"operator": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"negate_condition": SchemaFrontdoorNegateCondition(),

									"match_values": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 3,

										Elem: &pluginsdk.Schema{
											Type:    pluginsdk.TypeString,
											Default: string(track1.SslProtocolTLSv12),
											ValidateFunc: validation.StringInSlice([]string{
												string(track1.SslProtocolTLSv1),
												string(track1.SslProtocolTLSv11),
												string(track1.SslProtocolTLSv12),
											}, false),
										},
									},

									"transforms": SchemaFrontdoorRuleTransforms(),
								},
							},
						},
						//
						//
					},
				},
			},

			"rule_set_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFrontdoorRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	ruleSetId, err := parse.FrontdoorRuleSetID(d.Get("frontdoor_rule_set_id").(string))
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
			return tf.ImportAsExistsError("azurerm_frontdoor_rule", id.ID())
		}
	}

	matchProcessingBehaviorValue := track1.MatchProcessingBehavior(d.Get("match_processing_behavior").(string))
	props := track1.Rule{
		RuleProperties: &track1.RuleProperties{
			Actions:                 expandOptionalRuleDeliveryRuleActionArray(d.Get("actions").([]interface{})),
			Conditions:              expandRuleDeliveryRuleConditionArray(d.Get("conditions").([]interface{})),
			MatchProcessingBehavior: matchProcessingBehaviorValue,
			RuleSetName:             &ruleSetId.RuleSetName,
			Order:                   utils.Int32(int32(d.Get("order").(int))),
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

	return resourceFrontdoorRuleRead(d, meta)
}

func resourceFrontdoorRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RuleName)
	d.Set("frontdoor_rule_set_id", parse.NewFrontdoorRuleSetID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.RuleSetName).ID())

	if props := resp.RuleProperties; props != nil {
		d.Set("match_processing_behavior", props.MatchProcessingBehavior)
		d.Set("order", props.Order)
		d.Set("rule_set_name", props.RuleSetName)

		if err := d.Set("actions", flattenRuleDeliveryRuleActionArray(props.Actions)); err != nil {
			return fmt.Errorf("setting `actions`: %+v", err)
		}

		if err := d.Set("conditions", flattenRuleDeliveryRuleConditionArray(props.Conditions)); err != nil {
			return fmt.Errorf("setting `conditions`: %+v", err)
		}
	}

	return nil
}

func resourceFrontdoorRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRuleID(d.Id())
	if err != nil {
		return err
	}

	matchProcessingBehaviorValue := track1.MatchProcessingBehavior(d.Get("match_processing_behavior").(string))
	props := track1.RuleUpdateParameters{
		RuleUpdatePropertiesParameters: &track1.RuleUpdatePropertiesParameters{
			Actions:                 expandRequiredRuleDeliveryRuleActionArray(d.Get("actions").([]interface{})),
			Conditions:              expandRuleDeliveryRuleConditionArray(d.Get("conditions").([]interface{})),
			MatchProcessingBehavior: matchProcessingBehaviorValue,
			Order:                   utils.Int32(int32(d.Get("order").(int))),
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceFrontdoorRuleRead(d, meta)
}

func resourceFrontdoorRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRulesClient
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

func expandRuleDeliveryRuleConditionArray(input []interface{}) *[]track1.BasicDeliveryRuleCondition {
	results := make([]track1.BasicDeliveryRuleCondition, 0)
	if len(input) == 0 {
		return &results
	}

	conditions := utils.ExpandStringSlice(input)

	for _, condition := range *conditions {

		switch condition {
		case string(track1.NameSslProtocol):
			var drspc track1.DeliveryRuleSslProtocolCondition
			results = append(results, drspc)
		case string(track1.NameHostName):
			var drhnc track1.DeliveryRuleHostNameCondition
			results = append(results, drhnc)
		case string(track1.NameServerPort):
			var drspc track1.DeliveryRuleServerPortCondition
			results = append(results, drspc)
		case string(track1.NameClientPort):
			var drcpc track1.DeliveryRuleClientPortCondition
			results = append(results, drcpc)
		case string(track1.NameSocketAddr):
			var drsac track1.DeliveryRuleSocketAddrCondition
			results = append(results, drsac)
		case string(track1.NameIsDevice):
			var dridc track1.DeliveryRuleIsDeviceCondition
			results = append(results, dridc)
		case string(track1.NameCookies):
			var drcc track1.DeliveryRuleCookiesCondition
			results = append(results, drcc)
		case string(track1.NameHTTPVersion):
			var drhvc track1.DeliveryRuleHTTPVersionCondition
			results = append(results, drhvc)
		case string(track1.NameURLFileName):
			var drufnc track1.DeliveryRuleURLFileNameCondition
			results = append(results, drufnc)
		case string(track1.NameURLFileExtension):
			var drufec track1.DeliveryRuleURLFileExtensionCondition
			results = append(results, drufec)
		case string(track1.NameURLPath):
			var drupc track1.DeliveryRuleURLPathCondition
			results = append(results, drupc)
		case string(track1.NameRequestScheme):
			var drrsc track1.DeliveryRuleRequestSchemeCondition
			results = append(results, drrsc)
		case string(track1.NameRequestBody):
			var drrbc track1.DeliveryRuleRequestBodyCondition
			results = append(results, drrbc)
		case string(track1.NameRequestHeader):
			var drrhc track1.DeliveryRuleRequestHeaderCondition
			results = append(results, drrhc)
		case string(track1.NameRequestURI):
			var drruc track1.DeliveryRuleRequestURICondition
			results = append(results, drruc)
		case string(track1.NamePostArgs):
			var drpac track1.DeliveryRulePostArgsCondition
			results = append(results, drpac)
		case string(track1.NameQueryString):
			var drqsc track1.DeliveryRuleQueryStringCondition
			results = append(results, drqsc)
		case string(track1.NameRequestMethod):
			var drrmc track1.DeliveryRuleRequestMethodCondition
			results = append(results, drrmc)
		case string(track1.NameRemoteAddress):
			var drrac track1.DeliveryRuleRemoteAddressCondition
			results = append(results, drrac)
		default:
			var drc track1.DeliveryRuleCondition
			results = append(results, drc)
		}

		results = append(results, track1.DeliveryRuleCondition{
			Name: track1.Name(condition),
		})
	}

	return &results
}

// TODO: Get this actually working with new schema
func expandOptionalRuleDeliveryRuleActionArray(input []interface{}) *[]track1.BasicDeliveryRuleAction {
	result := make([]track1.BasicDeliveryRuleAction, 0)
	if len(input) == 0 {

		return &result
	}

	result = append(result, track1.BasicDeliveryRuleAction(expandRuleDeliveryRuleActions(input)))
	return &result
}

func expandRequiredRuleDeliveryRuleActionArray(input []interface{}) *[]track1.BasicDeliveryRuleAction {
	results := make([]track1.BasicDeliveryRuleAction, 0)
	if len(input) == 0 {
		return nil
	}

	// results := expandRuleDeliveryRuleActions(input)

	return &results
}

func expandRuleDeliveryRuleActions(input []interface{}) track1.BasicDeliveryRuleAction {
	// results := make([]track1.BasicDeliveryRuleAction, 0)
	// actions := utils.ExpandStringSlice(input)

	// for _, action := range *actions {
	// 	results = append(results, track1.DeliveryRuleAction{
	// 		Name: track1.NameBasicDeliveryRuleAction(action),
	// 	})
	// }

	return nil
}

func flattenRuleDeliveryRuleConditionArray(input *[]track1.BasicDeliveryRuleCondition) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// for _, condition := range *input {
	// 	results = append(results, condition.Name)
	// }

	return results
}

func flattenRuleDeliveryRuleActionArray(input *[]track1.BasicDeliveryRuleAction) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// for _, action := range *input {
	// 	results = append(results, string(action))
	// }

	return results
}
