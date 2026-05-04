// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

const SentinelAlertRuleFusionName = "BuiltInFusion"

func resourceSentinelAlertRuleFusion() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleFusionCreate,
		Read:   resourceSentinelAlertRuleFusionRead,
		Update: resourceSentinelAlertRuleFusionUpdate,
		Delete: resourceSentinelAlertRuleFusionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := alertrules.ParseAlertRuleID(id)
			return err
		}, importSentinelAlertRule(alertrules.AlertRuleKindFusion)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: alertrules.ValidateWorkspaceID,
			},

			"alert_rule_template_guid": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"source": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				// NOTE: O+C The API creates a source if omitted based on the `alert_rule_template_guid`
				// but overwriting this/reverting to the default can be done without issue so this can remain
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"sub_type": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
									},
									"severities_allowed": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(
												[]string{
													string(alertrules.AlertSeverityHigh),
													string(alertrules.AlertSeverityMedium),
													string(alertrules.AlertSeverityLow),
													string(alertrules.AlertSeverityInformational),
												},
												false,
											),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if !features.FivePointOh() {
		resource.Schema["name"] = &pluginsdk.Schema{
			Deprecated:   "the `name` is deprecated and will be removed in v5.0 version of the provider.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      SentinelAlertRuleFusionName,
			ValidateFunc: validation.StringIsNotEmpty,
		}
	}
	return resource
}

func resourceSentinelAlertRuleFusionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := SentinelAlertRuleFusionName
	if !features.FivePointOh() {
		name = d.Get("name").(string)
	}

	workspaceID, err := alertrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	// The only one fusion alert is enabled by default, so we do not do exisiting check here.
	// https://learn.microsoft.com/en-us/azure/sentinel/configure-fusion-rules#configure-scheduled-analytics-rules-for-fusion-detections
	params := alertrules.FusionAlertRule{
		Properties: &alertrules.FusionAlertRuleProperties{
			AlertRuleTemplateName: d.Get("alert_rule_template_guid").(string),
			Enabled:               d.Get("enabled").(bool),
			SourceSettings:        expandFusionSourceSettings(d.Get("source").([]interface{})),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleFusionRead(d, meta)
}

func resourceSentinelAlertRuleFusionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload, ok := resp.Model.(alertrules.FusionAlertRule)
	if !ok {
		return fmt.Errorf("retrieving %s: expected an alert rule of type `Fusion`, got %q", id, pointer.From(resp.Model.AlertRule().Type))
	}

	if payload.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if d.HasChange("alert_rule_template_guid") {
		payload.Properties.AlertRuleTemplateName = d.Get("alert_rule_template_guid").(string)
	}

	if d.HasChange("enabled") {
		payload.Properties.Enabled = d.Get("enabled").(bool)
	}

	if d.HasChange("source") {
		payload.Properties.SourceSettings = expandFusionSourceSettings(d.Get("source").([]interface{}))
	}

	// The `Description` is read-only but not specified on the Swagger, tracked on: https://github.com/Azure/azure-rest-api-specs/issues/31330
	payload.Properties.Description = nil
	payload.Properties.DisplayName = nil
	payload.Properties.LastModifiedUtc = nil
	payload.Properties.Severity = nil
	payload.Properties.Tactics = nil

	if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceSentinelAlertRuleFusionRead(d, meta)
}

func resourceSentinelAlertRuleFusionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
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

	workspaceId := alertrules.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

	if !features.FivePointOh() {
		d.Set("name", id.RuleId)
	}
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if model := resp.Model; model != nil {
		if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindFusion); err != nil {
			return fmt.Errorf("asserting alert rule of %s: %+v", id, err)
		}

		if rule, ok := model.(alertrules.FusionAlertRule); ok {
			if prop := rule.Properties; prop != nil {
				d.Set("enabled", prop.Enabled)
				d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)
				if err := d.Set("source", flattenFusionSourceSettings(prop.SourceSettings)); err != nil {
					return fmt.Errorf("setting `source`: %v", err)
				}
			}
		}
	}

	return nil
}

func resourceSentinelAlertRuleFusionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandFusionSourceSettings(input []interface{}) *[]alertrules.FusionSourceSettings {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.FusionSourceSettings, 0)

	for _, e := range input {
		e := e.(map[string]interface{})
		setting := alertrules.FusionSourceSettings{
			Enabled:        e["enabled"].(bool),
			SourceName:     e["name"].(string),
			SourceSubTypes: expandFusionSourceSubTypes(e["sub_type"].([]interface{})),
		}
		result = append(result, setting)
	}

	return &result
}

func expandFusionSourceSubTypes(input []interface{}) *[]alertrules.FusionSourceSubTypeSetting {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.FusionSourceSubTypeSetting, 0)

	for _, e := range input {
		e := e.(map[string]interface{})
		setting := alertrules.FusionSourceSubTypeSetting{
			Enabled:           e["enabled"].(bool),
			SourceSubTypeName: e["name"].(string),
			SeverityFilters: alertrules.FusionSubTypeSeverityFilter{
				Filters: expandFusionSubTypeSeverityFiltersItems(e["severities_allowed"].(*pluginsdk.Set).List()),
			},
		}
		result = append(result, setting)
	}

	return &result
}

func expandFusionSubTypeSeverityFiltersItems(input []interface{}) *[]alertrules.FusionSubTypeSeverityFiltersItem {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.FusionSubTypeSeverityFiltersItem, 0)

	// We can't simply remove the disabled properties in the request, as that will be reflected to the backend model (i.e. those unspecified severity will be absent also).
	// As any absent severity then will not be shown in the Portal when users try to edit the alert rule. The drop down menu won't show these absent severities...
	filters := map[string]bool{}
	for _, e := range alertrules.PossibleValuesForAlertSeverity() {
		filters[e] = false
	}

	for _, e := range input {
		filters[e.(string)] = true
	}

	for severity, enabled := range filters {
		item := alertrules.FusionSubTypeSeverityFiltersItem{
			Enabled:  enabled,
			Severity: alertrules.AlertSeverity(severity),
		}
		result = append(result, item)
	}

	return &result
}

func flattenFusionSourceSettings(input *[]alertrules.FusionSourceSettings) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, map[string]interface{}{
			"name":     e.SourceName,
			"enabled":  e.Enabled,
			"sub_type": flattenFusionSourceSubTypes(e.SourceSubTypes),
		})
	}

	return output
}

func flattenFusionSourceSubTypes(input *[]alertrules.FusionSourceSubTypeSetting) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, map[string]interface{}{
			"name":               e.SourceSubTypeName,
			"enabled":            e.Enabled,
			"severities_allowed": flattenFusionSubTypeSeverityFiltersItems(e.SeverityFilters.Filters),
		})
	}

	return output
}

func flattenFusionSubTypeSeverityFiltersItems(input *[]alertrules.FusionSubTypeSeverityFiltersItem) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		if e.Enabled {
			output = append(output, string(e.Severity))
		}
	}

	return output
}
