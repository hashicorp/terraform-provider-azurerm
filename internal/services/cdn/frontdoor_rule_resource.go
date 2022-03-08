package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorruleactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/frontdoorruleconditions"
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateFrontdoorRuleName,
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
			"actions": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 5,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						// Name: UrlRedirect
						// DeliveryRuleUrlRedirectActionParameters
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

						// Name: URLRewrite
						// URLRewriteAction
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

						// Name: ModifyRequestHeader
						// DeliveryRuleRequestHeaderAction
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
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									"origin_group_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validate.FrontdoorOriginGroupID,
									},

									"forwarding_protocol": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(track1.ForwardingProtocolMatchRequest),
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.ForwardingProtocolHTTPOnly),
											string(track1.ForwardingProtocolHTTPSOnly),
											string(track1.ForwardingProtocolMatchRequest),
										}, false),
									},

									"query_string_caching_behavior": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(track1.RuleQueryStringCachingBehaviorIgnoreQueryString),
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.RuleQueryStringCachingBehaviorIgnoreQueryString),
											string(track1.RuleQueryStringCachingBehaviorUseQueryString),
											string(track1.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
											string(track1.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
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
										Default:  string(track1.RuleCacheBehaviorHonorOrigin),
										ValidateFunc: validation.StringInSlice([]string{
											string(track1.RuleCacheBehaviorHonorOrigin),
											string(track1.RuleCacheBehaviorOverrideAlways),
											string(track1.RuleCacheBehaviorOverrideIfOriginMissing),
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
					}, // end of Actions
				},
			},

			// type BasicDeliveryRuleCondition interface {
			"conditions": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 10,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"remote_address_condition": {
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

						"url_file_extension_condition": {
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
						"socket_address_condition": {
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

			"frontdoor_rule_set_name": {
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
	actions, err := expandFrontdoorDeliveryRuleActions(d.Get("actions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "actions", err)
	}

	conditions, err := expandFrontdoorDeliveryRuleConditions(d.Get("conditions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "conditions", err)
	}

	props := track1.Rule{
		RuleProperties: &track1.RuleProperties{
			Actions:                 &actions,
			Conditions:              &conditions,
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
	d.Set("frontdoor_rule_set_id", ruleSetId.ID())

	if props := resp.RuleProperties; props != nil {
		d.Set("match_processing_behavior", props.MatchProcessingBehavior)
		d.Set("order", props.Order)

		// BUG: RuleSetName is not being returned by the API
		d.Set("frontdoor_rule_set_name", ruleSetId.RuleSetName)

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

func resourceFrontdoorRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRuleID(d.Id())
	if err != nil {
		return err
	}

	actions, err := expandFrontdoorDeliveryRuleActions(d.Get("actions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "actions", err)
	}

	conditions, err := expandFrontdoorDeliveryRuleConditions(d.Get("conditions").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding %q: %+v", "conditions", err)
	}

	matchProcessingBehaviorValue := track1.MatchProcessingBehavior(d.Get("match_processing_behavior").(string))
	props := track1.RuleUpdateParameters{
		RuleUpdatePropertiesParameters: &track1.RuleUpdatePropertiesParameters{
			Actions:                 &actions,
			Conditions:              &conditions,
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

func expandFrontdoorDeliveryRuleActions(input []interface{}) ([]track1.BasicDeliveryRuleAction, error) {
	results := make([]track1.BasicDeliveryRuleAction, 0)

	type expandfunc func(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error)

	actions := map[string]expandfunc{
		"route_configuration_override_action": frontdoorruleactions.ExpandFrontdoorRouteConfigurationOverrideAction,
		"request_header_action":               frontdoorruleactions.ExpandFrontdoorRequestHeaderAction,
		"response_header_action":              frontdoorruleactions.ExpandFrontdoorResponseHeaderAction,
		"url_redirect_action":                 frontdoorruleactions.ExpandFrontdoorUrlRedirectAction,
		"url_rewrite_action":                  frontdoorruleactions.ExpandFrontdoorUrlRewriteAction,
	}

	basicDeliveryRuleAction := input[0].(map[string]interface{})

	for actionName, expand := range actions {
		raw := basicDeliveryRuleAction[actionName].([]interface{})
		expanded, err := expand(raw)
		if err != nil {
			return nil, err
		}

		results = append(results, *expanded...)
	}

	// validate action block
	if err := validateActionsBlock(results); err != nil {
		return nil, err
	}

	return results, nil
}

func expandFrontdoorDeliveryRuleConditions(input []interface{}) ([]track1.BasicDeliveryRuleCondition, error) {
	results := make([]track1.BasicDeliveryRuleCondition, 0)

	type expandfunc func(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error)
	m := frontdoorruleconditions.InitializeFrontdoorConditionMappings()

	// TODO: Add validation for this error once I figure out what that validation should be
	// Error: creating Frontdoor Rule: (Rule Name "jclinerule" / Rule Set Name "jclineruleset" / Profile Name "jcline-profile" / Resource Group "jcline-afdx-profile"): cdn.RulesClient#Create: Failure sending request: StatusCode=400 -- Original Error: Code="BadRequest" Message="Parameters specified for RequestMethod condition are invalid.; \nProperty 'Rule.Conditions[2].Parameters.Selector' is required but it was not set"

	conditions := map[string]expandfunc{
		m.ClientPort.ConfigName:       frontdoorruleconditions.ExpandFrontdoorClientPortCondition,
		m.Cookies.ConfigName:          frontdoorruleconditions.ExpandFrontdoorCookiesCondition,
		m.HostName.ConfigName:         frontdoorruleconditions.ExpandFrontdoorHostNameCondition,
		m.HttpVersion.ConfigName:      frontdoorruleconditions.ExpandFrontdoorHttpVersionCondition,
		m.IsDevice.ConfigName:         frontdoorruleconditions.ExpandFrontdoorIsDeviceCondition,
		m.PostArgs.ConfigName:         frontdoorruleconditions.ExpandFrontdoorPostArgsCondition,
		m.QueryString.ConfigName:      frontdoorruleconditions.ExpandFrontdoorQueryStringCondition,
		m.RemoteAddress.ConfigName:    frontdoorruleconditions.ExpandFrontdoorRemoteAddressCondition,
		m.RequestBody.ConfigName:      frontdoorruleconditions.ExpandFrontdoorRequestBodyCondition,
		m.RequestHeader.ConfigName:    frontdoorruleconditions.ExpandFrontdoorRequestHeaderCondition,
		m.RequestMethod.ConfigName:    frontdoorruleconditions.ExpandFrontdoorRequestMethodCondition,
		m.RequestScheme.ConfigName:    frontdoorruleconditions.ExpandFrontdoorRequestSchemeCondition,
		m.RequestUri.ConfigName:       frontdoorruleconditions.ExpandFrontdoorRequestUriCondition,
		m.ServerPort.ConfigName:       frontdoorruleconditions.ExpandFrontdoorServerPortCondition,
		m.SocketAddress.ConfigName:    frontdoorruleconditions.ExpandFrontdoorSocketAddressCondition,
		m.SslProtocol.ConfigName:      frontdoorruleconditions.ExpandFrontdoorSslProtocolCondition,
		m.UrlFileExtension.ConfigName: frontdoorruleconditions.ExpandFrontdoorUrlFileExtensionCondition,
		m.UrlFilename.ConfigName:      frontdoorruleconditions.ExpandFrontdoorUrlFileNameCondition,
		m.UrlPath.ConfigName:          frontdoorruleconditions.ExpandFrontdoorUrlPathCondition,
	}

	// conditions := map[string]expandfunc{
	// 	"client_port_condition":        frontdoorruleconditions.ExpandFrontdoorClientPortCondition,
	// 	"cookies_condition":            frontdoorruleconditions.ExpandFrontdoorCookiesCondition,
	// 	"host_name_condition":          frontdoorruleconditions.ExpandFrontdoorHostNameCondition,
	// 	"http_version_condition":       frontdoorruleconditions.ExpandFrontdoorHttpVersionCondition,
	// 	"is_device_condition":          frontdoorruleconditions.ExpandFrontdoorIsDeviceCondition,
	// 	"postargs_condition":           frontdoorruleconditions.ExpandFrontdoorPostArgsCondition,
	// 	"query_string_condition":       frontdoorruleconditions.ExpandFrontdoorQueryStringCondition,
	// 	"remote_address_condition":     frontdoorruleconditions.ExpandFrontdoorRemoteAddressCondition,
	// 	"request_body_condition":       frontdoorruleconditions.ExpandFrontdoorRequestBodyCondition,
	// 	"request_header_condition":     frontdoorruleconditions.ExpandFrontdoorRequestHeaderCondition,
	// 	"request_method_condition":     frontdoorruleconditions.ExpandFrontdoorRequestMethodCondition,
	// 	"request_scheme_condition":     frontdoorruleconditions.ExpandFrontdoorRequestSchemeCondition,
	// 	"request_uri_condition":        frontdoorruleconditions.ExpandFrontdoorRequestUriCondition,
	// 	"server_port_condition":        frontdoorruleconditions.ExpandFrontdoorServerPortCondition,
	// 	"socket_address_condition":     frontdoorruleconditions.ExpandFrontdoorSocketAddressCondition,
	// 	"ssl_protocol_condition":       frontdoorruleconditions.ExpandFrontdoorSslProtocolCondition,
	// 	"url_file_extension_condition": frontdoorruleconditions.ExpandFrontdoorUrlFileExtensionCondition,
	// 	"url_filename_condition":       frontdoorruleconditions.ExpandFrontdoorUrlFileNameCondition,
	// 	"url_path_condition":           frontdoorruleconditions.ExpandFrontdoorUrlPathCondition,
	// }

	basicDeliveryRuleCondition := input[0].(map[string]interface{})

	for conditionName, expand := range conditions {
		raw := basicDeliveryRuleCondition[conditionName].([]interface{})
		expanded, err := expand(raw)
		if err != nil {
			return nil, err
		}

		results = append(results, *expanded...)
	}

	return results, nil
}

func flattenFrontdoorDeliveryRuleConditions(input *[]track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	m := frontdoorruleconditions.InitializeFrontdoorConditionMappings()

	for _, BasicDeliveryRuleCondition := range *input {
		result := make(map[string]interface{})
		foundRule := false

		// Client Port
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleClientPortCondition(); ok {
			foundRule = true
			flattened, err := frontdoorruleconditions.FlattenFrontdoorClientPortCondition(condition, m.ClientPort)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.ClientPort.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Cookies
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleCookiesCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorCookiesCondition(condition, m.Cookies)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.Cookies.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Host Name
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleHostNameCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorHostNameCondition(condition, m.HostName)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.HostName.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// HTTP Version
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleHTTPVersionCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorHttpVersionCondition(condition, m.HttpVersion)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.HttpVersion.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Is Device
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleIsDeviceCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorIsDeviceCondition(condition, m.IsDevice)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.IsDevice.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Post Args
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRulePostArgsCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorPostArgsCondition(condition, m.PostArgs)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.PostArgs.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Query String
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleQueryStringCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorQueryStringCondition(condition, m.QueryString)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.QueryString.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Remote Address
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRemoteAddressCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorRemoteAddressCondition(condition, m.RemoteAddress)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.RemoteAddress.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Request Body
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestBodyCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorRequestBodyCondition(condition, m.RequestBody)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.RequestBody.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Request Header
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestHeaderCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorRequestHeaderCondition(condition, m.RequestHeader)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.RequestHeader.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Request Method
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestMethodCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorRequestMethodCondition(condition, m.RequestMethod)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.RequestMethod.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Request Scheme
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestSchemeCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorRequestSchemeCondition(condition, m.RequestScheme)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.RequestScheme.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Request URI
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleRequestURICondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorRequestUriCondition(condition, m.RequestUri)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.RequestUri.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Server Port
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleServerPortCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorServerPortCondition(condition, m.ServerPort)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.ServerPort.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Socket Address
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleSocketAddrCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorSocketAddressCondition(condition, m.SocketAddress)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.SocketAddress.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// Ssl Protocol
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleSslProtocolCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorSslProtocolCondition(condition, m.SslProtocol)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.SslProtocol.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// URL File Extension
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLFileExtensionCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorUrlFileExtensionCondition(condition, m.UrlFileExtension)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.UrlFileExtension.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// URL Filename
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLFileNameCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorUrlFileNameCondition(condition, m.UrlFilename)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.UrlFilename.ConfigName] = flattened
				results = append(results, result)
			}
		}

		// URL Path
		if condition, ok := BasicDeliveryRuleCondition.AsDeliveryRuleURLPathCondition(); ok {
			foundRule = true

			flattened, err := frontdoorruleconditions.FlattenFrontdoorUrlPathCondition(condition, m.UrlPath)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result[m.UrlPath.ConfigName] = flattened
				results = append(results, result)
			}
		}

		if !foundRule {
			return nil, fmt.Errorf("unknown BasicDeliveryRuleAction encountered")
		}
	}

	log.Printf("\n\n\n\n\n\n\n\n\n\nXXXX-XX-XXTXX:XX:XX.XXX-0700 [DEBUG] plugin.terraform-provider-azurerm:***************************************************")
	log.Printf("  results == %+v", results)
	log.Printf("***************************************************\n\n\n\n\n\n\n\n\n\n")

	return results, nil
}

func flattenFrontdoorDeliveryRuleActions(input *[]track1.BasicDeliveryRuleAction) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, BasicDeliveryRuleAction := range *input {
		result := make(map[string]interface{})
		foundRule := false

		// Route Configuraton
		if action, ok := BasicDeliveryRuleAction.AsDeliveryRuleRouteConfigurationOverrideAction(); ok {
			foundRule = true
			flattened, err := frontdoorruleactions.FlattenFrontdoorRouteConfigurationOverrideAction(action)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result["route_configuration_override_action"] = flattened
				results = append(results, result)
			}
		}

		// Request Header
		if action, ok := BasicDeliveryRuleAction.AsDeliveryRuleRequestHeaderAction(); ok {
			foundRule = true

			flattened, err := frontdoorruleactions.FlattenFrontdoorRequestHeaderAction(action)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result["request_header_action"] = flattened
				results = append(results, result)
			}
		}

		// Response Header
		if action, ok := BasicDeliveryRuleAction.AsDeliveryRuleResponseHeaderAction(); ok {
			foundRule = true

			flattened, err := frontdoorruleactions.FlattenFrontdoorResponseHeaderAction(action)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result["response_header_action"] = flattened
				results = append(results, result)
			}
		}

		// URL Redirect
		if action, ok := BasicDeliveryRuleAction.AsURLRedirectAction(); ok {
			foundRule = true

			flattened, err := frontdoorruleactions.FlattenFrontdoorUrlRedirectAction(action)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result["url_redirect_action"] = flattened
				results = append(results, result)
			}
		}

		// URL Rewrite
		if action, ok := BasicDeliveryRuleAction.AsURLRewriteAction(); ok {
			foundRule = true

			flattened, err := frontdoorruleactions.FlattenFrontdoorUrlRewriteAction(action)
			if err != nil {
				return nil, err
			}

			if len(flattened) > 0 {
				result["url_rewrite_action"] = flattened
				results = append(results, result)
			}
		}

		if !foundRule {
			return nil, fmt.Errorf("unknown BasicDeliveryRuleAction encountered")
		}
	}

	log.Printf("\n\n\n\n\n\n\n\n\n\nXXXX-XX-XXTXX:XX:XX.XXX-0700 [DEBUG] plugin.terraform-provider-azurerm:***************************************************")
	log.Printf("  results == %+v", results)
	log.Printf("***************************************************\n\n\n\n\n\n\n\n\n\n")

	return results, nil
}
