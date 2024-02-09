// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorScheduledQueryRulesLog() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorScheduledQueryRulesLogCreateUpdate,
		Read:   resourceMonitorScheduledQueryRulesLogRead,
		Update: resourceMonitorScheduledQueryRulesLogCreateUpdate,
		Delete: resourceMonitorScheduledQueryRulesLogDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := scheduledqueryrules.ParseScheduledQueryRuleID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ScheduledQueryRulesLogUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringDoesNotContainAny("<>*%&:\\?+/"),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"authorized_resource_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"criteria": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dimension": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  "Include",
										ValidateFunc: validation.StringInSlice([]string{
											"Include",
										}, false),
									},
									"values": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
						"metric_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"data_source_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
			},
			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorScheduledQueryRulesLogCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	action := expandMonitorScheduledQueryRulesLogToMetricAction(d)
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scheduledqueryrules.NewScheduledQueryRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules_alert", id.ID())
		}
	}

	description := d.Get("description").(string)
	enabledRaw := d.Get("enabled").(bool)

	enabled := scheduledqueryrules.EnabledTrue
	if !enabledRaw {
		enabled = scheduledqueryrules.EnabledFalse
	}

	location := azure.NormalizeLocation(d.Get("location"))

	source := expandMonitorScheduledQueryRulesCommonSource(d)

	t := d.Get("tags").(map[string]interface{})

	parameters := scheduledqueryrules.LogSearchRuleResource{
		Location: location,
		Properties: scheduledqueryrules.LogSearchRule{
			Description: utils.String(description),
			Enabled:     pointer.To(enabled),
			Source:      source,
			Action:      action,
		},
		Tags: utils.ExpandPtrMapStringString(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorScheduledQueryRulesLogRead(d, meta)
}

func resourceMonitorScheduledQueryRulesLogRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scheduledqueryrules.ParseScheduledQueryRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Monitor %s: %+v", *id, err)
	}

	d.Set("name", id.ScheduledQueryRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		props := model.Properties

		d.Set("description", props.Description)
		if props.Enabled != nil && *props.Enabled == scheduledqueryrules.EnabledTrue {
			d.Set("enabled", true)
		} else {
			d.Set("enabled", false)
		}

		action, ok := props.Action.(scheduledqueryrules.LogToMetricAction)
		if !ok {
			return fmt.Errorf("wrong action type in %s: %T", *id, props.Action)
		}

		if err = d.Set("criteria", flattenAzureRmScheduledQueryRulesLogCriteria(action.Criteria)); err != nil {
			return fmt.Errorf("setting `criteria`: %+v", err)
		}

		if props.Source.AuthorizedResources != nil {
			d.Set("authorized_resource_ids", utils.FlattenStringSlice(props.Source.AuthorizedResources))
		}

		d.Set("data_source_id", props.Source.DataSourceId)

		if err = d.Set("tags", utils.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceMonitorScheduledQueryRulesLogDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scheduledqueryrules.ParseScheduledQueryRuleID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
		}
	}

	return nil
}

func expandMonitorScheduledQueryRulesLogCriteria(input []interface{}) []scheduledqueryrules.Criteria {
	criteria := make([]scheduledqueryrules.Criteria, 0)
	if len(input) == 0 {
		return criteria
	}

	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		dimensions := make([]scheduledqueryrules.Dimension, 0)
		for _, dimension := range v["dimension"].(*pluginsdk.Set).List() {
			if dimension == nil {
				continue
			}
			dVal, ok := dimension.(map[string]interface{})
			if !ok {
				continue
			}

			dimensions = append(dimensions, scheduledqueryrules.Dimension{
				Name:     dVal["name"].(string),
				Operator: scheduledqueryrules.Operator(dVal["operator"].(string)),
				Values:   expandStringValues(dVal["values"].([]interface{})),
			})
		}

		criteria = append(criteria, scheduledqueryrules.Criteria{
			MetricName: v["metric_name"].(string),
			Dimensions: &dimensions,
		})
	}
	return criteria
}

func expandMonitorScheduledQueryRulesLogToMetricAction(d *pluginsdk.ResourceData) *scheduledqueryrules.LogToMetricAction {
	criteriaRaw := d.Get("criteria").([]interface{})
	criteria := expandMonitorScheduledQueryRulesLogCriteria(criteriaRaw)

	action := scheduledqueryrules.LogToMetricAction{
		Criteria: criteria,
	}

	return &action
}
