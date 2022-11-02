package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2022-01-01-preview/securityinsight"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSentinelAlertRuleFusion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleFusionCreateUpdate,
		Read:   resourceSentinelAlertRuleFusionRead,
		Update: resourceSentinelAlertRuleFusionCreateUpdate,
		Delete: resourceSentinelAlertRuleFusionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.AlertRuleID(id)
			return err
		}, importSentinelAlertRule(securityinsight.AlertRuleKindFusion)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
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

			"source_setting": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				// Service will auto-fill this if not given in request, based on the "alert_rule_template_guid".
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"source_sub_type": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"enabled_severities": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(
												[]string{
													string(securityinsight.AlertSeverityHigh),
													string(securityinsight.AlertSeverityMedium),
													string(securityinsight.AlertSeverityLow),
													string(securityinsight.AlertSeverityInformational),
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
}

func resourceSentinelAlertRuleFusionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Fusion %q: %+v", id, err)
			}
		}

		id := alertRuleID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_fusion", *id)
		}
	}

	params := securityinsight.FusionAlertRule{
		Kind: securityinsight.KindBasicAlertRuleKindFusion,
		FusionAlertRuleProperties: &securityinsight.FusionAlertRuleProperties{
			AlertRuleTemplateName: utils.String(d.Get("alert_rule_template_guid").(string)),
			Enabled:               utils.Bool(d.Get("enabled").(bool)),
			SourceSettings:        expandFusionSourceSettings(d.Get("source_setting").([]interface{})),
		},
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Fusion %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindFusion); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.CreateOrUpdate(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name, params); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleFusionRead(d, meta)
}

func resourceSentinelAlertRuleFusionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Sentinel Alert Rule Fusion %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindFusion); err != nil {
		return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
	}
	rule := resp.Value.(securityinsight.FusionAlertRule)

	d.Set("name", id.Name)

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if prop := rule.FusionAlertRuleProperties; prop != nil {
		d.Set("enabled", prop.Enabled)
		d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)
		if err := d.Set("source_setting", flattenFusionSourceSettings(prop.SourceSettings)); err != nil {
			return fmt.Errorf("setting `source_setting`: %v", err)
		}
	}

	return nil
}

func resourceSentinelAlertRuleFusionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	return nil
}

func expandFusionSourceSettings(input []interface{}) *[]securityinsight.FusionSourceSettings {
	if len(input) == 0 {
		return nil
	}

	result := make([]securityinsight.FusionSourceSettings, 0)

	for _, e := range input {
		e := e.(map[string]interface{})
		setting := securityinsight.FusionSourceSettings{
			Enabled:        utils.Bool(true),
			SourceName:     utils.String(e["name"].(string)),
			SourceSubTypes: expandFusionSourceSubTypes(e["source_sub_type"].([]interface{})),
		}
		result = append(result, setting)
	}

	return &result
}

func expandFusionSourceSubTypes(input []interface{}) *[]securityinsight.FusionSourceSubTypeSetting {
	if len(input) == 0 {
		return nil
	}

	result := make([]securityinsight.FusionSourceSubTypeSetting, 0)

	for _, e := range input {
		e := e.(map[string]interface{})
		setting := securityinsight.FusionSourceSubTypeSetting{
			Enabled:           utils.Bool(true),
			SourceSubTypeName: utils.String(e["name"].(string)),
			SeverityFilters: &securityinsight.FusionSubTypeSeverityFilter{
				Filters: expandFusionSubTypeSeverityFiltersItems(e["enabled_severities"].(*pluginsdk.Set).List()),
			},
		}
		result = append(result, setting)
	}

	return &result
}

func expandFusionSubTypeSeverityFiltersItems(input []interface{}) *[]securityinsight.FusionSubTypeSeverityFiltersItem {
	if len(input) == 0 {
		return nil
	}

	result := make([]securityinsight.FusionSubTypeSeverityFiltersItem, 0)

	for _, e := range input {
		item := securityinsight.FusionSubTypeSeverityFiltersItem{
			Enabled:  utils.Bool(true),
			Severity: securityinsight.AlertSeverity(e.(string)),
		}
		result = append(result, item)
	}

	return &result
}

func flattenFusionSourceSettings(input *[]securityinsight.FusionSourceSettings) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var name string
		if e.SourceName != nil {
			name = *e.SourceName
		}
		output = append(output, map[string]interface{}{
			"name":            name,
			"source_sub_type": flattenFusionSourceSubTypes(e.SourceSubTypes),
		})
	}

	return output
}

func flattenFusionSourceSubTypes(input *[]securityinsight.FusionSourceSubTypeSetting) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var name string
		if e.SourceSubTypeName != nil {
			name = *e.SourceSubTypeName
		}
		var enabledSeverities []interface{}
		if e.SeverityFilters != nil {
			enabledSeverities = flattenFusionSubTypeSeverityFiltersItems(e.SeverityFilters.Filters)
		}
		output = append(output, map[string]interface{}{
			"name":               name,
			"enabled_severities": enabledSeverities,
		})
	}

	return output
}

func flattenFusionSubTypeSeverityFiltersItems(input *[]securityinsight.FusionSubTypeSeverityFiltersItem) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		if e.Enabled != nil && *e.Enabled {
			output = append(output, string(e.Severity))
		}
	}

	return output
}
