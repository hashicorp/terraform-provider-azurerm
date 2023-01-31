package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
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
	return validate.AlertRuleID
}

func (a AlertRuleThreatIntelligenceResource) CustomImporter() sdk.ResourceRunFunc {
	return importSentinelAlertRuleForTypedSdk(securityinsight.AlertRuleKindThreatIntelligence)
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

			id := parse.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, metaModel.Name)

			param := securityinsight.ThreatIntelligenceAlertRule{
				Kind: securityinsight.KindBasicAlertRuleKindThreatIntelligence,
				ThreatIntelligenceAlertRuleProperties: &securityinsight.ThreatIntelligenceAlertRuleProperties{
					Enabled:               utils.Bool(metaModel.Enabled),
					AlertRuleTemplateName: utils.String(metaModel.TemplateName),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, param); err != nil {
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

			id, err := parse.AlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %+v", err)
			}

			if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindThreatIntelligence); err != nil {
				return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
			}
			rule := resp.Value.(securityinsight.ThreatIntelligenceAlertRule)

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

			state := AlertRuleThreatIntelligenceModel{
				Name:        id.Name,
				WorkspaceId: workspaceId.ID(),
			}

			if prop := rule.ThreatIntelligenceAlertRuleProperties; prop != nil {
				if prop.Enabled != nil {
					state.Enabled = *prop.Enabled
				}
				if prop.AlertRuleTemplateName != nil {
					state.TemplateName = *prop.AlertRuleTemplateName
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

			id, err := parse.AlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
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

			id, err := parse.AlertRuleID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}
			resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				return fmt.Errorf("reading %+v", err)
			}
			if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindThreatIntelligence); err != nil {
				return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
			}
			rule := resp.Value.(securityinsight.ThreatIntelligenceAlertRule)

			if metadata.ResourceData.HasChange("enabled") {
				rule.ThreatIntelligenceAlertRuleProperties.Enabled = utils.Bool(metaModel.Enabled)
			}
			if metadata.ResourceData.HasChange("template_name") {
				rule.ThreatIntelligenceAlertRuleProperties.AlertRuleTemplateName = utils.String(metaModel.TemplateName)
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, rule); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}
}
