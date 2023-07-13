// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2021-09-01-preview/securityinsight" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	// Either "name" or "display_name" must have been specified, constrained by the pluginsdk.
	var resp securityinsight.BasicAlertRuleTemplate
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

	id := parse.NewSentinelAlertRuleTemplateID(subscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	switch template := resp.(type) {
	case securityinsight.MLBehaviorAnalyticsAlertRuleTemplate:
		setForMLBehaviorAnalyticsAlertRuleTemplate(d, id, &template)
	case securityinsight.FusionAlertRuleTemplate:
		setForFusionAlertRuleTemplate(d, id, &template)
	case securityinsight.MicrosoftSecurityIncidentCreationAlertRuleTemplate:
		err = setForMsSecurityIncidentAlertRuleTemplate(d, id, &template)
	case securityinsight.ScheduledAlertRuleTemplate:
		err = setForScheduledAlertRuleTemplate(d, id, &template)
	case securityinsight.NrtAlertRuleTemplate:
		err = setForNrtAlertRuleTemplate(d, id, &template)
	case securityinsight.ThreatIntelligenceAlertRuleTemplate:
		setForThreatIntelligenceAlertRuleTemplate(d, id, &template)
	default:
		return fmt.Errorf("unknown template type of Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q) ID", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroupName)
	}

	if err != nil {
		return fmt.Errorf("setting ResourceData for Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q) ID: %+v", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroupName, err)
	}

	return nil
}

func getAlertRuleTemplateByName(ctx context.Context, client *securityinsight.AlertRuleTemplatesClient, workspaceID *workspaces.WorkspaceId, name string) (res securityinsight.BasicAlertRuleTemplate, err error) {
	template, err := client.Get(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
	if err != nil {
		return nil, err
	}
	return template.Value, nil
}

func getAlertRuleTemplateByDisplayName(ctx context.Context, client *securityinsight.AlertRuleTemplatesClient, workspaceID *workspaces.WorkspaceId, displayName string) (res securityinsight.BasicAlertRuleTemplate, name *string, err error) {
	templates, err := client.ListComplete(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName)
	if err != nil {
		return nil, nil, err
	}
	var results []securityinsight.BasicAlertRuleTemplate
	for templates.NotDone() {
		template := templates.Value()
		switch template := template.(type) {
		case securityinsight.FusionAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == displayName {
				results = append(results, templates.Value())
				if template.Name != nil {
					name = template.Name
				}
			}
		case securityinsight.MLBehaviorAnalyticsAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == displayName {
				results = append(results, templates.Value())
				if template.Name != nil {
					name = template.Name
				}
			}
		case securityinsight.MicrosoftSecurityIncidentCreationAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == displayName {
				results = append(results, templates.Value())
				if template.Name != nil {
					name = template.Name
				}
			}
		case securityinsight.ScheduledAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == displayName {
				results = append(results, templates.Value())
				if template.Name != nil {
					name = template.Name
				}
			}
		case securityinsight.NrtAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == displayName {
				results = append(results, templates.Value())
				if template.Name != nil {
					name = template.Name
				}
			}
		case securityinsight.ThreatIntelligenceAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == displayName {
				results = append(results, templates.Value())
				if template.Name != nil {
					name = template.Name
				}
			}
		}

		if err := templates.NextWithContext(ctx); err != nil {
			return nil, nil, fmt.Errorf("iterating Alert Rule Templates: %+v", err)
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

func setForScheduledAlertRuleTemplate(d *pluginsdk.ResourceData, id parse.SentinelAlertRuleTemplateId, template *securityinsight.ScheduledAlertRuleTemplate) error {
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
	return d.Set("scheduled_template", flattenScheduledAlertRuleTemplate(template.ScheduledAlertRuleTemplateProperties))
}

func setForNrtAlertRuleTemplate(d *pluginsdk.ResourceData, id parse.SentinelAlertRuleTemplateId, template *securityinsight.NrtAlertRuleTemplate) error {
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
	return d.Set("nrt_template", flattenNrtAlertRuleTemplate(template.NrtAlertRuleTemplateProperties))
}

func setForMsSecurityIncidentAlertRuleTemplate(d *pluginsdk.ResourceData, id parse.SentinelAlertRuleTemplateId, template *securityinsight.MicrosoftSecurityIncidentCreationAlertRuleTemplate) error {
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
	return d.Set("security_incident_template", flattenMsSecurityIncidentAlertRuleTemplate(template.MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties))
}

func setForFusionAlertRuleTemplate(d *pluginsdk.ResourceData, id parse.SentinelAlertRuleTemplateId, template *securityinsight.FusionAlertRuleTemplate) {
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
}

func setForMLBehaviorAnalyticsAlertRuleTemplate(d *pluginsdk.ResourceData, id parse.SentinelAlertRuleTemplateId, template *securityinsight.MLBehaviorAnalyticsAlertRuleTemplate) {
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
}

func setForThreatIntelligenceAlertRuleTemplate(d *pluginsdk.ResourceData, id parse.SentinelAlertRuleTemplateId, template *securityinsight.ThreatIntelligenceAlertRuleTemplate) {
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
}

func flattenScheduledAlertRuleTemplate(input *securityinsight.ScheduledAlertRuleTemplateProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	tactics := []interface{}{}
	if input.Tactics != nil {
		tactics = flattenAlertRuleTacticsForTemplate(input.Tactics)
	}

	query := ""
	if input.Query != nil {
		query = *input.Query
	}

	queryFrequency := ""
	if input.QueryFrequency != nil {
		queryFrequency = *input.QueryFrequency
	}

	queryPeriod := ""
	if input.QueryPeriod != nil {
		queryPeriod = *input.QueryPeriod
	}

	triggerThreshold := 0
	if input.TriggerThreshold != nil {
		triggerThreshold = int(*input.TriggerThreshold)
	}

	return []interface{}{
		map[string]interface{}{
			"description":       description,
			"tactics":           tactics,
			"severity":          string(input.Severity),
			"query":             query,
			"query_frequency":   queryFrequency,
			"query_period":      queryPeriod,
			"trigger_operator":  string(input.TriggerOperator),
			"trigger_threshold": triggerThreshold,
		},
	}
}

func flattenNrtAlertRuleTemplate(input *securityinsight.NrtAlertRuleTemplateProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	tactics := []interface{}{}
	if input.Tactics != nil {
		tactics = flattenAlertRuleTacticsForTemplate(input.Tactics)
	}

	query := ""
	if input.Query != nil {
		query = *input.Query
	}

	return []interface{}{
		map[string]interface{}{
			"description": description,
			"tactics":     tactics,
			"severity":    string(input.Severity),
			"query":       query,
		},
	}
}

func flattenMsSecurityIncidentAlertRuleTemplate(input *securityinsight.MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	return []interface{}{
		map[string]interface{}{
			"description":    description,
			"product_filter": string(input.ProductFilter),
		},
	}
}

func flattenAlertRuleTacticsForTemplate(input *[]securityinsight.AttackTactic) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, string(e))
	}

	return output
}
