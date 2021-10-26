package monitor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-07-01-preview/insights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

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

	id := parse.NewScheduledQueryRulesID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ScheduledQueryRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] Scheduled Query Rule %q was not found in Resource Group %q: %+v", id.ScheduledQueryRuleName, id.ResourceGroup, err)
		}
		return fmt.Errorf("getting Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.ScheduledQueryRuleName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("description", resp.Description)
	if resp.Enabled == insights.EnabledTrue {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

	action, ok := resp.Action.(insights.AlertingAction)
	if !ok {
		return fmt.Errorf("wrong action type in %s: %T", id, resp.Action)
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

	if schedule := resp.Schedule; schedule != nil {
		if schedule.FrequencyInMinutes != nil {
			d.Set("frequency", schedule.FrequencyInMinutes)
		}
		if schedule.TimeWindowInMinutes != nil {
			d.Set("time_window", schedule.TimeWindowInMinutes)
		}
	}

	if source := resp.Source; source != nil {
		if source.AuthorizedResources != nil {
			d.Set("authorized_resource_ids", utils.FlattenStringSlice(source.AuthorizedResources))
		}
		if source.DataSourceID != nil {
			d.Set("data_source_id", source.DataSourceID)
		}
		if source.Query != nil {
			d.Set("query", source.Query)
		}
		d.Set("query_type", string(source.QueryType))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
