package logz

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			_, err := parse.LogzMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogzMonitorName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

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
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"100gb14days",
							}, false),
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

func resourceLogzMonitorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.MonitorClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLogzMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_logz_monitor", id.ID())
	}

	monitoringStatus := logz.MonitoringStatusDisabled
	if d.Get("enabled").(bool) {
		monitoringStatus = logz.MonitoringStatusEnabled
	}

	props := logz.MonitorResource{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &logz.MonitorProperties{
			LogzOrganizationProperties: expandMonitorOrganizationProperties(d),
			PlanData:                   expandMonitorPlanData(d.Get("plan").([]interface{})),
			UserInfo:                   expandUserInfo(d.Get("user").([]interface{})),
			MonitoringStatus:           monitoringStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.MonitorName, &props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of the %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogzMonitorRead(d, meta)
}

func resourceLogzMonitorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.MonitorClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] logz %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		flattenMonitorOrganizationProperties(d, props.LogzOrganizationProperties)
		d.Set("enabled", props.MonitoringStatus == logz.MonitoringStatusEnabled)
		if err := d.Set("plan", flattenMonitorPlanData(props.PlanData)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		if err := d.Set("user", flattenUserInfo(expandUserInfo(d.Get("user").([]interface{})))); err != nil {
			return fmt.Errorf("setting `user`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogzMonitorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.MonitorClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzMonitorID(d.Id())
	if err != nil {
		return err
	}

	body := logz.MonitorResourceUpdateParameters{
		Properties: &logz.MonitorUpdateProperties{},
	}

	if d.HasChange("enabled") {
		monitoringStatus := logz.MonitoringStatusDisabled
		if d.Get("enabled").(bool) {
			monitoringStatus = logz.MonitoringStatusEnabled
		}
		body.Properties.MonitoringStatus = monitoringStatus
	}

	if d.HasChange("tags") {
		body.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.MonitorName, &body); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceLogzMonitorRead(d, meta)
}

func resourceLogzMonitorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.MonitorClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzMonitorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the %s: %+v", id, err)
	}

	return nil
}

func expandMonitorOrganizationProperties(d *pluginsdk.ResourceData) *logz.OrganizationProperties {
	props := &logz.OrganizationProperties{}
	companyName := d.Get("company_name").(string)
	if companyName != "" {
		props.CompanyName = utils.String(companyName)
	}

	enterpriseAppID := d.Get("enterprise_app_id").(string)
	if enterpriseAppID != "" {
		props.EnterpriseAppID = utils.String(enterpriseAppID)
	}

	return props
}

func expandMonitorPlanData(input []interface{}) *logz.PlanData {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	effectiveDate, _ := time.Parse(time.RFC3339, v["effective_date"].(string))
	return &logz.PlanData{
		UsageType:     utils.String(v["usage_type"].(string)),
		BillingCycle:  utils.String(v["billing_cycle"].(string)),
		PlanDetails:   utils.String(v["plan_id"].(string)),
		EffectiveDate: &date.Time{Time: effectiveDate},
	}
}

func flattenMonitorOrganizationProperties(d *pluginsdk.ResourceData, input *logz.OrganizationProperties) {
	if input == nil {
		return
	}

	d.Set("company_name", input.CompanyName)
	d.Set("enterprise_app_id", input.EnterpriseAppID)
	d.Set("single_sign_on_url", input.SingleSignOnURL)
	d.Set("logz_organization_id", input.ID)
}

func flattenMonitorPlanData(input *logz.PlanData) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var billingCycle string
	if input.BillingCycle != nil {
		billingCycle = *input.BillingCycle
	}

	var effectiveDate string
	if input.EffectiveDate != nil {
		effectiveDate = input.EffectiveDate.Format(time.RFC3339)
	}

	var planDetails string
	if input.PlanDetails != nil {
		planDetails = *input.PlanDetails
	}

	var usageType string
	if input.UsageType != nil {
		usageType = *input.UsageType
	}

	return []interface{}{
		map[string]interface{}{
			"billing_cycle":  billingCycle,
			"effective_date": effectiveDate,
			"plan_id":        planDetails,
			"usage_type":     usageType,
		},
	}
}
