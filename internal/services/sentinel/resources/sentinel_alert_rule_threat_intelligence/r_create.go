// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_threat_intelligence

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertruletemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
)

func (a Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config AlertRuleThreatIntelligenceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := sentinelClient.From(metadata.Client).AlertRulesClient
			alertRuleTemplatesClient := sentinelClient.From(metadata.Client).AlertRuleTemplatesClient

			workspaceID, err := workspaces.ParseWorkspaceID(config.WorkspaceId)
			if err != nil {
				return err
			}

			id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, config.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
					Severity:              pointer.ToEnum[alertrules.AlertSeverity](string(props.Severity)),
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
