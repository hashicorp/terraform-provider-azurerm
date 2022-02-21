package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/rulesets"
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
			_, err := rules.ParseRuleID(id)
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
				ValidateFunc: rulesets.ValidateRuleSetID,
			},

			"actions": {
				Type:     pluginsdk.TypeList,
				Required: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(rules.DeliveryRuleActionCacheExpiration),
						string(rules.DeliveryRuleActionCacheKeyQueryString),
						string(rules.DeliveryRuleActionModifyRequestHeader),
						string(rules.DeliveryRuleActionModifyResponseHeader),
						string(rules.DeliveryRuleActionOriginGroupOverride),
						string(rules.DeliveryRuleActionRouteConfigurationOverride),
						string(rules.DeliveryRuleActionUrlRedirect),
						string(rules.DeliveryRuleActionUrlRewrite),
						string(rules.DeliveryRuleActionUrlSigning),
					}, false),
				},
			},

			"conditions": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(rules.MatchVariableClientPort),
						string(rules.MatchVariableCookies),
						string(rules.MatchVariableHostName),
						string(rules.MatchVariableHttpVersion),
						string(rules.MatchVariableIsDevice),
						string(rules.MatchVariablePostArgs),
						string(rules.MatchVariableQueryString),
						string(rules.MatchVariableRemoteAddress),
						string(rules.MatchVariableRequestBody),
						string(rules.MatchVariableRequestHeader),
						string(rules.MatchVariableRequestMethod),
						string(rules.MatchVariableRequestScheme),
						string(rules.MatchVariableRequestUri),
						string(rules.MatchVariableServerPort),
						string(rules.MatchVariableSocketAddr),
						string(rules.MatchVariableSslProtocol),
						string(rules.MatchVariableUrlFileExtension),
						string(rules.MatchVariableUrlFileName),
						string(rules.MatchVariableUrlPath),
					}, false),
				},
			},

			"match_processing_behavior": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(rules.MatchProcessingBehaviorContinue),
				ValidateFunc: validation.StringInSlice([]string{
					string(rules.MatchProcessingBehaviorContinue),
					string(rules.MatchProcessingBehaviorStop),
				}, false),
			},

			"order": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
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

	ruleSetId, err := rulesets.ParseRuleSetID(d.Get("frontdoor_rule_set_id").(string))
	if err != nil {
		return err
	}

	id := rules.NewRuleID(ruleSetId.SubscriptionId, ruleSetId.ResourceGroupName, ruleSetId.ProfileName, ruleSetId.RuleSetName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_rule", id.ID())
		}
	}

	matchProcessingBehaviorValue := rules.MatchProcessingBehavior(d.Get("match_processing_behavior").(string))
	props := rules.Rule{
		Properties: &rules.RuleProperties{
			Actions:                 expandOptionalRuleDeliveryRuleActionArray(d.Get("actions").([]interface{})),
			Conditions:              expandRuleDeliveryRuleConditionArray(d.Get("conditions").([]interface{})),
			MatchProcessingBehavior: &matchProcessingBehaviorValue,
			RuleSetName:             &ruleSetId.RuleSetName,
			Order:                   int64(d.Get("order").(int)),
		},
	}
	if err := client.CreateThenPoll(ctx, id, props); err != nil {

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorRuleRead(d, meta)
}

func resourceFrontdoorRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RuleName)
	d.Set("frontdoor_rule_set_id", rulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("match_processing_behavior", props.MatchProcessingBehavior)
			d.Set("order", props.Order)
			d.Set("rule_set_name", props.RuleSetName)

			if err := d.Set("actions", flattenRuleDeliveryRuleActionArray(&props.Actions)); err != nil {
				return fmt.Errorf("setting `actions`: %+v", err)
			}

			if err := d.Set("conditions", flattenRuleDeliveryRuleConditionArray(props.Conditions)); err != nil {
				return fmt.Errorf("setting `conditions`: %+v", err)
			}
		}
	}

	return nil
}

func resourceFrontdoorRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseRuleID(d.Id())
	if err != nil {
		return err
	}

	matchProcessingBehaviorValue := rules.MatchProcessingBehavior(d.Get("match_processing_behavior").(string))
	props := rules.RuleUpdateParameters{
		Properties: &rules.RuleUpdatePropertiesParameters{
			Actions:                 expandRequiredRuleDeliveryRuleActionArray(d.Get("actions").([]interface{})),
			Conditions:              expandRuleDeliveryRuleConditionArray(d.Get("conditions").([]interface{})),
			MatchProcessingBehavior: &matchProcessingBehaviorValue,
			Order:                   utils.Int64(int64(d.Get("order").(int))),
		},
	}
	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorRuleRead(d, meta)
}

func resourceFrontdoorRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseRuleID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {

		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandRuleDeliveryRuleConditionArray(input []interface{}) *[]rules.DeliveryRuleCondition {
	results := make([]rules.DeliveryRuleCondition, 0)
	if len(input) == 0 {
		return &results
	}

	conditions := utils.ExpandStringSlice(input)

	for _, condition := range *conditions {
		results = append(results, rules.DeliveryRuleCondition{
			Name: rules.MatchVariable(condition),
		})
	}

	return &results
}

func expandOptionalRuleDeliveryRuleActionArray(input []interface{}) []rules.DeliveryRuleAction {
	if len(input) == 0 {
		return make([]rules.DeliveryRuleAction, 0)
	}

	return expandRuleDeliveryRuleActions(input)
}

func expandRequiredRuleDeliveryRuleActionArray(input []interface{}) *[]rules.DeliveryRuleAction {
	if len(input) == 0 {
		return nil
	}

	results := expandRuleDeliveryRuleActions(input)

	return &results
}

func expandRuleDeliveryRuleActions(input []interface{}) []rules.DeliveryRuleAction {
	results := make([]rules.DeliveryRuleAction, 0)
	actions := utils.ExpandStringSlice(input)

	for _, action := range *actions {
		results = append(results, rules.DeliveryRuleAction(action))
	}

	return results
}

func flattenRuleDeliveryRuleConditionArray(input *[]rules.DeliveryRuleCondition) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, condition := range *input {
		results = append(results, condition.Name)
	}

	return results
}

func flattenRuleDeliveryRuleActionArray(input *[]rules.DeliveryRuleAction) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, action := range *input {
		results = append(results, string(action))
	}

	return results
}
