// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertruletemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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

var (
	_ sdk.ResourceWithCustomImporter = AlertRuleThreatIntelligenceResource{}
	_ sdk.ResourceWithUpdate         = AlertRuleThreatIntelligenceResource{}
)

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
			var config AlertRuleThreatIntelligenceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.Sentinel.AlertRulesClient
			alertRuleTemplatesClient := metadata.Client.Sentinel.AlertRuleTemplatesClient

			workspaceID, err := workspaces.ParseWorkspaceID(config.WorkspaceId)
			if err != nil {
				return err
			}

			id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				resp, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(resp.HttpResponse) {
						return fmt.Errorf("checking for existing %q: %+v", id, err)
					}
				}

				if !response.WasNotFound(resp.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_threat_intelligence", id.ID())
				}
			}

			templateID := alertruletemplates.NewAlertRuleTemplateID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, config.TemplateName)
			template, err := alertRuleTemplatesClient.Get(ctx, templateID)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", templateID, err)
			}

			if template.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", templateID)
			}

			v, ok := template.Model.(alertruletemplates.ThreatIntelligenceAlertRuleTemplate)
			if !ok {
				return fmt.Errorf("retrieving %s: unexpected template type (%T)", templateID, template.Model)
			}

			if v.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", templateID)
			}
			props := v.Properties

			tactics := make([]alertrules.AttackTactic, 0)
			for _, t := range pointer.From(props.Tactics) {
				tactics = append(tactics, alertrules.AttackTactic(t))
			}

			param := alertrules.ThreatIntelligenceAlertRule{
				Properties: &alertrules.ThreatIntelligenceAlertRuleProperties{
					Enabled:               config.Enabled,
					AlertRuleTemplateName: config.TemplateName,
					Severity:              pointer.To(alertrules.AlertSeverity(props.Severity)),
					DisplayName:           props.DisplayName,
					Description:           props.Description,
					Tactics:               &tactics,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
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

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindThreatIntelligence); err != nil {
				return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

			state := AlertRuleThreatIntelligenceModel{
				Name:        id.RuleId,
				WorkspaceId: workspaceId.ID(),
			}

			if rule, ok := resp.Model.(alertrules.ThreatIntelligenceAlertRule); ok {
				if prop := rule.Properties; prop != nil {
					state.Enabled = prop.Enabled
					state.TemplateName = prop.AlertRuleTemplateName
				}
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

			if _, err := client.Delete(ctx, *id); err != nil {
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
			var config AlertRuleThreatIntelligenceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.Sentinel.AlertRulesClient

			id, err := alertrules.ParseAlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindThreatIntelligence); err != nil {
				return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
			}

			rule := resp.Model.(alertrules.ThreatIntelligenceAlertRule)

			if rule.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			props := rule.Properties

			if metadata.ResourceData.HasChange("enabled") {
				props.Enabled = config.Enabled
			}

			param := alertrules.ThreatIntelligenceAlertRule{
				Properties: &alertrules.ThreatIntelligenceAlertRuleProperties{
					Enabled:               props.Enabled,
					AlertRuleTemplateName: props.AlertRuleTemplateName,
					Severity:              props.Severity,
					DisplayName:           props.DisplayName,
					Description:           props.Description,
					Tactics:               props.Tactics,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}
}
