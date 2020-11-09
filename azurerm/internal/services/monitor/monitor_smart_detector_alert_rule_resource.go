package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-06-01-preview/alertsmanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	commonValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorSmartDetectorAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorSmartDetectorAlertRuleCreateUpdate,
		Read:   resourceArmMonitorSmartDetectorAlertRuleRead,
		Update: resourceArmMonitorSmartDetectorAlertRuleCreateUpdate,
		Delete: resourceArmMonitorSmartDetectorAlertRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
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

			"detector_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"FailureAnomaliesDetector",
				}, false),
			},

			"scope": {
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

						"custom_email_subject": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"custom_webhook_payload": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"throttling": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: commonValidate.ISO8601Duration,
						},
					},
				},
			},
		},
	}
}

func resourceArmMonitorSmartDetectorAlertRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
				ID: utils.String(d.Get("detector_id").(string)),
			},
			Scope:        utils.ExpandStringSlice(d.Get("scope").(*schema.Set).List()),
			ActionGroups: expandArmMonitorSmartDetectorAlertRuleActionGroup(d.Get("action_group").([]interface{})),
			Throttling:   expandArmMonitorSmartDetectorAlertRuleThrottling(d.Get("throttling").([]interface{})),
		},
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
	return resourceArmMonitorSmartDetectorAlertRuleRead(d, meta)
}

func resourceArmMonitorSmartDetectorAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
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
		d.Set("scope", utils.FlattenStringSlice(props.Scope))

		if props.Detector != nil {
			d.Set("detector_id", props.Detector.ID)
		}

		if err := d.Set("action_group", flattenArmMonitorSmartDetectorAlertRuleActionGroup(props.ActionGroups)); err != nil {
			return fmt.Errorf("setting `action_group`: %+v", err)
		}
		if err := d.Set("throttling", flattenArmMonitorSmartDetectorAlertRuleThrottling(props.Throttling)); err != nil {
			return fmt.Errorf("setting `throttling`: %+v", err)
		}
	}

	return nil
}

func resourceArmMonitorSmartDetectorAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandArmMonitorSmartDetectorAlertRuleActionGroup(input []interface{}) *alertsmanagement.ActionGroupsInformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &alertsmanagement.ActionGroupsInformation{
		CustomEmailSubject:   utils.String(v["custom_email_subject"].(string)),
		CustomWebhookPayload: utils.String(v["custom_webhook_payload"].(string)),
		GroupIds:             utils.ExpandStringSlice(v["ids"].(*schema.Set).List()),
	}
}

func expandArmMonitorSmartDetectorAlertRuleThrottling(input []interface{}) *alertsmanagement.ThrottlingInformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &alertsmanagement.ThrottlingInformation{
		Duration: utils.String(v["duration"].(string)),
	}
}

func flattenArmMonitorSmartDetectorAlertRuleActionGroup(input *alertsmanagement.ActionGroupsInformation) []interface{} {
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
			"ids":                    utils.FlattenStringSlice(input.GroupIds),
			"custom_email_subject":   customEmailSubject,
			"custom_webhook_payload": CustomWebhookPayload,
		},
	}
}

func flattenArmMonitorSmartDetectorAlertRuleThrottling(input *alertsmanagement.ThrottlingInformation) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	var duration string
	if input.Duration != nil {
		duration = *input.Duration
	}
	return []interface{}{
		map[string]interface{}{
			"duration": duration,
		},
	}
}
