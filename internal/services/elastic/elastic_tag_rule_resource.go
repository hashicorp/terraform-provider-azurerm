package elastic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/legacysdk/elastic/mgmt/2020-07-01/elastic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceElasticTagRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceElasticTagRuleCreateorUpdate,
		Read:   resourceElasticTagRuleRead,
		Update: resourceElasticTagRuleCreateorUpdate,
		Delete: resourceElasticTagRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ElasticTagRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"monitor_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ElasticMonitorName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"rule_set_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  utils.String("default"),
			},

			"log_rules": {
				Type:     pluginsdk.TypeList,
				Optional: true,
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
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"value": {
										Type:     pluginsdk.TypeString,
										Required: true,
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

func resourceElasticTagRuleCreateorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("monitor_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	ruleSetName := d.Get("rule_set_name").(string)

	id := parse.NewElasticTagRuleID(subscriptionId, resourceGroup, name, ruleSetName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, ruleSetName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_elastic_tag_rule", *existing.ID)
		}
	}

	body := elastic.MonitoringTagRules{
		Properties: &elastic.MonitoringTagRulesProperties{
			LogRules: expandLogRules(d.Get("log_rules").([]interface{})),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, ruleSetName, &body); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceElasticTagRuleRead(d, meta)
}

func resourceElasticTagRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ElasticTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Elastic monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	d.Set("monitor_name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("rule_set_name", id.TagRuleName)

	if props := resp.Properties; props != nil {
		if err := d.Set("log_rules", flattenLogRules(props.LogRules)); err != nil {
			return fmt.Errorf("setting `log_rules`: %+v", err)
		}
	}
	d.SetId(id.ID())

	return nil
}

func resourceElasticTagRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err1 := parse.ElasticTagRuleID(d.Id())
	if err1 != nil {
		return err1
	}

	resp, err2 := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
	if err2 != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Elastic %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	d.Set("log_rules", nil)

	body := elastic.MonitoringTagRules{
		Properties: &elastic.MonitoringTagRulesProperties{
			LogRules: expandLogRules(d.Get("log_rules").([]interface{})),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName, &body); err != nil {
		return fmt.Errorf("removing Tag Rules configuration from Elastic Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	return nil
}

func expandLogRules(input []interface{}) *elastic.LogRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filtering_tag"].([]interface{})

	return &elastic.LogRules{
		SendAadLogs:          utils.Bool(v["send_aad_logs"].(bool)),
		SendSubscriptionLogs: utils.Bool(v["send_subscription_logs"].(bool)),
		SendActivityLogs:     utils.Bool(v["send_activity_logs"].(bool)),
		FilteringTags:        expandFilteringTag(filteringTag),
	}
}

func expandFilteringTag(input []interface{}) *[]elastic.FilteringTag {

	filteringTags := make([]elastic.FilteringTag, 0)

	for _, v := range input {
		config := v.(map[string]interface{})
		name := config["name"].(string)
		value := config["value"].(string)
		action := config["action"].(string)

		filteringTag := elastic.FilteringTag{
			Name:   utils.String(name),
			Value:  utils.String(value),
			Action: elastic.TagAction(action),
		}

		filteringTags = append(filteringTags, filteringTag)
	}

	return &filteringTags
}

func flattenLogRules(input *elastic.LogRules) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if input.SendAadLogs != nil {
		result["send_aad_logs"] = *input.SendAadLogs
	}

	if input.SendSubscriptionLogs != nil {
		result["send_subscription_logs"] = *input.SendSubscriptionLogs
	}

	if input.SendActivityLogs != nil {
		result["send_activity_logs"] = *input.SendActivityLogs
	}

	result["filtering_tag"] = flattenFilteringTags(input.FilteringTags)
	return append(results, result)

}

func flattenFilteringTags(input *[]elastic.FilteringTag) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	var t elastic.TagAction

	for _, filteringTagRules := range *input {
		result := make(map[string]interface{})

		if filteringTagRules.Name != nil {
			result["name"] = *filteringTagRules.Name
		}
		if filteringTagRules.Value != nil {
			result["value"] = *filteringTagRules.Value
		}
		if filteringTagRules.Action != "" {
			t = filteringTagRules.Action
			result["action"] = t
		}
		results = append(results, result)
	}
	return results
}
