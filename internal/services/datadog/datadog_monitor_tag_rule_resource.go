// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// @tombuildsstuff: in 4.0 consider inlining this within the `azurerm_datadog_monitors` resource
// since this appears to be a 1:1 with it (given the name defaults to `default`)

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
			_, err := rules.ParseTagRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"datadog_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: monitorsresource.ValidateMonitorID,
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
										ValidateFunc: validation.StringInSlice(rules.PossibleValuesForTagAction(), false),
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
										ValidateFunc: validation.StringInSlice(rules.PossibleValuesForTagAction(), false),
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
	client := meta.(*clients.Client).Datadog.Rules
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := monitorsresource.ParseMonitorID(d.Get("datadog_monitor_id").(string))
	if err != nil {
		return err
	}

	id := rules.NewTagRuleID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, d.Get("name").(string))
	existing, err := client.TagRulesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for an existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_datadog_monitor_tag_rule", id.ID())
	}

	payload := rules.MonitoringTagRules{
		Properties: &rules.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric").([]interface{})),
		},
	}
	if _, err := client.TagRulesCreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatadogTagRulesRead(d, meta)
}

func resourceDatadogTagRulesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.Rules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.TagRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
	}

	monitorId := monitorsresource.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName)
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
	client := meta.(*clients.Client).Datadog.Rules
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	payload := rules.MonitoringTagRules{
		Properties: &rules.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric").([]interface{})),
		},
	}
	if _, err := client.TagRulesCreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDatadogTagRulesRead(d, meta)
}

func resourceDatadogTagRulesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.Rules
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	// Tag Rules can't be removed on their own, they can only be nil'd out
	payload := rules.MonitoringTagRules{
		Properties: &rules.MonitoringTagRulesProperties{
			LogRules:    expandLogRules(d.Get("log").([]interface{})),
			MetricRules: expandMetricRules(d.Get("metric").([]interface{})),
		},
	}
	if _, err := client.TagRulesCreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("removing %s: %+v", *id, err)
	}

	return nil
}

func expandLogRules(input []interface{}) *rules.LogRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filter"].([]interface{})

	return &rules.LogRules{
		SendAadLogs:          utils.Bool(v["aad_log_enabled"].(bool)),
		SendSubscriptionLogs: utils.Bool(v["subscription_log_enabled"].(bool)),
		SendResourceLogs:     utils.Bool(v["resource_log_enabled"].(bool)),
		FilteringTags:        expandFilteringTag(filteringTag),
	}
}

func expandMetricRules(input []interface{}) *rules.MetricRules {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	filteringTag := v["filter"].([]interface{})

	return &rules.MetricRules{
		FilteringTags: expandFilteringTag(filteringTag),
	}
}

func expandFilteringTag(input []interface{}) *[]rules.FilteringTag {
	filteringTags := make([]rules.FilteringTag, 0)

	for _, v := range input {
		config := v.(map[string]interface{})

		filteringTags = append(filteringTags, rules.FilteringTag{
			Name:   utils.String(config["name"].(string)),
			Value:  utils.String(config["value"].(string)),
			Action: pointer.To(rules.TagAction(config["action"].(string))),
		})
	}

	return &filteringTags
}

func flattenLogRules(input *rules.LogRules) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		aadLogEnabled := false
		if input.SendAadLogs != nil {
			aadLogEnabled = *input.SendAadLogs
		}

		subscriptionLogEnabled := false
		if input.SendSubscriptionLogs != nil {
			subscriptionLogEnabled = *input.SendSubscriptionLogs
		}

		resourceLogEnabled := false
		if input.SendResourceLogs != nil {
			resourceLogEnabled = *input.SendResourceLogs
		}

		results = append(results, map[string]interface{}{
			"aad_log_enabled":          aadLogEnabled,
			"filter":                   flattenFilteringTags(input.FilteringTags),
			"resource_log_enabled":     resourceLogEnabled,
			"subscription_log_enabled": subscriptionLogEnabled,
		})
	}

	return results
}

func flattenMetricRules(input *rules.MetricRules) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"filter": flattenFilteringTags(input.FilteringTags),
		},
	}
}

func flattenFilteringTags(input *[]rules.FilteringTag) []interface{} {
	results := make([]interface{}, 0)
	if input != nil {
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
	}
	return results
}
