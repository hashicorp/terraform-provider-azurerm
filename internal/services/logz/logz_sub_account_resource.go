package logz

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogzSubAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogzSubAccountCreate,
		Read:   resourceLogzSubAccountRead,
		Update: resourceLogzSubAccountUpdate,
		Delete: resourceLogzSubAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LogzSubAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogzMonitorName,
			},

			"logz_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogzMonitorID,
			},

			"user": SchemaUserInfo(),

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceLogzSubAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := parse.LogzMonitorID(d.Get("logz_monitor_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLogzSubAccountID(monitorId.SubscriptionId, monitorId.ResourceGroup, monitorId.MonitorName, d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.AccountName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_logz_sub_account", id.ID())
	}

	monitoringStatus := logz.MonitoringStatusDisabled
	if d.Get("enabled").(bool) {
		monitoringStatus = logz.MonitoringStatusEnabled
	}

	monitorClient := meta.(*clients.Client).Logz.MonitorClient
	resp, err := monitorClient.Get(ctx, monitorId.ResourceGroup, monitorId.MonitorName)
	if err != nil {
		return fmt.Errorf("checking for existing %s: %+v", monitorId, err)
	}

	props := logz.MonitorResource{
		Location: resp.Location,
		Properties: &logz.MonitorProperties{
			UserInfo:         expandUserInfo(d.Get("user").([]interface{})),
			MonitoringStatus: monitoringStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if properties := resp.Properties; properties != nil {
		props.Properties.PlanData = properties.PlanData
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.MonitorName, id.AccountName, &props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of the %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogzSubAccountRead(d, meta)
}

func resourceLogzSubAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzSubAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.AccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] logz %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("logz_monitor_id", parse.NewLogzMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName).ID())

	if props := resp.Properties; props != nil {
		d.Set("enabled", props.MonitoringStatus == logz.MonitoringStatusEnabled)
		if err := d.Set("user", flattenUserInfo(expandUserInfo(d.Get("user").([]interface{})))); err != nil {
			return fmt.Errorf("setting `user`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogzSubAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzSubAccountID(d.Id())
	if err != nil {
		return err
	}

	props := logz.MonitorResourceUpdateParameters{
		Properties: &logz.MonitorUpdateProperties{},
	}

	if d.HasChange("enabled") {
		monitoringStatus := logz.MonitoringStatusDisabled
		if d.Get("enabled").(bool) {
			monitoringStatus = logz.MonitoringStatusEnabled
		}
		props.Properties.MonitoringStatus = monitoringStatus
	}

	if d.HasChange("tags") {
		props.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.MonitorName, id.AccountName, &props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceLogzSubAccountRead(d, meta)
}

func resourceLogzSubAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzSubAccountID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.MonitorName, id.AccountName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the %s: %+v", id, err)
	}

	// API has bug, which appears to be eventually consistent. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/18572
	log.Printf("[DEBUG] Waiting for %s to be fully deleted..", *id)
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   subAccountDeletedRefreshFunc(ctx, client, *id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 20,
		Timeout:                   time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
	}

	return nil
}

func subAccountDeletedRefreshFunc(ctx context.Context, client *logz.SubAccountClient, id parse.LogzSubAccountId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.AccountName)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}
