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
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorProfileRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorProfileRuleCreate,
		Read:   resourceFrontdoorProfileRuleRead,
		Update: resourceFrontdoorProfileRuleUpdate,
		Delete: resourceFrontdoorProfileRuleDelete,

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

			"cdn_rule_set_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: rulesets.ValidateRuleSetID,
			},

			"actions": {
				Type:     pluginsdk.TypeList,
				Required: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"conditions": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"deployment_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"match_processing_behavior": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"order": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"rule_set_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFrontdoorProfileRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	ruleSetId, err := rulesets.ParseRuleSetID(d.Get("cdn_rule_set_id").(string))
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
			return tf.ImportAsExistsError("azurerm_cdn_rule", id.ID())
		}
	}

	matchProcessingBehaviorValue := rules.MatchProcessingBehavior(d.Get("match_processing_behavior").(string))
	props := rules.Rule{
		Properties: &rules.RuleProperties{
			Actions:                 expandRuleDeliveryRuleActionArray(d.Get("actions").([]interface{})),
			Conditions:              expandRuleDeliveryRuleConditionArray(d.Get("conditions").([]interface{})),
			MatchProcessingBehavior: &matchProcessingBehaviorValue,
			Order:                   int64(d.Get("order").(int)),
		},
	}
	if err := client.CreateThenPoll(ctx, id, props); err != nil {

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorProfileRuleRead(d, meta)
}

func resourceFrontdoorProfileRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileRulesClient
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

	d.Set("cdn_rule_set_id", rulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {

			if err := d.Set("actions", flattenRuleDeliveryRuleActionArray(&props.Actions)); err != nil {
				return fmt.Errorf("setting `actions`: %+v", err)
			}

			if err := d.Set("conditions", flattenRuleDeliveryRuleConditionArray(props.Conditions)); err != nil {
				return fmt.Errorf("setting `conditions`: %+v", err)
			}
			d.Set("deployment_status", props.DeploymentStatus)
			d.Set("match_processing_behavior", props.MatchProcessingBehavior)
			d.Set("order", props.Order)
			d.Set("provisioning_state", props.ProvisioningState)
			d.Set("rule_set_name", props.RuleSetName)
		}
	}
	return nil
}

func resourceFrontdoorProfileRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileRulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseRuleID(d.Id())
	if err != nil {
		return err
	}

	matchProcessingBehaviorValue := rules.MatchProcessingBehavior(d.Get("match_processing_behavior").(string))
	props := rules.RuleUpdateParameters{
		Properties: &rules.RuleUpdatePropertiesParameters{
			Actions:                 expandRuleDeliveryRuleActionArray(d.Get("actions").([]interface{})),
			Conditions:              expandRuleDeliveryRuleConditionArray(d.Get("conditions").([]interface{})),
			MatchProcessingBehavior: &matchProcessingBehaviorValue,
			Order:                   utils.Int64(int64(d.Get("order").(int))),
		},
	}
	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorProfileRuleRead(d, meta)
}

func resourceFrontdoorProfileRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileRulesClient
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
	for _, item := range input {
		v := item.(map[string]interface{})
		nameValue := rules.MatchVariable(v["name"].(string))
		results = append(results, rules.DeliveryRuleCondition{
			Name: nameValue,
		})
	}
	return &results
}

func expandRuleDeliveryRuleActionArray(input []interface{}) *[]rules.DeliveryRuleAction {
	results := make([]rules.DeliveryRuleAction, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		nameValue := rules.DeliveryRuleAction(v["name"].(string))
		results = append(results, rules.DeliveryRuleAction{
			Name: nameValue,
		})
	}
	return &results
}

func flattenRuleDeliveryRuleConditionArray(inputs *[]rules.DeliveryRuleCondition) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})
		result["name"] = input.Name
		results = append(results, result)
	}

	return results
}

func flattenRuleDeliveryRuleActionArray(inputs *[]rules.DeliveryRuleAction) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})
		result["name"] = input.Name
		results = append(results, result)
	}

	return results
}
