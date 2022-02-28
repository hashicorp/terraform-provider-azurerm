package elastic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/sdk/2020-07-01/monitorsresource"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/sdk/2020-07-01/rules"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceElasticTagRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceElasticTagRuleCreate,
		Read:   resourceElasticTagRuleRead,
		Update: resourceElasticTagRuleUpdate,
		Delete: resourceElasticTagRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := rules.ParseTagRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: monitorsresource.ValidateMonitorID,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  utils.String("default"),
			},

			"log_rules": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"send_aad_logs": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"send_subscription_logs": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"send_activity_logs": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"filtering_tag": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"action": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Include",
											"Exclude",
										}, false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceElasticTagRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := monitorsresource.ParseMonitorID(d.Get("monitor_id").(string))
	if err != nil {
		return err
	}

	id := rules.NewTagRuleID(subscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, d.Get("name").(string))
	existing, err := client.TagRulesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_elastic_tag_rule", id.ID())
	}

	body := rules.MonitoringTagRules{
		Properties: &rules.MonitoringTagRulesProperties{
			LogRules: expandLogRules(d.Get("log_rules").([]interface{})),
		},
	}
	if _, err := client.TagRulesCreateOrUpdate(ctx, id, body); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceElasticTagRuleRead(d, meta)
}

func resourceElasticTagRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	body := rules.MonitoringTagRules{
		Properties: &rules.MonitoringTagRulesProperties{
			LogRules: expandLogRules(d.Get("log_rules").([]interface{})),
		},
	}
	if _, err := client.TagRulesCreateOrUpdate(ctx, *id, body); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceElasticTagRuleRead(d, meta)
}

func resourceElasticTagRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.TagRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.RuleSetName)
	d.Set("monitor_id", monitorsresource.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("log_rules", flattenLogRules(props.LogRules)); err != nil {
				return fmt.Errorf("setting `log_rules`: %+v", err)
			}
		}
	}

	return nil
}

func resourceElasticTagRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	if err := client.TagRulesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandLogRules(input []interface{}) *rules.LogRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filtering_tag"].([]interface{})

	return &rules.LogRules{
		SendAadLogs:          utils.Bool(v["send_aad_logs"].(bool)),
		SendSubscriptionLogs: utils.Bool(v["send_subscription_logs"].(bool)),
		SendActivityLogs:     utils.Bool(v["send_activity_logs"].(bool)),
		FilteringTags:        expandFilteringTag(filteringTag),
	}
}

func expandFilteringTag(input []interface{}) *[]rules.FilteringTag {
	filteringTags := make([]rules.FilteringTag, 0)

	for _, v := range input {
		config := v.(map[string]interface{})
		name := config["name"].(string)
		value := config["value"].(string)
		action := rules.TagAction(config["action"].(string))

		filteringTags = append(filteringTags, rules.FilteringTag{
			Name:   utils.String(name),
			Value:  utils.String(value),
			Action: &action,
		})
	}

	return &filteringTags
}

func flattenLogRules(input *rules.LogRules) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return make([]interface{}, 0)
	}

	sendAadLogs := false
	if input.SendAadLogs != nil {
		sendAadLogs = *input.SendAadLogs
	}

	sendSubscriptionLogs := false
	if input.SendSubscriptionLogs != nil {
		sendSubscriptionLogs = *input.SendSubscriptionLogs
	}

	sendActivityLogs := false
	if input.SendActivityLogs != nil {
		sendActivityLogs = *input.SendActivityLogs
	}

	return append(results, map[string]interface{}{
		"filtering_tag":          flattenFilteringTags(input.FilteringTags),
		"send_aad_logs":          sendAadLogs,
		"send_activity_logs":     sendActivityLogs,
		"send_subscription_logs": sendSubscriptionLogs,
	})
}

func flattenFilteringTags(input *[]rules.FilteringTag) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, filteringTagRules := range *input {
		action := ""
		if filteringTagRules.Action != nil {
			action = string(*filteringTagRules.Action)
		}

		name := ""
		if filteringTagRules.Name != nil {
			name = *filteringTagRules.Name
		}

		value := ""
		if filteringTagRules.Value != nil {
			value = *filteringTagRules.Value
		}
		results = append(results, map[string]interface{}{
			"action": action,
			"name":   name,
			"value":  value,
		})
	}
	return results
}
