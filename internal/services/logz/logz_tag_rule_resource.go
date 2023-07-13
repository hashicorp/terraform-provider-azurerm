// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logz

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogzTagRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogzTagRuleCreateUpdate,
		Read:   resourceLogzTagRuleRead,
		Update: resourceLogzTagRuleCreateUpdate,
		Delete: resourceLogzTagRuleDelete,

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
			"logz_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: monitors.ValidateMonitorID,
			},

			"tag_filter": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(tagrules.TagActionInclude),
								string(tagrules.TagActionExclude),
							}, false),
						},

						"value": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"send_aad_logs": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"send_activity_logs": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"send_subscription_logs": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceLogzTagRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := monitors.ParseMonitorID(d.Get("logz_monitor_id").(string))
	if err != nil {
		return err
	}

	id := tagrules.NewTagRuleID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, "default")
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_logz_tag_rule", id.ID())
		}
	}

	payload := tagrules.MonitoringTagRules{
		Properties: &tagrules.MonitoringTagRulesProperties{
			LogRules: expandTagRuleLogRules(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogzTagRuleRead(d, meta)
}

func resourceLogzTagRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("logz_monitor_id", monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if logRules := props.LogRules; logRules != nil {
				d.Set("send_aad_logs", logRules.SendAadLogs)
				d.Set("send_activity_logs", logRules.SendActivityLogs)
				d.Set("send_subscription_logs", logRules.SendSubscriptionLogs)
				if err := d.Set("tag_filter", flattenTagRuleFilteringTagArray(logRules.FilteringTags)); err != nil {
					return fmt.Errorf("setting `tag_filter`: %+v", err)
				}
			}
		}
	}

	return nil
}

func resourceLogzTagRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandTagRuleLogRules(d *pluginsdk.ResourceData) *tagrules.LogRules {
	return &tagrules.LogRules{
		SendAadLogs:          pointer.To(d.Get("send_aad_logs").(bool)),
		SendSubscriptionLogs: pointer.To(d.Get("send_subscription_logs").(bool)),
		SendActivityLogs:     pointer.To(d.Get("send_activity_logs").(bool)),
		FilteringTags:        expandTagRuleFilteringTagArray(d.Get("tag_filter").([]interface{})),
	}
}

func expandTagRuleFilteringTagArray(input []interface{}) *[]tagrules.FilteringTag {
	results := make([]tagrules.FilteringTag, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, tagrules.FilteringTag{
			Name:   utils.String(v["name"].(string)),
			Value:  utils.String(v["value"].(string)),
			Action: pointer.To(tagrules.TagAction(v["action"].(string))),
		})
	}

	return &results
}

func flattenTagRuleFilteringTagArray(input *[]tagrules.FilteringTag) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}
		action := ""
		if item.Action != nil {
			action = string(*item.Action)
		}
		var value string
		if item.Value != nil {
			value = *item.Value
		}
		results = append(results, map[string]interface{}{
			"name":   name,
			"action": action,
			"value":  value,
		})
	}

	return results
}
