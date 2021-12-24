package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAfdRuleSetRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdRuleSetRuleCreate,
		Read:   resourceAfdRuleSetRuleRead,
		Update: resourceAfdRuleSetRuleUpdate,
		Delete: resourceAfdRuleSetRuleDelete,

		SchemaVersion: 1,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AfdRuleSetRulesID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"rule_set_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AfdRuleSetsID,
			},

			"order": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			"action": {
				Type:     pluginsdk.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.NameBasicDeliveryRuleActionNameURLRedirect),
								string(cdn.NameBasicDeliveryRuleActionNameURLRewrite),
								string(cdn.NameBasicDeliveryRuleActionNameURLSigning),
							}, false),
						},
						// RedirectType - The redirect type the rule will use when redirecting traffic. Possible values include: 'RedirectTypeMoved', 'RedirectTypeFound', 'RedirectTypeTemporaryRedirect', 'RedirectTypePermanentRedirect'
						"redirect_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.RedirectTypeFound),
								string(cdn.RedirectTypeMoved),
								string(cdn.RedirectTypePermanentRedirect),
								string(cdn.RedirectTypeTemporaryRedirect),
							}, false),
						},
						// DestinationProtocol - Protocol to use for the redirect. The default value is MatchRequest. Possible values include: 'DestinationProtocolMatchRequest', 'DestinationProtocolHTTP', 'DestinationProtocolHTTPS'
						"destination_protocol": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  cdn.DestinationProtocolMatchRequest,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.DestinationProtocolMatchRequest),
								string(cdn.DestinationProtocolHTTP),
								string(cdn.DestinationProtocolHTTPS),
							}, false),
						},
						// CustomPath - The full path to redirect. Path cannot be empty and must start with /. Leave empty to use the incoming path as destination path.
						"custom_path": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						// CustomHostname - Host to redirect. Leave empty to use the incoming host as the destination host.
						"custom_hostname": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						// CustomQueryString - The set of query strings to be placed in the redirect URL. Setting this value would replace any existing query string; leave empty to preserve the incoming query string. Query string must be in <key>=<value> format. ? and & will be added automatically so do not include them.
						"custom_querystring": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						// CustomFragment - Fragment to add to the redirect URL. Fragment is the part of the URL that comes after #. Do not include the #.
						"custom_fragment": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						// AlgorithmSHA256
						"algorithm": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  cdn.AlgorithmSHA256,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.AlgorithmSHA256),
							}, false),
						},

						// URLSigningParamIdentifier defines how to identify a parameter for a specific purpose e.g. expires
						"param_indicator": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.ParamIndicatorExpires),
								string(cdn.ParamIndicatorKeyID),
								string(cdn.ParamIndicatorSignature),
							}, false),
						},
						"param_name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"destination": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"source_pattern": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"preserve_unmatched_path": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"condition": {
				Type:     pluginsdk.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"sample_size": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  4,
						},
					},
				},
			},
		},
	}
}

func resourceAfdRuleSetRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDRuleSetRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// rule name
	ruleName := d.Get("name").(string)

	// parse rule_set_id
	ruleSetId := d.Get("rule_set_id").(string)
	ruleSet, err := parse.AfdRuleSetsID(ruleSetId)
	if err != nil {
		return err
	}

	id := parse.NewAfdRuleSetRulesID(ruleSet.SubscriptionId, ruleSet.ResourceGroup, ruleSet.ProfileName, ruleSet.RuleSetName, ruleName)

	rule := cdn.Rule{}
	ruleProperties := cdn.RuleProperties{}

	// order
	order := int32(d.Get("order").(int))
	ruleProperties.Order = utils.Int32(order)

	// conditions
	ruleProperties.Conditions = expandRuleConditions(d.Get("condition").([]interface{}))

	// actions
	ruleProperties.Actions = expandRuleActions(d.Get("action").([]interface{}))

	rule.RuleProperties = &ruleProperties

	future, err := client.Create(ctx, ruleSet.ResourceGroup, ruleSet.ProfileName, ruleSet.RuleSetName, ruleName, rule)
	if err != nil {
		return fmt.Errorf("creating rule %s in rule set %s: %+v", ruleName, ruleSet.RuleSetName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of rule %s in rule set %s: %+v", ruleName, ruleSet.RuleSetName, err)
	}

	d.SetId(id.ID())

	return resourceAfdRuleSetRuleRead(d, meta)
}

func expandRuleConditions(conditions []interface{}) *[]cdn.BasicDeliveryRuleCondition {

	return nil
}

func expandRuleActions(actions []interface{}) *[]cdn.BasicDeliveryRuleAction {

	results := make([]cdn.BasicDeliveryRuleAction, 0)

	for _, action := range actions {

		config := action.(map[string]interface{})
		actionType := config["type"].(string)

		switch actionType {
		case string(cdn.NameBasicDeliveryRuleActionNameURLRedirect):
			ruleActionType := cdn.URLRedirectAction{}
			ruleActionType.Name = cdn.NameBasicDeliveryRuleAction(actionType)

			params := cdn.URLRedirectActionParameters{}

			if customFragment := config["custom_fragment"].(string); customFragment != "" {
				params.CustomFragment = &customFragment
			}

			if customHostname := config["custom_hostname"].(string); customHostname != "" {
				params.CustomHostname = &customHostname
			}

			if customPath := config["custom_path"].(string); customPath != "" {
				params.CustomPath = &customPath
			}

			if customQueryString := config["custom_querystring"].(string); customQueryString != "" {
				params.CustomQueryString = &customQueryString
			}

			if destinationProtocol := config["destination_protocol"].(string); destinationProtocol != "" {
				params.DestinationProtocol = cdn.DestinationProtocol(destinationProtocol)
			}

			if redirectType := config["redirect_type"].(string); redirectType != "" {
				params.RedirectType = cdn.RedirectType(redirectType)
			}

			params.OdataType = utils.String("#Microsoft.Azure.Cdn.Models.DeliveryRuleUrlRedirectActionParameters")

			ruleActionType.Parameters = &params

			results = append(results, ruleActionType)

		case string(cdn.NameBasicDeliveryRuleActionNameURLSigning):
			ruleActionType := cdn.URLSigningAction{}
			params := cdn.URLSigningActionParameters{}

			if algorithm := config["algorithm"].(string); algorithm != "" {
				params.Algorithm = cdn.Algorithm(algorithm)
			}

			paramIdentifierArray := make([]cdn.URLSigningParamIdentifier, 0)
			paramIdentifier := cdn.URLSigningParamIdentifier{}

			if paramName := config["param_name"].(string); paramName != "" {
				paramIdentifier.ParamName = &paramName
			}

			if paramIndicator := config["param_indicator"].(string); paramIndicator != "" {
				paramIdentifier.ParamIndicator = cdn.ParamIndicator(paramIndicator)
			}

			params.ParameterNameOverride = &paramIdentifierArray

			params.OdataType = utils.String("#Microsoft.Azure.Cdn.Models.DeliveryRuleUrlSigningActionParameters")

			ruleActionType.Parameters = &params

			results = append(results, ruleActionType)

		case string(cdn.NameBasicDeliveryRuleActionNameURLRewrite):
			ruleActionType := cdn.URLRewriteAction{}
			params := cdn.URLRewriteActionParameters{}
			ruleActionType.Parameters = &params

			if destination := config["destination"].(string); destination != "" {
				params.Destination = &destination
			}

			preserveUnmatchedPath := config["preserve_unmatched_path"].(bool)
			params.PreserveUnmatchedPath = &preserveUnmatchedPath

			if sourcePattern := config["source_pattern"].(string); sourcePattern != "" {
				params.SourcePattern = &sourcePattern
			}

			params.OdataType = utils.String("#Microsoft.Azure.Cdn.Models.DeliveryRuleUrlRewriteActionParameters")

			results = append(results, ruleActionType)

		default:
			log.Fatalf("%s is not implemented, yet.", string(cdn.NameBasicDeliveryRuleActionNameURLSigning))

		}
	}

	return &results
}

func resourceAfdRuleSetRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDRuleSetRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdRuleSetRulesID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("action", resp.Actions)
	d.Set("condition", resp.Conditions)

	d.Set("name", resp.Name)

	return nil
}

func resourceAfdRuleSetRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDRuleSetRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdRuleSetRulesID(d.Id())
	if err != nil {
		return err
	}

	ruleUpdates := cdn.RuleUpdateParameters{}
	ruleUpdateProperties := cdn.RuleUpdatePropertiesParameters{}

	if d.HasChange("order") {
		order := int32(d.Get("order").(int))
		ruleUpdateProperties.Order = utils.Int32(order)
	}

	// update conditions
	if d.HasChange("condition") {
		ruleUpdateProperties.Conditions = expandRuleConditions(d.Get("condition").([]interface{}))
	}

	// update actions
	if d.HasChange("action") {
		ruleUpdateProperties.Actions = expandRuleActions(d.Get("action").([]interface{}))
	}

	ruleUpdates.RuleUpdatePropertiesParameters = &ruleUpdateProperties

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName, ruleUpdates)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceAfdRuleSetRuleRead(d, meta)
}

func resourceAfdRuleSetRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDRuleSetRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdRuleSetRulesID(d.Id())
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

	return err
}
