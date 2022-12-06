package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog"
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
			_, err := parse.DatadogTagRulesID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"datadog_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatadogMonitorID,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  utils.String("default"),
			},

			"log": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"aad_log_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"subscription_log_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"resource_log_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"filter": {
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

			"metric": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"filter": {
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

func resourceDatadogTagRulesCreateorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.TagRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	datadogMonitorId := d.Get("datadog_monitor_id").(string)
	ruleSetName := d.Get("name").(string)
	id, err := parse.DatadogMonitorID(datadogMonitorId)
	if err != nil {
		return err
	}

	tagRulesid := parse.NewDatadogTagRulesID(id.SubscriptionId, id.ResourceGroup, id.MonitorName, ruleSetName).ID()

	existing, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, ruleSetName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing Datadog Monitor %q (Resource Group %q): %+v", id.ResourceGroup, id.MonitorName, err)
		}
	}

	body := datadog.MonitoringTagRules{
		Properties: &datadog.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric").([]interface{})),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MonitorName, ruleSetName, &body); err != nil {
		return fmt.Errorf("configuring Tag Rules on Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	d.SetId(tagRulesid)
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

	monitorId := parse.NewDatadogMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName)
	d.Set("datadog_monitor_id", monitorId.ID())
	d.Set("name", id.TagRuleName)

	if props := resp.Properties; props != nil {
		if err := d.Set("log", flattenLogRules(props.LogRules)); err != nil {
			return fmt.Errorf("setting `log`: %+v", err)
		}
		if err := d.Set("metric", flattenMetricRules(props.MetricRules)); err != nil {
			return fmt.Errorf("setting `metric`: %+v", err)
		}
	}

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

	d.Set("log", nil)
	d.Set("metric", nil)

	body := datadog.MonitoringTagRules{
		Properties: &datadog.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric").([]interface{})),
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
	filteringTag := v["filter"].([]interface{})

	return &datadog.LogRules{
		SendAadLogs:          utils.Bool(v["aad_log_enabled"].(bool)),
		SendSubscriptionLogs: utils.Bool(v["subscription_log_enabled"].(bool)),
		SendResourceLogs:     utils.Bool(v["resource_log_enabled"].(bool)),
		FilteringTags:        expandFilteringTag(filteringTag),
	}
}

func expandMetricRules(input []interface{}) *datadog.MetricRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filter"].([]interface{})

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
		result["aad_log_enabled"] = *input.SendAadLogs
	}

	if input.SendSubscriptionLogs != nil {
		result["subscription_log_enabled"] = *input.SendSubscriptionLogs
	}

	if input.SendResourceLogs != nil {
		result["resource_log_enabled"] = *input.SendResourceLogs
	}

	result["filter"] = flattenFilteringTags(input.FilteringTags)
	return append(results, result)

}

func flattenMetricRules(input *datadog.MetricRules) []interface{} {
	results := make([]interface{}, 0)

	if input == nil {
		return make([]interface{}, 0)
	}
	result := make(map[string]interface{})

	result["filter"] = flattenFilteringTags(input.FilteringTags)
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
