// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logz

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/subaccount"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
			_, err := subaccount.ParseAccountID(id)
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
				ValidateFunc: monitors.ValidateMonitorID,
			},

			"user": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"first_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 50),
						},

						"last_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 50),
						},

						"phone_number": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 40),
						},
					},
				},
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": commonschema.Tags(),
		},
	}
}
func resourceLogzSubAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountClient
	monitorClient := meta.(*clients.Client).Logz.MonitorClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := monitors.ParseMonitorID(d.Get("logz_monitor_id").(string))
	if err != nil {
		return err
	}

	id := subaccount.NewAccountID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_logz_sub_account", id.ID())
	}

	monitoringStatus := subaccount.MonitoringStatusDisabled
	if d.Get("enabled").(bool) {
		monitoringStatus = subaccount.MonitoringStatusEnabled
	}

	parentMonitor, err := monitorClient.Get(ctx, *monitorId)
	if err != nil {
		return fmt.Errorf("retrieving parent %s: %+v", *monitorId, err)
	}
	if parentMonitor.Model == nil {
		return fmt.Errorf("retrieving parent %s: model was nil", *monitorId)
	}

	payload := subaccount.LogzMonitorResource{
		Location: location.Normalize(parentMonitor.Model.Location),
		Properties: &subaccount.MonitorProperties{
			UserInfo:         expandSubAccountUserInfo(d.Get("user").([]interface{})),
			MonitoringStatus: pointer.To(monitoringStatus),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if props := parentMonitor.Model.Properties; props != nil {
		payload.Properties.PlanData = mapMonitorPlanDataToSubAccountPlanData(props.PlanData)
	}

	if err := client.CreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogzSubAccountRead(d, meta)
}

func resourceLogzSubAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subaccount.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("logz_monitor_id", monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("enabled", props.MonitoringStatus != nil && *props.MonitoringStatus == subaccount.MonitoringStatusEnabled)

			if err := d.Set("user", flattenSubAccountUserInfo(props.UserInfo)); err != nil {
				return fmt.Errorf("setting `user_info`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceLogzSubAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subaccount.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	payload := subaccount.LogzMonitorResourceUpdateParameters{}

	if d.HasChange("enabled") {
		monitoringStatus := subaccount.MonitoringStatusDisabled
		if d.Get("enabled").(bool) {
			monitoringStatus = subaccount.MonitoringStatusEnabled
		}
		payload.Properties = &subaccount.MonitorUpdateProperties{
			MonitoringStatus: pointer.To(monitoringStatus),
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceLogzSubAccountRead(d, meta)
}

func resourceLogzSubAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subaccount.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	// API has bug, which appears to be eventually consistent. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/18572
	log.Printf("[DEBUG] Waiting for %s to be fully deleted..", *id)
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
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

func subAccountDeletedRefreshFunc(ctx context.Context, client *subaccount.SubAccountClient, id subaccount.AccountId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}

func expandSubAccountUserInfo(input []interface{}) *subaccount.UserInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &subaccount.UserInfo{
		FirstName:    utils.String(v["first_name"].(string)),
		LastName:     utils.String(v["last_name"].(string)),
		EmailAddress: utils.String(v["email"].(string)),
		PhoneNumber:  utils.String(v["phone_number"].(string)),
	}
}

func flattenSubAccountUserInfo(input *subaccount.UserInfo) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	firstName := ""
	if input.FirstName != nil {
		firstName = *input.FirstName
	}

	lastName := ""
	if input.LastName != nil {
		lastName = *input.LastName
	}

	email := ""
	if input.EmailAddress != nil {
		email = *input.EmailAddress
	}

	phoneNumber := ""
	if input.PhoneNumber != nil {
		phoneNumber = *input.PhoneNumber
	}

	return []interface{}{
		map[string]interface{}{
			"first_name":   firstName,
			"last_name":    lastName,
			"email":        email,
			"phone_number": phoneNumber,
		},
	}
}

func mapMonitorPlanDataToSubAccountPlanData(input *monitors.PlanData) *subaccount.PlanData {
	if input == nil {
		return nil
	}

	return &subaccount.PlanData{
		BillingCycle:  input.BillingCycle,
		EffectiveDate: input.EffectiveDate,
		PlanDetails:   input.PlanDetails,
		UsageType:     input.UsageType,
	}
}
