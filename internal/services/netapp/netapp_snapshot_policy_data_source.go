package netapp

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceNetAppSnapshotPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNetAppSnapshotPolicyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SnapshotName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocationForDataSource(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"hourly_schedule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshots_to_keep": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"minute": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"daily_schedule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshots_to_keep": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"hour": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"minute": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"weekly_schedule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshots_to_keep": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"days_of_week": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"hour": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"minute": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"monthly_schedule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshots_to_keep": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"days_of_month": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeInt,
							},
						},

						"hour": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"minute": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceNetAppSnapshotPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSnapshotPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("account_name", id.NetAppAccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.SnapshotPolicyProperties; props != nil {
		d.Set("enabled", props.Enabled)
		if err := d.Set("hourly_schedule", flattenNetAppVolumeSnapshotPolicyHourlySchedule(props.HourlySchedule)); err != nil {
			return fmt.Errorf("setting `hourly_schedule`: %+v", err)
		}
		if err := d.Set("daily_schedule", flattenNetAppVolumeSnapshotPolicyDailySchedule(props.DailySchedule)); err != nil {
			return fmt.Errorf("setting `daily_schedule`: %+v", err)
		}
		if err := d.Set("weekly_schedule", flattenNetAppVolumeSnapshotPolicyWeeklySchedule(props.WeeklySchedule)); err != nil {
			return fmt.Errorf("setting `weekly_schedule`: %+v", err)
		}
		if err := d.Set("monthly_schedule", flattenNetAppVolumeSnapshotPolicyMonthlySchedule(props.MonthlySchedule)); err != nil {
			return fmt.Errorf("setting `monthly_schedule`: %+v", err)
		}
	}

	return nil
}
