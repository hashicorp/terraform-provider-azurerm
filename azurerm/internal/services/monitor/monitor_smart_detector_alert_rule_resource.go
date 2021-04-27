package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-06-01-preview/alertsmanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	commonValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMonitorSmartDetectorAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorSmartDetectorAlertRuleCreateUpdate,
		Read:   resourceMonitorSmartDetectorAlertRuleRead,
		Update: resourceMonitorSmartDetectorAlertRuleCreateUpdate,
		Delete: resourceMonitorSmartDetectorAlertRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SmartDetectorAlertRuleID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"detector_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"FailureAnomaliesDetector",
				}, false),
			},

			"scope_resource_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
				Set: set.HashStringIgnoreCase,
			},

			"severity": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(alertsmanagement.Sev0),
						string(alertsmanagement.Sev1),
						string(alertsmanagement.Sev2),
						string(alertsmanagement.Sev3),
						string(alertsmanagement.Sev4),
					}, false),
			},

			"frequency": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: commonValidate.ISO8601Duration,
			},

			"action_group": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.ActionGroupID,
							},
							Set: set.HashStringIgnoreCase,
						},

						"email_subject": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"webhook_payload": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},
					},
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"throttling_duration": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: commonValidate.ISO8601Duration,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorSmartDetectorAlertRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, utils.Bool(true))
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Monitor Smart Detector Alert Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_smart_detector_alert_rule", *existing.ID)
		}
	}

	state := alertsmanagement.AlertRuleStateDisabled
	if d.Get("enabled").(bool) {
		state = alertsmanagement.AlertRuleStateEnabled
	}

	actionRule := alertsmanagement.AlertRule{
		// the location is always global from the portal
		Location: utils.String(location.Normalize("Global")),
		AlertRuleProperties: &alertsmanagement.AlertRuleProperties{
			Description: utils.String(d.Get("description").(string)),
			State:       state,
			Severity:    alertsmanagement.Severity(d.Get("severity").(string)),
			Frequency:   utils.String(d.Get("frequency").(string)),
			Detector: &alertsmanagement.Detector{
				ID: utils.String(d.Get("detector_type").(string)),
			},
			Scope:        utils.ExpandStringSlice(d.Get("scope_resource_ids").(*schema.Set).List()),
			ActionGroups: expandMonitorSmartDetectorAlertRuleActionGroup(d.Get("action_group").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("throttling_duration"); ok {
		actionRule.AlertRuleProperties.Throttling = &alertsmanagement.ThrottlingInformation{
			Duration: utils.String(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, actionRule); err != nil {
		return fmt.Errorf("creating/updating Monitor Smart Detector Alert Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		return fmt.Errorf("retrieving Monitor Smart Detector Alert Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Monitor Smart Detector Alert Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceMonitorSmartDetectorAlertRuleRead(d, meta)
}

func resourceMonitorSmartDetectorAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SmartDetectorAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, utils.Bool(true))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Monitor Smart Detector Alert Rule %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Monitor Smart Detector Alert Rule %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if props := resp.AlertRuleProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("enabled", props.State == alertsmanagement.AlertRuleStateEnabled)
		d.Set("frequency", props.Frequency)
		d.Set("severity", string(props.Severity))
		d.Set("scope_resource_ids", utils.FlattenStringSlice(props.Scope))

		if props.Detector != nil {
			d.Set("detector_type", props.Detector.ID)
		}

		throttlingDuration := ""
		if props.Throttling != nil && props.Throttling.Duration != nil {
			throttlingDuration = *props.Throttling.Duration
		}
		d.Set("throttling_duration", throttlingDuration)

		if err := d.Set("action_group", flattenMonitorSmartDetectorAlertRuleActionGroup(props.ActionGroups)); err != nil {
			return fmt.Errorf("setting `action_group`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMonitorSmartDetectorAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SmartDetectorAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Monitor Smart Detector Alert Rule %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}

func expandMonitorSmartDetectorAlertRuleActionGroup(input []interface{}) *alertsmanagement.ActionGroupsInformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &alertsmanagement.ActionGroupsInformation{
		CustomEmailSubject:   utils.String(v["email_subject"].(string)),
		CustomWebhookPayload: utils.String(v["webhook_payload"].(string)),
		GroupIds:             utils.ExpandStringSlice(v["ids"].(*schema.Set).List()),
	}
}

func flattenMonitorSmartDetectorAlertRuleActionGroup(input *alertsmanagement.ActionGroupsInformation) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var customEmailSubject, CustomWebhookPayload string
	if input.CustomEmailSubject != nil {
		customEmailSubject = *input.CustomEmailSubject
	}
	if input.CustomWebhookPayload != nil {
		CustomWebhookPayload = *input.CustomWebhookPayload
	}

	return []interface{}{
		map[string]interface{}{
			"ids":             utils.FlattenStringSlice(input.GroupIds),
			"email_subject":   customEmailSubject,
			"webhook_payload": CustomWebhookPayload,
		},
	}
}
