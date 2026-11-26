// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access the properties of an AlertingAction scheduled query rule.",
		Read:        dataSourceMonitorScheduledQueryRulesAlertRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the scheduled query rule.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"authorized_resource_ids": {
				Description: "The list of Resource IDs referred into query.",
				Type:        pluginsdk.TypeSet,
				Computed:    true,
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
				Description: "The resource URI over which log search query is to be run.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the scheduled query rule.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Whether this scheduled query rule is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},
			"frequency": {
				Description: "Frequency at which rule condition should be evaluated.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},
			"query": {
				Description: "Log search query.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"query_type": {
				Description: "The type of query results.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"severity": {
				Description: "Severity of the alert.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},
			"throttling": {
				Description: "Time for which alerts should be throttled or suppressed.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},
			"time_window": {
				Description: "Time window for which data needs to be fetched for query.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
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
							Description: "Evaluation operation for rule.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"threshold": {
							Description: "Result or count threshold based on which rule should be triggered.",
							Type:        pluginsdk.TypeFloat,
							Computed:    true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMonitorScheduledQueryRulesAlertRead(d *pluginsdk.ResourceData, meta any) error {
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

		d.Set("authorized_resource_ids", pluginsdk.FlattenSlice(props.Source.AuthorizedResources))
		d.Set("data_source_id", props.Source.DataSourceId)
		d.Set("query", props.Source.Query)
		d.Set("query_type", pointer.FromEnum(props.Source.QueryType))

		if err = d.Set("tags", pluginsdk.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}
