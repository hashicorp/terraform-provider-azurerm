package sentinel

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceSentinelAlertRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSentinelAlertRuleTemplateRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},

			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
			},

			"scheduled_template": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tactics": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_frequency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"security_incident_template": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSentinelAlertRuleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRuleTemplatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	// Either "name" or "display_name" must have been specified, constrained by the schema.
	var resp securityinsight.BasicAlertRuleTemplate
	var nameToLog string
	if name != "" {
		nameToLog = name
		resp, err = getAlertRuleTemplateByName(ctx, client, workspaceID, name)
	} else {
		nameToLog = displayName
		resp, err = getAlertRuleTemplateByDisplayName(ctx, client, workspaceID, displayName)
	}
	if err != nil {
		return fmt.Errorf("retrieving Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q): %+v", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroup, err)
	}

	switch template := resp.(type) {
	case securityinsight.MLBehaviorAnalyticsAlertRuleTemplate:
		err = setForMLBehaviorAnalyticsAlertRuleTemplate(d, &template)
	case securityinsight.FusionAlertRuleTemplate:
		err = setForFusionAlertRuleTemplate(d, &template)
	case securityinsight.MicrosoftSecurityIncidentCreationAlertRuleTemplate:
		err = setForMsSecurityIncidentAlertRuleTemplate(d, &template)
	case securityinsight.ScheduledAlertRuleTemplate:
		err = setForScheduledAlertRuleTemplate(d, &template)
	default:
		return fmt.Errorf("unknown template type of Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q) ID", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroup)
	}

	if err != nil {
		return fmt.Errorf("setting ResourceData for Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q) ID", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroup)
	}

	return nil
}

func getAlertRuleTemplateByName(ctx context.Context, client *securityinsight.AlertRuleTemplatesClient, workspaceID *loganalyticsParse.LogAnalyticsWorkspaceId, name string) (res securityinsight.BasicAlertRuleTemplate, err error) {
	template, err := client.Get(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName, name)
	if err != nil {
		return nil, err
	}
	return template.Value, nil
}

func getAlertRuleTemplateByDisplayName(ctx context.Context, client *securityinsight.AlertRuleTemplatesClient, workspaceID *loganalyticsParse.LogAnalyticsWorkspaceId, name string) (res securityinsight.BasicAlertRuleTemplate, err error) {
	templates, err := client.ListComplete(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName)
	if err != nil {
		return nil, err
	}
	var results []securityinsight.BasicAlertRuleTemplate
	for templates.NotDone() {
		template := templates.Value()
		switch template := template.(type) {
		case securityinsight.FusionAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == name {
				results = append(results, templates.Value())
			}
		case securityinsight.MLBehaviorAnalyticsAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == name {
				results = append(results, templates.Value())
			}
		case securityinsight.MicrosoftSecurityIncidentCreationAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == name {
				results = append(results, templates.Value())
			}
		case securityinsight.ScheduledAlertRuleTemplate:
			if template.DisplayName != nil && *template.DisplayName == name {
				results = append(results, templates.Value())
			}
		}

		if err := templates.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("iterating Alert Rule Templates: %+v", err)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no Alert Rule Template found with display name: %s", name)
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("more than one Alert Rule Template found with display name: %s", name)
	}
	return results[0], nil
}

func setForScheduledAlertRuleTemplate(d *schema.ResourceData, template *securityinsight.ScheduledAlertRuleTemplate) error {
	if template.ID == nil || *template.ID == "" {
		return errors.New("empty or nil ID")
	}
	id, err := parse.SentinelAlertRuleTemplateID(*template.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
	return d.Set("scheduled_template", flattenScheduledAlertRuleTemplate(template.ScheduledAlertRuleTemplateProperties))
}

func setForMsSecurityIncidentAlertRuleTemplate(d *schema.ResourceData, template *securityinsight.MicrosoftSecurityIncidentCreationAlertRuleTemplate) error {
	if template.ID == nil || *template.ID == "" {
		return errors.New("empty or nil ID")
	}
	id, err := parse.SentinelAlertRuleTemplateID(*template.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
	return d.Set("security_incident_template", flattenMsSecurityIncidentAlertRuleTemplate(template.MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties))
}

func setForFusionAlertRuleTemplate(d *schema.ResourceData, template *securityinsight.FusionAlertRuleTemplate) error {
	if template.ID == nil || *template.ID == "" {
		return errors.New("empty or nil ID")
	}
	id, err := parse.SentinelAlertRuleTemplateID(*template.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
	return nil
}

func setForMLBehaviorAnalyticsAlertRuleTemplate(d *schema.ResourceData, template *securityinsight.MLBehaviorAnalyticsAlertRuleTemplate) error {
	if template.ID == nil || *template.ID == "" {
		return errors.New("empty or nil ID")
	}
	id, err := parse.SentinelAlertRuleTemplateID(*template.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID())
	d.Set("name", template.Name)
	d.Set("display_name", template.DisplayName)
	return nil
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
		tactics = flattenAlertRuleScheduledTactics(input.Tactics)
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
