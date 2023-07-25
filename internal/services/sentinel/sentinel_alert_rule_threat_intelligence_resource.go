// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"time"

	alertruletemplates "github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2021-09-01-preview/securityinsight" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AlertRuleThreatIntelligenceModel struct {
	Name         string `tfschema:"name"`
	WorkspaceId  string `tfschema:"log_analytics_workspace_id"`
	TemplateName string `tfschema:"alert_rule_template_guid"`
	Enabled      bool   `tfschema:"enabled"`
}

type AlertRuleThreatIntelligenceResource struct{}

var _ sdk.ResourceWithCustomImporter = AlertRuleThreatIntelligenceResource{}
var _ sdk.ResourceWithUpdate = AlertRuleThreatIntelligenceResource{}

func (a AlertRuleThreatIntelligenceResource) ModelObject() interface{} {
	return &AlertRuleThreatIntelligenceModel{}
}

func (a AlertRuleThreatIntelligenceResource) ResourceType() string {
	return "azurerm_sentinel_alert_rule_threat_intelligence"
}

func (a AlertRuleThreatIntelligenceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return alertrules.ValidateAlertRuleID
}

func (a AlertRuleThreatIntelligenceResource) CustomImporter() sdk.ResourceRunFunc {
	return importSentinelAlertRuleForTypedSdk(alertrules.AlertRuleKindThreatIntelligence)
}

func (a AlertRuleThreatIntelligenceResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}

func (a AlertRuleThreatIntelligenceResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (a AlertRuleThreatIntelligenceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel AlertRuleThreatIntelligenceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.Sentinel.AlertRulesClient

			workspaceID, err := workspaces.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return err
			}

			id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, metaModel.Name)

			resp, err := client.AlertRulesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_threat_intelligence", id.ID())
			}

			template, err := fetchAlertRuleThreatIntelligenceTemplate(ctx, metadata.Client.Sentinel.AlertRuleTemplatesClient, *workspaceID, metaModel.TemplateName)
			if err != nil {
				return fmt.Errorf("fetching severity from template: %+v", err)
			}

			tactics := make([]alertrules.AttackTactic, 0)
			if template.Tactics != nil {
				for _, t := range *template.Tactics {
					tactics = append(tactics, alertrules.AttackTactic(t))
				}
			}

			severity := alertrules.AlertSeverity(template.Severity)
			param := alertrules.ThreatIntelligenceAlertRule{
				Properties: &alertrules.ThreatIntelligenceAlertRuleProperties{
					Enabled:               metaModel.Enabled,
					AlertRuleTemplateName: metaModel.TemplateName,
					Severity:              &severity,
					DisplayName:           template.DisplayName,
					Description:           template.Description,
					Tactics:               &tactics,
				},
			}

			if _, err := client.AlertRulesCreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %q: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (a AlertRuleThreatIntelligenceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.AlertRulesClient

			id, err := alertrules.ParseAlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			resp, err := client.AlertRulesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %+v", err)
			}

			if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindThreatIntelligence); err != nil {
				return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
			}
			rule := (*resp.Model).(alertrules.ThreatIntelligenceAlertRule)

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

			state := AlertRuleThreatIntelligenceModel{
				Name:        id.RuleId,
				WorkspaceId: workspaceId.ID(),
			}

			if prop := rule.Properties; prop != nil {
				state.Enabled = prop.Enabled
				state.TemplateName = prop.AlertRuleTemplateName
			}

			return metadata.Encode(&state)
		},
	}
}

func (a AlertRuleThreatIntelligenceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.AlertRulesClient

			id, err := alertrules.ParseAlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			if _, err := client.AlertRulesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %+v", err)
			}

			return nil
		},
	}
}

func (a AlertRuleThreatIntelligenceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel AlertRuleThreatIntelligenceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.Sentinel.AlertRulesClient

			id, err := alertrules.ParseAlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}
			resp, err := client.AlertRulesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %+v", err)
			}
			if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindThreatIntelligence); err != nil {
				return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
			}
			rule := (*resp.Model).(alertrules.ThreatIntelligenceAlertRule)

			if metadata.ResourceData.HasChange("enabled") {
				rule.Properties.Enabled = metaModel.Enabled
			}
			if metadata.ResourceData.HasChange("template_name") {
				rule.Properties.AlertRuleTemplateName = metaModel.TemplateName
			}

			param := alertrules.ThreatIntelligenceAlertRule{
				Properties: &alertrules.ThreatIntelligenceAlertRuleProperties{
					Enabled:               rule.Properties.Enabled,
					AlertRuleTemplateName: rule.Properties.AlertRuleTemplateName,
					Severity:              rule.Properties.Severity,
					DisplayName:           rule.Properties.DisplayName,
					Description:           rule.Properties.Description,
					Tactics:               rule.Properties.Tactics,
				},
			}

			if _, err := client.AlertRulesCreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}
}

// workaround for https://github.com/Azure/azure-rest-api-specs/issues/16615
func fetchAlertRuleThreatIntelligenceTemplate(ctx context.Context, templateClient *alertruletemplates.AlertRuleTemplatesClient, workspaceId workspaces.WorkspaceId, templateName string) (alertruletemplates.ThreatIntelligenceAlertRuleTemplate, error) {
	foo := alertruletemplates.ThreatIntelligenceAlertRuleTemplate{}
	id := parse.NewSentinelAlertRuleTemplateID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, templateName)

	resp, err := templateClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.AlertRuleTemplateName)
	if err != nil {
		return foo, fmt.Errorf("reading %q: %+v", id, err)
	}

	v, ok := resp.Value.(alertruletemplates.ThreatIntelligenceAlertRuleTemplate)
	if !ok {
		return foo, fmt.Errorf("reading %q: type mismatch", id)
	}

	return v, nil
}
