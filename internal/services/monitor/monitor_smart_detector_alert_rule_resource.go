// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-06-01-preview/alertsmanagement" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorSmartDetectorAlertRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorSmartDetectorAlertRuleCreateUpdate,
		Read:   resourceMonitorSmartDetectorAlertRuleRead,
		Update: resourceMonitorSmartDetectorAlertRuleCreateUpdate,
		Delete: resourceMonitorSmartDetectorAlertRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SmartDetectorAlertRuleV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SmartDetectorAlertRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"detector_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"FailureAnomaliesDetector",
					"RequestPerformanceDegradationDetector",
					"DependencyPerformanceDegradationDetector",
					"ExceptionVolumeChangedDetector",
					"TraceSeverityDetector",
					"MemoryLeakDetector",
				}, false),
			},

			"scope_resource_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
				Set: set.HashStringIgnoreCase,
			},

			"severity": {
				Type:     pluginsdk.TypeString,
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonValidate.ISO8601Duration,
			},

			"action_group": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ids": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.ActionGroupID,
							},
							Set: set.HashStringIgnoreCase,
						},

						"email_subject": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"webhook_payload": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
						},
					},
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"throttling_duration": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonValidate.ISO8601Duration,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorSmartDetectorAlertRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSmartDetectorAlertRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, utils.Bool(true))
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_monitor_smart_detector_alert_rule", id.ID())
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
			Scope:        utils.ExpandStringSlice(d.Get("scope_resource_ids").(*pluginsdk.Set).List()),
			ActionGroups: expandMonitorSmartDetectorAlertRuleActionGroup(d.Get("action_group").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("throttling_duration"); ok {
		actionRule.AlertRuleProperties.Throttling = &alertsmanagement.ThrottlingInformation{
			Duration: utils.String(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, actionRule); err != nil {
		return fmt.Errorf("creating/updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorSmartDetectorAlertRuleRead(d, meta)
}

func resourceMonitorSmartDetectorAlertRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("retrieving Monitor %s: %+v", *id, err)
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

func resourceMonitorSmartDetectorAlertRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SmartDetectorAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
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
		GroupIds:             utils.ExpandStringSlice(v["ids"].(*pluginsdk.Set).List()),
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
