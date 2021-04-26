package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataShareCreateUpdate,
		Read:   resourceDataShareRead,
		Update: resourceDataShareCreateUpdate,
		Delete: resourceDataShareDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ShareID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ShareName(),
			},

			"account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountID,
			},

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(datashare.CopyBased),
					string(datashare.InPlace),
				}, false),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"snapshot_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.SnapshotScheduleName(),
						},

						"recurrence": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(datashare.Day),
								string(datashare.Hour),
							}, false),
						},

						"start_time": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validation.IsRFC3339Time,
							DiffSuppressFunc: suppress.RFC3339Time,
						},
					},
				},
			},

			"terms": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDataShareCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountId, err := parse.AccountID(d.Get("account_id").(string))
	if err != nil {
		return err
	}

	resourceId := parse.NewShareID(subscriptionId, accountId.ResourceGroup, accountId.Name, name).ID()
	if d.IsNewResource() {
		existing, err := client.Get(ctx, accountId.ResourceGroup, accountId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing DataShare %q (Resource Group %q / accountName %q): %+v", name, accountId.ResourceGroup, accountId.Name, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_share", resourceId)
		}
	}

	share := datashare.Share{
		ShareProperties: &datashare.ShareProperties{
			ShareKind:   datashare.ShareKind(d.Get("kind").(string)),
			Description: utils.String(d.Get("description").(string)),
			Terms:       utils.String(d.Get("terms").(string)),
		},
	}

	if _, err := client.Create(ctx, accountId.ResourceGroup, accountId.Name, name, share); err != nil {
		return fmt.Errorf("creating Data Share %q (Account %q / Resource Group %q): %+v", name, accountId.Name, accountId.ResourceGroup, err)
	}

	d.SetId(resourceId)

	if d.HasChange("snapshot_schedule") {
		// only one dependent sync setting is allowed in one data share
		o, _ := d.GetChange("snapshot_schedule")
		if origins := o.([]interface{}); len(origins) > 0 {
			origin := origins[0].(map[string]interface{})
			if originName, ok := origin["name"].(string); ok && originName != "" {
				future, err := syncClient.Delete(ctx, accountId.ResourceGroup, accountId.Name, name, originName)
				if err != nil {
					return fmt.Errorf("deleting DataShare %q snapshot schedule (Resource Group %q / accountName %q): %+v", name, accountId.ResourceGroup, accountId.Name, err)
				}
				if err = future.WaitForCompletionRef(ctx, syncClient.Client); err != nil {
					return fmt.Errorf("waiting for DataShare %q snapshot schedule (Resource Group %q / accountName %q) to be deleted: %+v", name, accountId.ResourceGroup, accountId.Name, err)
				}
			}
		}
	}

	if snapshotSchedule := expandAzureRmDataShareSnapshotSchedule(d.Get("snapshot_schedule").([]interface{})); snapshotSchedule != nil {
		if _, err := syncClient.Create(ctx, accountId.ResourceGroup, accountId.Name, name, d.Get("snapshot_schedule.0.name").(string), snapshotSchedule); err != nil {
			return fmt.Errorf("creating DataShare %q snapshot schedule (Resource Group %q / accountName %q): %+v", name, accountId.ResourceGroup, accountId.Name, err)
		}
	}

	return resourceDataShareRead(d, meta)
}

func resourceDataShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ShareID(d.Id())
	if err != nil {
		return err
	}

	dataShare, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(dataShare.Response) {
			log.Printf("[INFO] DataShare %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataShare %q (Resource Group %q / accountName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, err)
	}

	accountId := parse.NewAccountID(subscriptionId, id.ResourceGroup, id.AccountName)

	d.Set("name", id.Name)
	d.Set("account_id", accountId.ID())

	if props := dataShare.ShareProperties; props != nil {
		d.Set("kind", props.ShareKind)
		d.Set("description", props.Description)
		d.Set("terms", props.Terms)
	}

	settings := make([]datashare.ScheduledSynchronizationSetting, 0)
	syncIterator, err := syncClient.ListByShareComplete(ctx, id.ResourceGroup, id.AccountName, id.Name, "")
	if err != nil {
		return fmt.Errorf("listing Snapshot Schedules for Data Share %q (Account %q / Resource Group %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
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
	if err := d.Set("snapshot_schedule", flattenAzureRmDataShareSnapshotSchedule(settings)); err != nil {
		return fmt.Errorf("setting `snapshot_schedule`: %+v", err)
	}

	return nil
}

func resourceDataShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ShareID(d.Id())
	if err != nil {
		return err
	}

	// sync setting will not automatically be deleted after the data share is deleted
	if _, ok := d.GetOk("snapshot_schedule"); ok {
		syncFuture, err := syncClient.Delete(ctx, id.ResourceGroup, id.AccountName, id.Name, d.Get("snapshot_schedule.0.name").(string))
		if err != nil {
			return fmt.Errorf("deleting DataShare %q snapshot schedule (Resource Group %q / accountName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, err)
		}
		if err = syncFuture.WaitForCompletionRef(ctx, syncClient.Client); err != nil {
			return fmt.Errorf("waiting for DataShare %q snapshot schedule (Resource Group %q / accountName %q) to be deleted: %+v", id.Name, id.ResourceGroup, id.AccountName, err)
		}
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting DataShare %q (Resource Group %q / accountName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for DataShare %q (Resource Group %q / accountName %q) to be deleted: %+v", id.Name, id.ResourceGroup, id.AccountName, err)
	}

	return nil
}

func expandAzureRmDataShareSnapshotSchedule(input []interface{}) *datashare.ScheduledSynchronizationSetting {
	if len(input) == 0 {
		return nil
	}

	snapshotSchedule := input[0].(map[string]interface{})

	startTime, _ := time.Parse(time.RFC3339, snapshotSchedule["start_time"].(string))

	return &datashare.ScheduledSynchronizationSetting{
		Kind: datashare.KindBasicSynchronizationSettingKindScheduleBased,
		ScheduledSynchronizationSettingProperties: &datashare.ScheduledSynchronizationSettingProperties{
			RecurrenceInterval:  datashare.RecurrenceInterval(snapshotSchedule["recurrence"].(string)),
			SynchronizationTime: &date.Time{Time: startTime},
		},
	}
}

func flattenAzureRmDataShareSnapshotSchedule(input []datashare.ScheduledSynchronizationSetting) []interface{} {
	output := make([]interface{}, 0)

	for _, setting := range input {
		name := ""
		if setting.Name != nil {
			name = *setting.Name
		}

		startTime := ""
		if setting.SynchronizationTime != nil && !setting.SynchronizationTime.IsZero() {
			startTime = setting.SynchronizationTime.Format(time.RFC3339)
		}

		output = append(output, map[string]interface{}{
			"name":       name,
			"recurrence": string(setting.RecurrenceInterval),
			"start_time": startTime,
		})
	}

	return output
}
