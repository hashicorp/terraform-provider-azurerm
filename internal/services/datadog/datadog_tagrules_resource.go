package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDatadogTagRules() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatadogTagRulesCreateorUpdate,
		Read:   resourceDatadogTagRulesRead,
		Update: resourceDatadogTagRulesCreateorUpdate,
		Delete: resourceDatadogTagRulesDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatadogMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DatadogMonitorsName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"rule_set_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  utils.String("default"),
			},

			"id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
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

						"send_resource_logs": {
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

			"metric_rules": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
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

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDatadogTagRulesCreateorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Datadog.TagRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	ruleSetName := d.Get("rule_set_name").(string)

	id := parse.NewDatadogTagRulesID(subscriptionId, resourceGroup, name, ruleSetName).ID()

	existing, err := client.Get(ctx, resourceGroup, name, ruleSetName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing Datadog Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	body := datadog.MonitoringTagRules{
		Properties: &datadog.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log_rules").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric_rules").([]interface{})),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, ruleSetName, &body); err != nil {
		return fmt.Errorf("configuring Tag Rules on Datadog Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id)
	return resourceDatadogTagRulesRead(d, meta)
}

func resourceDatadogTagRulesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.TagRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatadogTagRulesID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Datadog monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	if props := resp.Properties; props != nil {
		if err := d.Set("log_rules", flattenLogRules(props.LogRules)); err != nil {
			return fmt.Errorf("setting `log_rules`: %+v", err)
		}
		if err := d.Set("metric_rules", flattenMetricRules(props.MetricRules)); err != nil {
			return fmt.Errorf("setting `metric_rules`: %+v", err)
		}
		d.Set("provisioning_state", props.ProvisioningState)
	}
	d.Set("type", resp.Type)
	d.Set("name", resp.Name)
	d.Set("id", resp.ID)

	return nil
}

func resourceDatadogTagRulesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.TagRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatadogTagRulesID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] datadog %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	d.Set("log_rules", nil)
	d.Set("metric_rules", nil)

	body := datadog.MonitoringTagRules{
		Properties: &datadog.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log_rules").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric_rules").([]interface{})),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName, &body); err != nil {
		return fmt.Errorf("removing Tag Rules configuration from Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	return nil
}

func expandLogRules(input []interface{}) *datadog.LogRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filtering_tag"].([]interface{})

	return &datadog.LogRules{
		SendAadLogs:          utils.Bool(v["send_aad_logs"].(bool)),
		SendSubscriptionLogs: utils.Bool(v["send_subscription_logs"].(bool)),
		SendResourceLogs:     utils.Bool(v["send_resource_logs"].(bool)),
		FilteringTags:        expandFilteringTag(filteringTag),
	}
}

func expandMetricRules(input []interface{}) *datadog.MetricRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filtering_tag"].([]interface{})

	return &datadog.MetricRules{
		FilteringTags: expandFilteringTag(filteringTag),
	}
}

func expandFilteringTag(input []interface{}) *[]datadog.FilteringTag {

	filteringTags := make([]datadog.FilteringTag, 0)

	for _, v := range input {
		config := v.(map[string]interface{})
		name := config["name"].(string)
		value := config["value"].(string)
		action := config["action"].(string)

		filteringTag := datadog.FilteringTag{
			Name:   utils.String(name),
			Value:  utils.String(value),
			Action: datadog.TagAction(action),
		}

		filteringTags = append(filteringTags, filteringTag)
	}

	return &filteringTags
}

func flattenLogRules(input *datadog.LogRules) []interface{} {
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

	if input.SendResourceLogs != nil {
		result["send_resource_logs"] = *input.SendResourceLogs
	}

	result["filtering_tag"] = flattenFilteringTags(input.FilteringTags)
	return append(results, result)

}

func flattenMetricRules(input *datadog.MetricRules) []interface{} {
	results := make([]interface{}, 0)

	if input == nil {
		return make([]interface{}, 0)
	}
	result := make(map[string]interface{})

	result["filtering_tag"] = flattenFilteringTags(input.FilteringTags)
	return append(results, result)
}

func flattenFilteringTags(input *[]datadog.FilteringTag) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	var t datadog.TagAction

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
