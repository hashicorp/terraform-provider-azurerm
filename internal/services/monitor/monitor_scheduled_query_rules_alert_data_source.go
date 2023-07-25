// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceMonitorScheduledQueryRulesAlert() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMonitorScheduledQueryRulesAlertRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"authorized_resource_ids": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"action": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"action_group": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"custom_webhook_payload": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"email_subject": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
			"data_source_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
			"frequency": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"query": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"query_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"severity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"throttling": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"time_window": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"trigger": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"metric_trigger": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"metric_column": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"metric_trigger_type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"threshold": {
										Type:     pluginsdk.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"operator": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"threshold": {
							Type:     pluginsdk.TypeFloat,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceMonitorScheduledQueryRulesAlertRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scheduledqueryrules.NewScheduledQueryRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("[DEBUG] %s: %+v", id, err)
		}
		return fmt.Errorf("getting Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.ScheduledQueryRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		props := model.Properties
		d.Set("description", props.Description)
		if props.Enabled != nil && *props.Enabled == scheduledqueryrules.EnabledTrue {
			d.Set("enabled", true)
		} else {
			d.Set("enabled", false)
		}

		action, ok := props.Action.(scheduledqueryrules.AlertingAction)
		if !ok {
			return fmt.Errorf("wrong action type in %s: %T", id, props.Action)
		}
		if action.AznsAction != nil {
			if err = d.Set("action", flattenAzureRmScheduledQueryRulesAlertAction(action.AznsAction)); err != nil {
				return fmt.Errorf("setting `action`: %+v", err)
			}
		}
		severity, err := strconv.Atoi(string(action.Severity))
		if err != nil {
			return fmt.Errorf("converting action.Severity %q to int in %s: %+v", action.Severity, id, err)
		}
		d.Set("severity", severity)
		d.Set("throttling", action.ThrottlingInMin)
		if err = d.Set("trigger", flattenAzureRmScheduledQueryRulesAlertTrigger(action.Trigger)); err != nil {
			return fmt.Errorf("setting `trigger`: %+v", err)
		}

		if schedule := props.Schedule; schedule != nil {
			d.Set("frequency", schedule.FrequencyInMinutes)
			d.Set("time_window", schedule.TimeWindowInMinutes)
		}

		d.Set("authorized_resource_ids", utils.FlattenStringSlice(props.Source.AuthorizedResources))
		d.Set("data_source_id", props.Source.DataSourceId)
		d.Set("query", props.Source.Query)
		d.Set("query_type", string(pointer.From(props.Source.QueryType)))

		if err = d.Set("tags", utils.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}

	}

	return nil
}
