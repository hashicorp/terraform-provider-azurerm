package datashare

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDataShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataShareRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ShareName(),
			},

			"account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountID,
			},

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"snapshot_schedule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"recurrence": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"start_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"terms": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDataShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountId, err := parse.AccountID(d.Get("account_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, accountId.ResourceGroup, accountId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("DataShare %q (Account %q / Resource Group %q) was not found", name, accountId.Name, accountId.ResourceGroup)
		}
		return fmt.Errorf("retrieving DataShare %q (Account %q / Resource Group %q): %+v", name, accountId.Name, accountId.ResourceGroup, err)
	}

	dataShareId := parse.NewShareID(subscriptionId, accountId.ResourceGroup, accountId.Name, name).ID()
	d.SetId(dataShareId)

	d.Set("name", name)
	d.Set("account_id", accountId.ID())

	if props := resp.ShareProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("kind", string(props.ShareKind))
		d.Set("terms", props.Terms)
	}

	settings := make([]datashare.ScheduledSynchronizationSetting, 0)
	syncIterator, err := syncClient.ListByShareComplete(ctx, accountId.ResourceGroup, accountId.Name, name, "")
	if err != nil {
		return fmt.Errorf("listing Snapshot Schedules for Data Share %q (Account %q / Resource Group %q): %+v", name, accountId.Name, accountId.ResourceGroup, err)
	}
	for syncIterator.NotDone() {
		item, ok := syncIterator.Value().AsScheduledSynchronizationSetting()
		if ok && item != nil {
			settings = append(settings, *item)
		}

		if err := syncIterator.NextWithContext(ctx); err != nil {
			return fmt.Errorf("retrieving next Snapshot Schedule: %+v", err)
		}
	}

	if err := d.Set("snapshot_schedule", flattenDataShareDataSourceSnapshotSchedule(settings)); err != nil {
		return fmt.Errorf("setting `snapshot_schedule`: %+v", err)
	}

	return nil
}

func flattenDataShareDataSourceSnapshotSchedule(input []datashare.ScheduledSynchronizationSetting) []interface{} {
	output := make([]interface{}, 0)

	for _, sync := range input {
		name := ""
		if sync.Name != nil {
			name = *sync.Name
		}

		startTime := ""
		if sync.SynchronizationTime != nil && !sync.SynchronizationTime.IsZero() {
			startTime = sync.SynchronizationTime.Format(time.RFC3339)
		}

		output = append(output, map[string]interface{}{
			"name":       name,
			"recurrence": string(sync.RecurrenceInterval),
			"start_time": startTime,
		})
	}

	return output
}
