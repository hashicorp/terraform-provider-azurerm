package cdn

import (
	"fmt"
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
						"sample_size": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  4,
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
