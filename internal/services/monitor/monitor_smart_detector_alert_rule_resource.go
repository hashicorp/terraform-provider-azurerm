// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
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
			_, err := smartdetectoralertrules.ParseSmartDetectorAlertRuleID(id)
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
						string(smartdetectoralertrules.SeveritySevZero),
						string(smartdetectoralertrules.SeveritySevOne),
						string(smartdetectoralertrules.SeveritySevTwo),
						string(smartdetectoralertrules.SeveritySevThree),
						string(smartdetectoralertrules.SeveritySevFour),
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMonitorSmartDetectorAlertRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := smartdetectoralertrules.NewSmartDetectorAlertRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, smartdetectoralertrules.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_smart_detector_alert_rule", id.ID())
		}
	}

	state := smartdetectoralertrules.AlertRuleStateDisabled
	if d.Get("enabled").(bool) {
		state = smartdetectoralertrules.AlertRuleStateEnabled
	}

	actionRule := smartdetectoralertrules.AlertRule{
		// the location is always global from the portal
		Location: utils.String(location.Normalize("Global")),
		Properties: &smartdetectoralertrules.AlertRuleProperties{
			Description: pointer.To(d.Get("description").(string)),
			State:       state,
			Severity:    smartdetectoralertrules.Severity(d.Get("severity").(string)),
			Frequency:   d.Get("frequency").(string),
			Detector: smartdetectoralertrules.Detector{
				Id: d.Get("detector_type").(string),
			},
			Scope:        pointer.From(utils.ExpandStringSlice(d.Get("scope_resource_ids").(*pluginsdk.Set).List())),
			ActionGroups: pointer.From(expandMonitorSmartDetectorAlertRuleActionGroup(d.Get("action_group").([]interface{}))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("throttling_duration"); ok {
		actionRule.Properties.Throttling = &smartdetectoralertrules.ThrottlingInformation{
			Duration: pointer.To(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, actionRule); err != nil {
		return fmt.Errorf("creating/updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorSmartDetectorAlertRuleRead(d, meta)
}

func resourceMonitorSmartDetectorAlertRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := smartdetectoralertrules.ParseSmartDetectorAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, smartdetectoralertrules.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SmartDetectorAlertRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("enabled", props.State == smartdetectoralertrules.AlertRuleStateEnabled)
			d.Set("frequency", props.Frequency)
			d.Set("severity", string(props.Severity))
			d.Set("scope_resource_ids", props.Scope)
			d.Set("detector_type", props.Detector.Id)

			throttlingDuration := ""
			if props.Throttling != nil && props.Throttling.Duration != nil {
				throttlingDuration = *props.Throttling.Duration
			}
			d.Set("throttling_duration", throttlingDuration)

			actionGroup, err := flattenMonitorSmartDetectorAlertRuleActionGroup(&props.ActionGroups)
			if err != nil {
				return fmt.Errorf("flatten `action_group`: %+v", err)
			}
			if err := d.Set("action_group", actionGroup); err != nil {
				return fmt.Errorf("setting `action_group`: %+v", err)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceMonitorSmartDetectorAlertRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := smartdetectoralertrules.ParseSmartDetectorAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}

func expandMonitorSmartDetectorAlertRuleActionGroup(input []interface{}) *smartdetectoralertrules.ActionGroupsInformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &smartdetectoralertrules.ActionGroupsInformation{
		CustomEmailSubject:   utils.String(v["email_subject"].(string)),
		CustomWebhookPayload: utils.String(v["webhook_payload"].(string)),
		GroupIds:             pointer.From(utils.ExpandStringSlice(v["ids"].(*pluginsdk.Set).List())),
	}
}

func flattenMonitorSmartDetectorAlertRuleActionGroup(input *smartdetectoralertrules.ActionGroupsInformation) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	var customEmailSubject, CustomWebhookPayload string
	if input.CustomEmailSubject != nil {
		customEmailSubject = *input.CustomEmailSubject
	}
	if input.CustomWebhookPayload != nil {
		CustomWebhookPayload = *input.CustomWebhookPayload
	}

	groupIds := make([]string, 0)
	for _, idRaw := range input.GroupIds {
		id, err := parse.ActionGroupIDInsensitively(idRaw)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", idRaw, err)
		}
		groupIds = append(groupIds, id.ID())
	}

	return []interface{}{
		map[string]interface{}{
			"ids":             groupIds,
			"email_subject":   customEmailSubject,
			"webhook_payload": CustomWebhookPayload,
		},
	}, nil
}
