// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogTagRules() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatadogTagRulesCreate,
		Read:   resourceDatadogTagRulesRead,
		Update: resourceDatadogTagRulesUpdate,
		Delete: resourceDatadogTagRulesDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := tagrules.ParseTagRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"datadog_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tagrules.ValidateMonitorID,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "default",
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
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(tagrules.PossibleValuesForTagAction(), false),
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
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(tagrules.PossibleValuesForTagAction(), false),
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

func resourceDatadogTagRulesCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.TagRules
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := tagrules.ParseMonitorID(d.Get("datadog_monitor_id").(string))
	if err != nil {
		return err
	}

	id := tagrules.NewTagRuleID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for an existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) && !isDefaultSettings(existing.Model) {
			return tf.ImportAsExistsError("azurerm_datadog_monitor_tag_rule", id.ID())
		}
	}

	payload := tagrules.MonitoringTagRules{
		Properties: &tagrules.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric").([]interface{})),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatadogTagRulesRead(d, meta)
}

func resourceDatadogTagRulesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.TagRules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
	}

	monitorId := tagrules.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName)
	d.Set("datadog_monitor_id", monitorId.ID())
	d.Set("name", id.TagRuleName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("log", flattenLogRules(props.LogRules)); err != nil {
				return fmt.Errorf("setting `log`: %+v", err)
			}
			if err := d.Set("metric", flattenMetricRules(props.MetricRules)); err != nil {
				return fmt.Errorf("setting `metric`: %+v", err)
			}
		}
	}

	return nil
}

func resourceDatadogTagRulesUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.TagRules
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	payload := tagrules.MonitoringTagRules{
		Properties: &tagrules.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric").([]interface{})),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDatadogTagRulesRead(d, meta)
}

func resourceDatadogTagRulesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.TagRules
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	// Tag Rules can't be removed on their own, they can only be nil'd out
	payload := tagrules.MonitoringTagRules{
		Properties: &tagrules.MonitoringTagRulesProperties{
			LogRules: &tagrules.LogRules{
				SendAadLogs:          pointer.To(false),
				SendSubscriptionLogs: pointer.To(false),
				SendResourceLogs:     pointer.To(false),
				FilteringTags:        &[]tagrules.FilteringTag{},
			},
			MetricRules: &tagrules.MetricRules{
				FilteringTags: &[]tagrules.FilteringTag{},
			},
		},
	}
	if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("removing %s: %+v", *id, err)
	}

	return nil
}

func expandLogRules(input []interface{}) *tagrules.LogRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filter"].([]interface{})

	return &tagrules.LogRules{
		SendAadLogs:          pointer.To(v["aad_log_enabled"].(bool)),
		SendSubscriptionLogs: pointer.To(v["subscription_log_enabled"].(bool)),
		SendResourceLogs:     pointer.To(v["resource_log_enabled"].(bool)),
		FilteringTags:        expandFilteringTag(filteringTag),
	}
}

func expandMetricRules(input []interface{}) *tagrules.MetricRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filter"].([]interface{})

	return &tagrules.MetricRules{
		FilteringTags: expandFilteringTag(filteringTag),
	}
}

func expandFilteringTag(input []interface{}) *[]tagrules.FilteringTag {
	filteringTags := make([]tagrules.FilteringTag, 0)

	for _, v := range input {
		config := v.(map[string]interface{})

		filteringTags = append(filteringTags, tagrules.FilteringTag{
			Name:   pointer.To(config["name"].(string)),
			Value:  pointer.To(config["value"].(string)),
			Action: pointer.ToEnum[tagrules.TagAction](config["action"].(string)),
		})
	}

	return &filteringTags
}

func flattenLogRules(input *tagrules.LogRules) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		results = append(results, map[string]interface{}{
			"aad_log_enabled":          pointer.From(input.SendAadLogs),
			"filter":                   flattenFilteringTags(input.FilteringTags),
			"resource_log_enabled":     pointer.From(input.SendResourceLogs),
			"subscription_log_enabled": pointer.From(input.SendSubscriptionLogs),
		})
	}

	return results
}

func flattenMetricRules(input *tagrules.MetricRules) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"filter": flattenFilteringTags(input.FilteringTags),
		},
	}
}

func flattenFilteringTags(input *[]tagrules.FilteringTag) []interface{} {
	results := make([]interface{}, 0)
	if input != nil {
		for _, filteringTagRules := range *input {
			action := ""
			if filteringTagRules.Action != nil {
				action = string(*filteringTagRules.Action)
			}
			results = append(results, map[string]interface{}{
				"action": action,
				"name":   pointer.From(filteringTagRules.Name),
				"value":  pointer.From(filteringTagRules.Value),
			})
		}
	}
	return results
}

func isDefaultSettings(input *tagrules.MonitoringTagRules) bool {
	if input == nil {
		return false
	}

	if input.Properties == nil || input.Properties.LogRules == nil || input.Properties.MetricRules == nil {
		return false
	}

	logRules := input.Properties.LogRules
	metricRules := input.Properties.MetricRules
	return (logRules.SendAadLogs != nil && !*logRules.SendAadLogs) &&
		(logRules.SendSubscriptionLogs != nil && !*logRules.SendSubscriptionLogs) &&
		(logRules.SendResourceLogs != nil && !*logRules.SendResourceLogs) &&
		(logRules.FilteringTags != nil && len(*logRules.FilteringTags) == 0) &&
		(metricRules.FilteringTags != nil && len(*metricRules.FilteringTags) == 0)
}
