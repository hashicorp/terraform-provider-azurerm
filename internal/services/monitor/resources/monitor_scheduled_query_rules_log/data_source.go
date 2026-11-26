// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_log

import (
	"fmt"
	"time"

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
		Description: "Use this data source to access the properties of a LogToMetricAction scheduled query rule.",
		Read:        dataSourceMonitorScheduledQueryRulesLogRead,

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
				Description: "A list of IDs of Resources referred into query.",
				Type:        pluginsdk.TypeSet,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"criteria": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dimension": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Description: "Name of the dimension.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"operator": {
										Description: "Operator for dimension values.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"values": {
										Description: "List of dimension values.",
										Type:        pluginsdk.TypeList,
										Computed:    true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"metric_name": {
							Description: "Name of the metric.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMonitorScheduledQueryRulesLogRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scheduledqueryrules.NewScheduledQueryRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("[DEBUG] %s was not found : %+v", id, err)
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

		action, ok := props.Action.(scheduledqueryrules.LogToMetricAction)
		if !ok {
			return fmt.Errorf("wrong action type in %s: %T", id, props.Action)
		}
		if err = d.Set("criteria", flattenAzureRmScheduledQueryRulesLogCriteria(action.Criteria)); err != nil {
			return fmt.Errorf("setting `criteria`: %+v", err)
		}

		if schedule := props.Schedule; schedule != nil {
			d.Set("frequency", schedule.FrequencyInMinutes)
			d.Set("time_window", schedule.TimeWindowInMinutes)
		}

		if props.Source.AuthorizedResources != nil {
			d.Set("authorized_resource_ids", pluginsdk.FlattenSlice(props.Source.AuthorizedResources))
			d.Set("data_source_id", props.Source.DataSourceId)
		}

		if err = d.Set("tags", pluginsdk.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}
