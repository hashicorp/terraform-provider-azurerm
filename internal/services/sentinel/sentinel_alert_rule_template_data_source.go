// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertruletemplates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSentinelAlertRuleTemplate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSentinelAlertRuleTemplateRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SentinelAlertRuleTemplateV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"scheduled_template": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tactics": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"severity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"query": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"query_frequency": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"query_period": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"trigger_operator": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"trigger_threshold": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"security_incident_template": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"product_filter": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"nrt_template": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tactics": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"severity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"query": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSentinelAlertRuleTemplateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRuleTemplatesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	var resp alertruletemplates.AlertRuleTemplate
	var nameToLog string
	if name != "" {
		nameToLog = name
		resp, err = getAlertRuleTemplateByName(ctx, client, workspaceID, name)
		if err != nil {
			return fmt.Errorf("an Alert Rule Template named %q was not found", name)
		}
	} else {
		nameToLog = displayName
		var realName *string
		resp, realName, err = getAlertRuleTemplateByDisplayName(ctx, client, workspaceID, displayName)
		if err != nil {
			return fmt.Errorf("an Alert Rule Template with the Display Name %q was not found", displayName)
		}
		name = *realName
	}

	id := alertruletemplates.NewAlertRuleTemplateID(meta.(*clients.Client).Account.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
	d.SetId(id.ID())

	switch template := resp.(type) {
	case alertruletemplates.MLBehaviorAnalyticsAlertRuleTemplate:
		setForMLBehaviorAnalyticsAlertRuleTemplate(d, &template)
	case alertruletemplates.FusionAlertRuleTemplate:
		setForFusionAlertRuleTemplate(d, &template)
	case alertruletemplates.MicrosoftSecurityIncidentCreationAlertRuleTemplate:
		setForMsSecurityIncidentAlertRuleTemplate(d, &template)
	case alertruletemplates.ScheduledAlertRuleTemplate:
		setForScheduledAlertRuleTemplate(d, &template)
	case alertruletemplates.NrtAlertRuleTemplate:
		setForNrtAlertRuleTemplate(d, &template)
	case alertruletemplates.ThreatIntelligenceAlertRuleTemplate:
		setForThreatIntelligenceAlertRuleTemplate(d, &template)
	default:
		return fmt.Errorf("unknown template type of Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q) ID", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroupName)
	}

	if err != nil {
		return fmt.Errorf("setting ResourceData for Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q) ID: %+v", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroupName, err)
	}

	return nil
}

func getAlertRuleTemplateByName(ctx context.Context, client *alertruletemplates.AlertRuleTemplatesClient, workspaceID *workspaces.WorkspaceId, name string) (res alertruletemplates.AlertRuleTemplate, err error) {
	id := alertruletemplates.NewAlertRuleTemplateID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
	template, err := client.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return template.Model, nil
}

func getAlertRuleTemplateByDisplayName(ctx context.Context, client *alertruletemplates.AlertRuleTemplatesClient, workspaceID *workspaces.WorkspaceId, displayName string) (res alertruletemplates.AlertRuleTemplate, name *string, err error) {
	templates, err := client.ListComplete(ctx, alertruletemplates.WorkspaceId(*workspaceID))
	if err != nil {
		return nil, nil, err
	}
	var results []alertruletemplates.AlertRuleTemplate
	for _, item := range templates.Items {
		switch template := item.(type) {
		case alertruletemplates.FusionAlertRuleTemplate:
			if template.Properties != nil && pointer.From(template.Properties.DisplayName) == displayName {
				results = append(results, item)
				if template.Name != nil {
					name = template.Name
				}
			}
		case alertruletemplates.MLBehaviorAnalyticsAlertRuleTemplate:
			if template.Properties != nil && pointer.From(template.Properties.DisplayName) == displayName {
				results = append(results, item)
				if template.Name != nil {
					name = template.Name
				}
			}
		case alertruletemplates.MicrosoftSecurityIncidentCreationAlertRuleTemplate:
			if template.Properties != nil && pointer.From(template.Properties.DisplayName) == displayName {
				results = append(results, item)
				if template.Name != nil {
					name = template.Name
				}
			}
		case alertruletemplates.ScheduledAlertRuleTemplate:
			if template.Properties != nil && pointer.From(template.Properties.DisplayName) == displayName {
				results = append(results, item)
				if template.Name != nil {
					name = template.Name
				}
			}
		case alertruletemplates.NrtAlertRuleTemplate:
			if template.Properties != nil && pointer.From(template.Properties.DisplayName) == displayName {
				results = append(results, item)
				if template.Name != nil {
					name = template.Name
				}
			}
		case alertruletemplates.ThreatIntelligenceAlertRuleTemplate:
			if template.Properties != nil && pointer.From(template.Properties.DisplayName) == displayName {
				results = append(results, item)
				if template.Name != nil {
					name = template.Name
				}
			}
		}
	}

	if len(results) == 0 {
		return nil, name, fmt.Errorf("no Alert Rule Template found with display name: %s", displayName)
	}
	if len(results) > 1 {
		return nil, name, fmt.Errorf("more than one Alert Rule Template found with display name: %s", displayName)
	}
	return results[0], name, nil
}

func setForScheduledAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.ScheduledAlertRuleTemplate) {
	d.Set("name", template.Name)
	if props := template.Properties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("scheduled_template", flattenScheduledAlertRuleTemplate(props))
	}
}

func setForNrtAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.NrtAlertRuleTemplate) {
	d.Set("name", template.Name)
	if props := template.Properties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("nrt_template", flattenNrtAlertRuleTemplate(props))
	}
}

func setForMsSecurityIncidentAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.MicrosoftSecurityIncidentCreationAlertRuleTemplate) {
	d.Set("name", template.Name)
	if props := template.Properties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("security_incident_template", flattenMsSecurityIncidentAlertRuleTemplate(props))
	}
}

func setForFusionAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.FusionAlertRuleTemplate) {
	d.Set("name", template.Name)
	if props := template.Properties; props != nil {
		d.Set("display_name", props.DisplayName)
	}
}

func setForMLBehaviorAnalyticsAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.MLBehaviorAnalyticsAlertRuleTemplate) {
	d.Set("name", template.Name)
	if props := template.Properties; props != nil {
		d.Set("display_name", props.DisplayName)
	}
}

func setForThreatIntelligenceAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.ThreatIntelligenceAlertRuleTemplate) {
	d.Set("name", template.Name)
	if props := template.Properties; props != nil {
		d.Set("display_name", props.DisplayName)
	}
}

func flattenScheduledAlertRuleTemplate(input *alertruletemplates.ScheduledAlertRuleTemplateProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"description":       pointer.From(input.Description),
			"tactics":           pointer.FromEnumSlice(input.Tactics),
			"severity":          pointer.FromEnum(input.Severity),
			"query":             pointer.From(input.Query),
			"query_frequency":   pointer.From(input.QueryFrequency),
			"query_period":      pointer.From(input.QueryPeriod),
			"trigger_operator":  pointer.FromEnum(input.TriggerOperator),
			"trigger_threshold": pointer.From(input.TriggerThreshold),
		},
	}
}

func flattenNrtAlertRuleTemplate(input *alertruletemplates.NrtAlertRuleTemplateProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"description": pointer.From(input.Description),
			"tactics":     pointer.FromEnumSlice(input.Tactics),
			"severity":    string(input.Severity),
			"query":       input.Query,
		},
	}
}

func flattenMsSecurityIncidentAlertRuleTemplate(input *alertruletemplates.MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"description":    pointer.From(input.Description),
			"product_filter": pointer.FromEnum(input.ProductFilter),
		},
	}
}
