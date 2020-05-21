package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDataShare() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDataShareRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DatashareName(),
			},

			"account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DatashareAccountID,
			},

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"snapshot_schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"recurrence": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"terms": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmDataShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountID := d.Get("account_id").(string)
	accountId, err := parse.DataShareAccountID(accountID)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, accountId.ResourceGroup, accountId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] DataShare %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataShare %q (Resource Group %q / accountName %q): %+v", name, accountId.ResourceGroup, accountId.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading DataShare %q (Resource Group %q / accountName %q): ID is empty", name, accountId.ResourceGroup, accountId.Name)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("account_id", accountID)
	if props := resp.ShareProperties; props != nil {
		d.Set("kind", props.ShareKind)
		d.Set("description", props.Description)
		d.Set("terms", props.Terms)
	}

	if syncIterator, err := syncClient.ListByShareComplete(ctx, accountId.ResourceGroup, accountId.Name, name, ""); syncIterator.NotDone() {
		if err != nil {
			return fmt.Errorf("listing DataShare %q snapshot schedule (Resource Group %q / accountName %q): %+v", name, accountId.ResourceGroup, accountId.Name, err)
		}
		if syncName := syncIterator.Value().(datashare.ScheduledSynchronizationSetting).Name; syncName != nil && *syncName != "" {
			syncResp, err := syncClient.Get(ctx, accountId.ResourceGroup, accountId.Name, name, *syncName)
			if err != nil {
				return fmt.Errorf("reading DataShare %q snapshot schedule (Resource Group %q / accountName %q): %+v", name, accountId.ResourceGroup, accountId.Name, err)
			}
			if schedule := syncResp.Value.(datashare.ScheduledSynchronizationSetting); schedule.ID != nil && *schedule.ID != "" {
				if err := d.Set("snapshot_schedule", flattenAzureRmDataShareSnapshotSchedule(&schedule)); err != nil {
					return fmt.Errorf("setting `snapshot_schedule`: %+v", err)
				}
			}
		}
		if err := syncIterator.NextWithContext(ctx); err != nil {
			return fmt.Errorf("listing DataShare %q snapshot schedule (Resource Group %q / accountName %q): %+v", name, accountId.ResourceGroup, accountId.Name, err)
		}
		if syncIterator.NotDone() {
			return fmt.Errorf("more than one DataShare %q snapshot schedule (Resource Group %q / accountName %q) is returned", name, accountId.ResourceGroup, accountId.Name)
		}
	}

	return nil
}
