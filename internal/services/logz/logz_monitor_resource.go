// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logz

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogzMonitor() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogzMonitorCreate,
		Read:   resourceLogzMonitorRead,
		Update: resourceLogzMonitorUpdate,
		Delete: resourceLogzMonitorDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := monitors.ParseMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogzMonitorName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"company_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enterprise_app_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"single_sign_on_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"logz_organization_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"plan": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"billing_cycle": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"MONTHLY",
								"WEEKLY",
							}, false),
						},

						"effective_date": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.RFC3339Time,
							ValidateFunc:     validation.IsRFC3339Time,
						},

						"plan_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								PlanId100gb14days,
							}, false),
							Default: PlanId100gb14days,
						},

						"usage_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PAYG",
								"COMMITTED",
							}, false),
						},
					},
				},
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

func resourceLogzMonitorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.MonitorClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := monitors.NewMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_logz_monitor", id.ID())
	}

	monitoringStatus := monitors.MonitoringStatusDisabled
	if d.Get("enabled").(bool) {
		monitoringStatus = monitors.MonitoringStatusEnabled
	}

	planData, err := expandMonitorPlanData(d.Get("plan").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `plan` %s: %+v", id, err)
	}

	payload := monitors.LogzMonitorResource{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &monitors.MonitorProperties{
			LogzOrganizationProperties: expandMonitorOrganizationProperties(d),
			PlanData:                   planData,
			UserInfo:                   expandMonitorUserInfo(d.Get("user").([]interface{})),
			MonitoringStatus:           pointer.To(monitoringStatus),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogzMonitorRead(d, meta)
}

func resourceLogzMonitorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.MonitorClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitors.ParseMonitorID(d.Id())
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

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			if org := props.LogzOrganizationProperties; org != nil {
				d.Set("company_name", org.CompanyName)
				d.Set("enterprise_app_id", org.EnterpriseAppId)
				d.Set("single_sign_on_url", org.SingleSignOnUrl)
				d.Set("logz_organization_id", org.Id)
			}

			d.Set("enabled", props.MonitoringStatus != nil && *props.MonitoringStatus == monitors.MonitoringStatusEnabled)

			planData, err := flattenMonitorPlanData(props.PlanData)
			if err != nil {
				return fmt.Errorf("flatten `plan`: %+v", err)
			}

			if err := d.Set("plan", planData); err != nil {
				return fmt.Errorf("setting `plan`: %+v", err)
			}

			if err := d.Set("user", flattenMonitorUserInfo(props.UserInfo)); err != nil {
				return fmt.Errorf("setting `user`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceLogzMonitorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.MonitorClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitors.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	payload := monitors.LogzMonitorResourceUpdateParameters{
		Properties: &monitors.MonitorUpdateProperties{},
	}

	if d.HasChange("enabled") {
		monitoringStatus := monitors.MonitoringStatusDisabled
		if d.Get("enabled").(bool) {
			monitoringStatus = monitors.MonitoringStatusEnabled
		}
		payload.Properties.MonitoringStatus = pointer.To(monitoringStatus)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceLogzMonitorRead(d, meta)
}

func resourceLogzMonitorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.MonitorClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitors.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandMonitorOrganizationProperties(d *pluginsdk.ResourceData) *monitors.LogzOrganizationProperties {
	props := &monitors.LogzOrganizationProperties{}
	companyName := d.Get("company_name").(string)
	if companyName != "" {
		props.CompanyName = utils.String(companyName)
	}

	enterpriseAppID := d.Get("enterprise_app_id").(string)
	if enterpriseAppID != "" {
		props.EnterpriseAppId = utils.String(enterpriseAppID)
	}

	return props
}

func expandMonitorPlanData(input []interface{}) (*monitors.PlanData, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]interface{})
	effectiveDate, _ := time.Parse(time.RFC3339, v["effective_date"].(string))
	planDetails, err := getPlanDetails(v["plan_id"].(string))
	if err != nil {
		return nil, err
	}

	data := &monitors.PlanData{
		UsageType:    utils.String(v["usage_type"].(string)),
		BillingCycle: utils.String(v["billing_cycle"].(string)),
		PlanDetails:  &planDetails,
	}
	data.SetEffectiveDateAsTime(effectiveDate)
	return data, nil
}

func flattenMonitorPlanData(input *monitors.PlanData) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), nil
	}

	var billingCycle string
	if input.BillingCycle != nil {
		billingCycle = *input.BillingCycle
	}

	var effectiveDate string
	date, err := input.GetEffectiveDateAsTime()
	if err != nil {
		return nil, fmt.Errorf("parsing EffectiveDate: %+v", err)
	}
	if date != nil {
		effectiveDate = date.Format(time.RFC3339)
	}

	var planId string
	if input.PlanDetails != nil {
		var err error
		planId, err = getPlanId(*input.PlanDetails)
		if err != nil {
			return nil, err
		}
	}

	var usageType string
	if input.UsageType != nil {
		usageType = *input.UsageType
	}

	return []interface{}{
		map[string]interface{}{
			"billing_cycle":  billingCycle,
			"effective_date": effectiveDate,
			"plan_id":        planId,
			"usage_type":     usageType,
		},
	}, nil
}

func expandMonitorUserInfo(input []interface{}) *monitors.UserInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &monitors.UserInfo{
		FirstName:    utils.String(v["first_name"].(string)),
		LastName:     utils.String(v["last_name"].(string)),
		EmailAddress: utils.String(v["email"].(string)),
		PhoneNumber:  utils.String(v["phone_number"].(string)),
	}
}

func flattenMonitorUserInfo(input *monitors.UserInfo) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
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
