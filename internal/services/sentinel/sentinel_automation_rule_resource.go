// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSentinelAutomationRule() *pluginsdk.Resource {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: automationrules.ValidateWorkspaceID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"order": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"triggers_on": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(automationrules.TriggersOnIncidents),
			ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForTriggersOn(), false),
		},

		"triggers_when": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(automationrules.TriggersWhenCreated),
			ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForTriggersWhen(), false),
		},

		"expiration": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"condition_json": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// We can't use the pluginsdk.SuppressJsonDiff here as the "condition_json" is always an array, while that function assume its input is an object.
			// Once https://github.com/hashicorp/terraform-plugin-sdk/pull/1102 is merged, we can switch to pluginsdk.SuppressJsonDiff.
			DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
				return utils.NormalizeJson(old) == utils.NormalizeJson(new)
			},
			ValidateFunc: validation.StringIsJSON,
		},

		"action_incident": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"status": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForIncidentStatus(), false),
					},

					"classification": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(automationrules.IncidentClassificationUndetermined),
							string(automationrules.IncidentClassificationBenignPositive) + "_" + string(automationrules.IncidentClassificationReasonSuspiciousButExpected),
							string(automationrules.IncidentClassificationFalsePositive) + "_" + string(automationrules.IncidentClassificationReasonIncorrectAlertLogic),
							string(automationrules.IncidentClassificationFalsePositive) + "_" + string(automationrules.IncidentClassificationReasonInaccurateData),
							string(automationrules.IncidentClassificationTruePositive) + "_" + string(automationrules.IncidentClassificationReasonSuspiciousActivity),
						}, false),
					},

					"classification_comment": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"labels": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"owner_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"severity": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(automationrules.PossibleValuesForIncidentSeverity(), false),
					},
				},
			},
			AtLeastOneOf: []string{"action_incident", "action_playbook"},
		},

		"action_playbook": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"logic_app_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workflows.ValidateWorkflowID,
					},

					"tenant_id": {
						Type: pluginsdk.TypeString,
						// NOTE: O+C We'll use the current tenant id if this property is absent.
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
			AtLeastOneOf: []string{"action_incident", "action_playbook"},
		},
	}
	return &pluginsdk.Resource{
		Create: resourceSentinelAutomationRuleCreateOrUpdate,
		Read:   resourceSentinelAutomationRuleRead,
		Update: resourceSentinelAutomationRuleCreateOrUpdate,
		Delete: resourceSentinelAutomationRuleDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SentinelAutomationRuleV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: schema,
	}
}

func resourceSentinelAutomationRuleCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	tenantId := meta.(*clients.Client).Account.TenantId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceId, err := automationrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := automationrules.NewAutomationRuleID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_sentinel_automation_rule", id.ID())
		}
	}

	actions, err := expandAutomationRuleActions(d, tenantId)
	if err != nil {
		return err
	}
	params := automationrules.AutomationRule{
		Properties: automationrules.AutomationRuleProperties{
			DisplayName: d.Get("display_name").(string),
			Order:       int64(d.Get("order").(int)),
			TriggeringLogic: automationrules.AutomationRuleTriggeringLogic{
				IsEnabled:    d.Get("enabled").(bool),
				TriggersOn:   automationrules.TriggersOn(d.Get("triggers_on").(string)),
				TriggersWhen: automationrules.TriggersWhen(d.Get("triggers_when").(string)),
			},
			Actions: actions,
		},
	}

	conditions, err := expandAutomationRuleConditionsFromJSON(d.Get("condition_json").(string))
	if err != nil {
		return fmt.Errorf("expanding `condition_json`: %v", err)
	}
	params.Properties.TriggeringLogic.Conditions = conditions

	if expiration := d.Get("expiration").(string); expiration != "" {
		t, _ := time.Parse(time.RFC3339, expiration)
		params.Properties.TriggeringLogic.SetExpirationTimeUtcAsTime(t)
	}

	_, err = client.CreateOrUpdate(ctx, id, params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAutomationRuleRead(d, meta)
}

func resourceSentinelAutomationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationrules.ParseAutomationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AutomationRuleId)
	d.Set("log_analytics_workspace_id", automationrules.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())
	if model := resp.Model; model != nil {
		prop := model.Properties
		d.Set("display_name", prop.DisplayName)
		d.Set("order", prop.Order)

		tl := prop.TriggeringLogic
		d.Set("enabled", tl.IsEnabled)
		d.Set("triggers_on", string(tl.TriggersOn))
		d.Set("triggers_when", string(tl.TriggersWhen))
		d.Set("expiration", tl.ExpirationTimeUtc)

		conditionJSON, err := flattenAutomationRuleConditionsToJSON(tl.Conditions)
		if err != nil {
			return fmt.Errorf("flattening `condition_json`: %v", err)
		}
		d.Set("condition_json", conditionJSON)

		actionIncident, actionPlaybook := flattenAutomationRuleActions(prop.Actions)

		if err := d.Set("action_incident", actionIncident); err != nil {
			return fmt.Errorf("setting `action_incident`: %v", err)
		}
		if err := d.Set("action_playbook", actionPlaybook); err != nil {
			return fmt.Errorf("setting `action_playbook`: %v", err)
		}
	}

	return nil
}

func resourceSentinelAutomationRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationrules.ParseAutomationRuleID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandAutomationRuleConditionsFromJSON(input string) (*[]automationrules.AutomationRuleCondition, error) {
	if input == "" {
		return nil, nil
	}
	triggerLogic := &automationrules.AutomationRuleTriggeringLogic{}
	err := triggerLogic.UnmarshalJSON([]byte(fmt.Sprintf(`{ "conditions": %s }`, input)))
	if err != nil {
		return nil, err
	}
	return triggerLogic.Conditions, nil
}

func flattenAutomationRuleConditionsToJSON(input *[]automationrules.AutomationRuleCondition) (string, error) {
	if input == nil || len(*input) == 0 {
		return "", nil
	}
	result, err := json.Marshal(input)
	return string(result), err
}

func expandAutomationRuleActions(d *pluginsdk.ResourceData, defaultTenantId string) ([]automationrules.AutomationRuleAction, error) {
	actionIncident, err := expandAutomationRuleActionIncident(d.Get("action_incident").([]interface{}))
	if err != nil {
		return nil, err
	}
	actionPlaybook := expandAutomationRuleActionPlaybook(d.Get("action_playbook").([]interface{}), defaultTenantId)

	if len(actionIncident)+len(actionPlaybook) == 0 {
		return nil, nil
	}

	out := make([]automationrules.AutomationRuleAction, 0, len(actionIncident)+len(actionPlaybook))
	out = append(out, actionIncident...)
	out = append(out, actionPlaybook...)
	return out, nil
}

func flattenAutomationRuleActions(input []automationrules.AutomationRuleAction) (actionIncident []interface{}, actionPlaybook []interface{}) {
	actionIncident = make([]interface{}, 0)
	actionPlaybook = make([]interface{}, 0)

	for _, action := range input {
		switch action := action.(type) {
		case automationrules.AutomationRuleModifyPropertiesAction:
			actionIncident = append(actionIncident, flattenAutomationRuleActionIncident(action))
		case automationrules.AutomationRuleRunPlaybookAction:
			actionPlaybook = append(actionPlaybook, flattenAutomationRuleActionPlaybook(action))
		}
	}

	return
}

func expandAutomationRuleActionIncident(input []interface{}) ([]automationrules.AutomationRuleAction, error) {
	if len(input) == 0 {
		return nil, nil
	}

	out := make([]automationrules.AutomationRuleAction, 0, len(input))
	for _, b := range input {
		b := b.(map[string]interface{})

		status := automationrules.IncidentStatus(b["status"].(string))
		l := strings.Split(b["classification"].(string), "_")
		classification, clr := l[0], ""
		if len(l) == 2 {
			clr = l[1]
		}
		classificationComment := b["classification_comment"].(string)

		// sanity check on classification
		if status == automationrules.IncidentStatusClosed && classification == "" {
			return nil, fmt.Errorf("`classification` is required when `status` is set to `Closed`")
		}
		if status != automationrules.IncidentStatusClosed {
			if classification != "" {
				return nil, fmt.Errorf("`classification` can't be set when `status` is not set to `Closed`")
			}
			if classificationComment != "" {
				return nil, fmt.Errorf("`classification_comment` can't be set when `status` is not set to `Closed`")
			}
		}

		var labelsPtr *[]automationrules.IncidentLabel
		if labelStrsPtr := utils.ExpandStringSlice(b["labels"].([]interface{})); labelStrsPtr != nil && len(*labelStrsPtr) > 0 {
			labels := make([]automationrules.IncidentLabel, 0, len(*labelStrsPtr))
			for _, label := range *labelStrsPtr {
				labels = append(labels, automationrules.IncidentLabel{
					LabelName: label,
				})
			}
			labelsPtr = &labels
		}

		var ownerPtr *automationrules.IncidentOwnerInfo
		if ownerIdStr := b["owner_id"].(string); ownerIdStr != "" {
			ownerPtr = &automationrules.IncidentOwnerInfo{
				ObjectId: utils.String(ownerIdStr),
			}
		}

		severity := b["severity"].(string)

		// sanity check on the whole incident action
		if severity == "" && ownerPtr == nil && labelsPtr == nil && status == "" {
			return nil, fmt.Errorf("at least one of `severity`, `owner_id`, `labels` or `status` should be specified")
		}

		classificationPtr := automationrules.IncidentClassification(classification)
		clrPtr := automationrules.IncidentClassificationReason(clr)
		severityPtr := automationrules.IncidentSeverity(severity)
		out = append(out, automationrules.AutomationRuleModifyPropertiesAction{
			Order: int64(b["order"].(int)),
			ActionConfiguration: &automationrules.IncidentPropertiesAction{
				Status:                &status,
				Classification:        &classificationPtr,
				ClassificationComment: &classificationComment,
				ClassificationReason:  &clrPtr,
				Labels:                labelsPtr,
				Owner:                 ownerPtr,
				Severity:              &severityPtr,
			},
		})
	}

	return out, nil
}

func flattenAutomationRuleActionIncident(input automationrules.AutomationRuleModifyPropertiesAction) map[string]interface{} {
	var (
		status      string
		clsf        string
		clsfComment string
		clsfReason  string
		labels      []interface{}
		owner       string
		severity    string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		if cfg.Status != nil {
			status = string(*cfg.Status)
		}
		if cfg.Classification != nil {
			clsf = string(*cfg.Classification)
		}
		if cfg.ClassificationComment != nil {
			clsfComment = *cfg.ClassificationComment
		}
		if cfg.ClassificationReason != nil {
			clsfReason = string(*cfg.ClassificationReason)
		}

		if cfg.Labels != nil {
			for _, label := range *cfg.Labels {
				labels = append(labels, label.LabelName)
			}
		}

		if cfg.Owner != nil && cfg.Owner.ObjectId != nil {
			owner = *cfg.Owner.ObjectId
		}

		if cfg.Severity != nil {
			severity = string(*cfg.Severity)
		}
	}

	classification := clsf
	if clsfReason != "" {
		classification = classification + "_" + clsfReason
	}

	return map[string]interface{}{
		"order":                  input.Order,
		"status":                 status,
		"classification":         classification,
		"classification_comment": clsfComment,
		"labels":                 labels,
		"owner_id":               owner,
		"severity":               severity,
	}
}

func expandAutomationRuleActionPlaybook(input []interface{}, defaultTenantId string) []automationrules.AutomationRuleAction {
	out := make([]automationrules.AutomationRuleAction, 0, len(input))
	for _, b := range input {
		b := b.(map[string]interface{})

		tid := defaultTenantId
		if t := b["tenant_id"].(string); t != "" {
			tid = t
		}

		out = append(out, automationrules.AutomationRuleRunPlaybookAction{
			Order: int64(b["order"].(int)),
			ActionConfiguration: &automationrules.PlaybookActionProperties{
				LogicAppResourceId: b["logic_app_id"].(string),
				TenantId:           &tid,
			},
		})
	}
	return out
}

func flattenAutomationRuleActionPlaybook(input automationrules.AutomationRuleRunPlaybookAction) map[string]interface{} {
	var (
		logicAppId string
		tenantId   string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		logicAppId = cfg.LogicAppResourceId

		if cfg.TenantId != nil {
			tenantId = *cfg.TenantId
		}
	}

	return map[string]interface{}{
		"order":        input.Order,
		"logic_app_id": logicAppId,
		"tenant_id":    tenantId,
	}
}
